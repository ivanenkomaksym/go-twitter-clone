package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"twitter-clone/internal/authn"
	"twitter-clone/internal/config"
	"twitter-clone/internal/messaging"
	"twitter-clone/internal/models"
	feedrepo "twitter-clone/internal/repositories/feed"
	tweetrepo "twitter-clone/internal/repositories/tweet"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	watermillHTTP "github.com/ThreeDotsLabs/watermill-http/pkg/http"
)

func StartRouter(configuration config.Configuration,
	tweetRepo tweetrepo.TweetRepository,
	feedRepo feedrepo.FeedRepository,
	messageHandler messaging.MessageHandler,
	authenticationValidator authn.IAuthenticationValidator) {
	logger := watermill.NewStdLogger(false, false)

	pub, sub, err := messageHandler.SetupMessageRouter(configuration, feedRepo, logger)
	if err != nil {
		panic(err)
	}

	normalizedDomain := strings.TrimPrefix(strings.TrimPrefix(configuration.AllowOrigin, "http://"), "https://")

	oauth2Router := authn.OAuth2Router{
		Authentication:          configuration.Authentication,
		RedirectURI:             configuration.RedirectURI,
		Domain:                  normalizedDomain,
		AuthenticationValidator: authenticationValidator,
	}

	httpRouter := Router{
		Config:                  configuration,
		AuthenticationValidator: authenticationValidator,
		OAuth2Router:            oauth2Router,
		Subscriber:              sub,
		Publisher:               Publisher{Publisher: pub},
		TweetRepo:               tweetRepo,
		FeedRepo:                feedRepo,
		Logger:                  logger,
	}

	mux := httpRouter.Mux()

	err = http.ListenAndServe(configuration.ApiServer.ApplicationUrl, mux)
	if err != nil {
		panic(err)
	}
}

type Router struct {
	Config                  config.Configuration
	AuthenticationValidator authn.IAuthenticationValidator
	OAuth2Router            authn.OAuth2Router
	Subscriber              message.Subscriber
	Publisher               IPublisher
	TweetRepo               tweetrepo.TweetRepository
	FeedRepo                feedrepo.FeedRepository
	Logger                  watermill.LoggerAdapter
}

func (router Router) Mux() *chi.Mux {
	r := chi.NewRouter()

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{router.Config.AllowOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	sseRouter, err := watermillHTTP.NewSSERouter(
		watermillHTTP.SSERouterConfig{
			UpstreamSubscriber: router.Subscriber,
			ErrorHandler:       watermillHTTP.DefaultErrorHandler,
		},
		router.Logger,
	)
	if err != nil {
		panic(err)
	}

	tweetStream := TweetStreamAdapter{
		repo:                    router.TweetRepo,
		authenticationValidator: router.AuthenticationValidator,
		logger:                  router.Logger,
	}
	feedStream := FeedStreamAdapter{
		repo:                    router.FeedRepo,
		authenticationValidator: router.AuthenticationValidator,
		logger:                  router.Logger,
	}
	allTweetsStream := AllTweetsStreamAdapter{
		repo:                    router.TweetRepo,
		authenticationValidator: router.AuthenticationValidator,
	}
	allFeedsStream := AllFeedsStreamAdapter{
		repo:                    router.FeedRepo,
		authenticationValidator: router.AuthenticationValidator,
		logger:                  router.Logger,
	}

	tweetHandler := sseRouter.AddHandler(messaging.TweetCreatedTopic, tweetStream)
	feedHandler := sseRouter.AddHandler(messaging.FeedUpdatedTopic, feedStream)
	allTweetsHandler := sseRouter.AddHandler(messaging.TweetUpdatedTopic, allTweetsStream)
	allFeedsHandler := sseRouter.AddHandler(messaging.FeedUpdatedTopic, allFeedsStream)

	if router.Config.Authentication.Enable {
		r.Route("/", func(r chi.Router) {
			r.Get("/auth/google/login", router.OAuth2Router.OauthGoogleLogin)
			r.Get("/auth/google/logout", router.OAuth2Router.OauthGoogleLogout)
			r.Get("/auth/google/callback", router.OAuth2Router.OauthGoogleCallback)
			r.Get("/auth/google/userinfo", router.OAuth2Router.OauthUserInfo)
		})
	}

	r.Route("/api", func(r chi.Router) {
		r.Post("/tweets", router.CreateTweet)
		r.Get("/tweets", allTweetsHandler)
		r.Get("/tweets/{tweetId}", tweetHandler)
		r.Delete("/tweets/{tweetId}", router.DeleteTweet)
		r.Get("/feeds/{name}", feedHandler)
		r.Get("/feeds", allFeedsHandler)
	})

	go func() {
		err := sseRouter.Run(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	<-sseRouter.Running()

	return r
}

func (router Router) CreateTweet(w http.ResponseWriter, r *http.Request) {
	user := router.AuthenticationValidator.ValidateAuthentication(w, r)
	if user == nil {
		return
	}

	var createTweetRequest models.CreateTweetRequest
	err := render.Decode(r, &createTweetRequest)
	if err != nil {
		logAndWriteError(router.Logger, w, err)
		return
	}

	createdTweet := router.TweetRepo.CreateTweet(createTweetRequest, *user)
	if createdTweet == nil {
		return
	}

	event := messaging.TweetCreated{
		Tweet:      *createdTweet,
		OccurredAt: time.Now().UTC(),
	}

	err = router.Publisher.Publish(messaging.TweetCreatedTopic, event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdTweet); err != nil {
		logAndWriteError(router.Logger, w, err)
	}
}

func (router Router) DeleteTweet(w http.ResponseWriter, r *http.Request) {
	user := router.AuthenticationValidator.ValidateAuthentication(w, r)
	if user == nil {
		return
	}

	tweetId := chi.URLParam(r, "tweetId")
	tweetToDelete := router.TweetRepo.GetTweetById(tweetId)
	if tweetToDelete == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var deleted = router.TweetRepo.DeleteTweet(tweetId)

	if !deleted {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {

		event := messaging.TweetDeleted{
			DeletedTweet: *tweetToDelete,
			OccurredAt:   time.Now().UTC(),
		}

		router.Logger.Info("Publishing tweet deleted event", watermill.LogFields{"event": event})
		err := router.Publisher.Publish(messaging.TweetDeletedTopic, event)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(204)
}

func logAndWriteError(logger watermill.LoggerAdapter, w http.ResponseWriter, err error) {
	logger.Error("Error", err, nil)
	w.WriteHeader(http.StatusInternalServerError)
}

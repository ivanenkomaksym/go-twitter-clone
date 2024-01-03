package api

import (
	"context"
	"net/http"
	"time"
	"twitter-clone/internal/config"
	"twitter-clone/internal/messaging"
	"twitter-clone/internal/models"
	"twitter-clone/internal/repositories"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	watermillHTTP "github.com/ThreeDotsLabs/watermill-http/pkg/http"
)

func StartRouter(configuration config.Configuration, tweetRepo repositories.TweetRepository, feedRepo repositories.FeedRepository) {
	logger := watermill.NewStdLogger(false, false)

	pub, sub, err := messaging.SetupMessageRouter(feedRepo, logger)
	if err != nil {
		panic(err)
	}

	httpRouter := Router{
		Subscriber: sub,
		Publisher:  Publisher{Publisher: pub},
		TweetRepo:  tweetRepo,
		FeedRepo:   feedRepo,
		Logger:     logger,
	}

	mux := httpRouter.Mux()

	err = http.ListenAndServe(configuration.ApiServer.ApplicationUrl, mux)
	if err != nil {
		panic(err)
	}
}

type Router struct {
	Subscriber message.Subscriber
	Publisher  Publisher
	TweetRepo  repositories.TweetRepository
	FeedRepo   repositories.FeedRepository
	Logger     watermill.LoggerAdapter
}

func (router Router) Mux() *chi.Mux {
	r := chi.NewRouter()

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

	tweetStream := TweetStreamAdapter{repo: router.TweetRepo, logger: router.Logger}
	feedStream := FeedStreamAdapter{repo: router.FeedRepo, logger: router.Logger}
	allTweetsStream := AllTweetsStreamAdapter{repo: router.TweetRepo}
	allFeedsStream := AllFeedsStreamAdapter{repo: router.FeedRepo, logger: router.Logger}

	tweetHandler := sseRouter.AddHandler(messaging.TweetCreatedTopic, tweetStream)
	_ = sseRouter.AddHandler(messaging.FeedUpdatedTopic, feedStream)
	allTweetsHandler := sseRouter.AddHandler(messaging.TweetUpdatedTopic, allTweetsStream)
	_ = sseRouter.AddHandler(messaging.FeedUpdatedTopic, allFeedsStream)

	r.Route("/api", func(r chi.Router) {
		r.Post("/tweets", router.CreateTweet)
		r.Get("/tweets", allTweetsHandler)
		r.Get("/tweets/{tweetId}", tweetHandler)
		r.Delete("/tweets/{tweetId}", router.DeleteTweet)
	})

	go func() {
		err = sseRouter.Run(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	<-sseRouter.Running()

	return r
}

func (router Router) CreateTweet(w http.ResponseWriter, r *http.Request) {
	var newTweet models.Tweet
	err := render.Decode(r, &newTweet)
	if err != nil {
		logAndWriteError(router.Logger, w, err)
		return
	}

	created := router.TweetRepo.CreateTweet(newTweet)
	if created == nil {
		return
	}

	event := messaging.TweetCreated{
		Tweet:      newTweet,
		OccurredAt: time.Now().UTC(),
	}

	err = router.Publisher.Publish(messaging.TweetCreatedTopic, event)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(201)
}

func (router Router) DeleteTweet(w http.ResponseWriter, r *http.Request) {
	tweetId := chi.URLParam(r, "tweetId")
	var deleted = router.TweetRepo.DeleteTweet(tweetId)

	if !deleted {
		w.WriteHeader(404)
	}

	w.WriteHeader(204)
}

func logAndWriteError(logger watermill.LoggerAdapter, w http.ResponseWriter, err error) {
	logger.Error("Error", err, nil)
	w.WriteHeader(500)
}

package api

import (
	"context"
	"net/http"
	"time"
	"twitter-clone/internal/messaging"
	"twitter-clone/internal/models"
	"twitter-clone/internal/repositories"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	watermillHTTP "github.com/ThreeDotsLabs/watermill-http/pkg/http"
)

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
	allFeedsStream := AllFeedsStreamAdapter{repo: router.FeedRepo, logger: router.Logger}

	tweetHandler := sseRouter.AddHandler(messaging.TweetCreatedTopic, tweetStream)
	_ = sseRouter.AddHandler(messaging.FeedUpdatedTopic, feedStream)
	_ = sseRouter.AddHandler(messaging.FeedUpdatedTopic, allFeedsStream)

	r.Route("/api", func(r chi.Router) {
		r.Post("/tweets", router.CreateTweet)
		r.Get("/tweets/{tweetId}", tweetHandler)
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
		return
	}
}

func logAndWriteError(logger watermill.LoggerAdapter, w http.ResponseWriter, err error) {
	logger.Error("Error", err, nil)
	w.WriteHeader(500)
}

package messaging

import (
	"context"
	"twitter-clone/internal/config"
	repositories "twitter-clone/internal/repositories/feed"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

type PubSubMessageHandler struct {
}

func (n *PubSubMessageHandler) SetupMessageRouter(
	configuration config.Configuration,
	feedRepo repositories.FeedRepository,
	logger watermill.LoggerAdapter,
) (message.Publisher, message.Subscriber, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, nil, err
	}
	router.AddMiddleware(middleware.Recoverer)

	// Google Pub/Sub Publisher setup
	pub, err := googlecloud.NewPublisher(googlecloud.PublisherConfig{
		ProjectID: configuration.ProjectId,
		Marshaler: googlecloud.DefaultMarshalerUnmarshaler{},
	}, logger)
	if err != nil {
		return nil, nil, err
	}

	// Google Pub/Sub Subscriber setup
	sub, err := googlecloud.NewSubscriber(googlecloud.SubscriberConfig{
		ProjectID: configuration.ProjectId,
		GenerateSubscriptionName: func(topic string) string {
			return topic + "-sse"
		},
	}, logger)
	if err != nil {
		return nil, nil, err
	}

	routerSub, err := googlecloud.NewSubscriber(googlecloud.SubscriberConfig{
		ProjectID: configuration.ProjectId,
		GenerateSubscriptionName: func(topic string) string {
			return topic + "-router"
		},
	}, logger)
	if err != nil {
		return nil, nil, err
	}

	// Add handler to process incoming messages
	router.AddHandler(
		UpdateFeedsOnNewTweetCreated,
		TweetCreatedTopic,
		routerSub,
		FeedUpdatedTopic,
		pub,
		func(msg *message.Message) (messages []*message.Message, err error) {
			return TweetCreatedHandler(msg, feedRepo, logger)
		},
	)

	router.AddHandler(
		UpdateFeedsOnTweetDeleted,
		TweetDeletedTopic,
		routerSub,
		FeedUpdatedTopic,
		pub,
		func(msg *message.Message) (messages []*message.Message, err error) {
			return TweetDeletedHandler(msg, feedRepo, logger)
		},
	)

	go func() {
		err = router.Run(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	<-router.Running()

	return pub, sub, nil
}

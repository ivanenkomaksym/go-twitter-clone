package messaging

import (
	"context"
	"twitter-clone/internal/config"
	repositories "twitter-clone/internal/repositories/feed"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/nats-io/stan.go"
)

type NATSMessageHandler struct {
}

func (n *NATSMessageHandler) SetupMessageRouter(
	configuration config.Configuration,
	feedRepo repositories.FeedRepository,
	logger watermill.LoggerAdapter,
) (message.Publisher, message.Subscriber, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, nil, err
	}
	router.AddMiddleware(middleware.Recoverer)

	natsURL := stan.NatsURL(configuration.NATSUrl)
	pub, err := nats.NewStreamingPublisher(nats.StreamingPublisherConfig{
		ClusterID:   "test-cluster",
		ClientID:    "publisher",
		StanOptions: []stan.Option{natsURL},
		Marshaler:   nats.GobMarshaler{},
	}, logger)
	if err != nil {
		return nil, nil, err
	}

	sub, err := nats.NewStreamingSubscriber(nats.StreamingSubscriberConfig{
		ClusterID:   "test-cluster",
		ClientID:    "subscriber",
		StanOptions: []stan.Option{natsURL},
		Unmarshaler: nats.GobMarshaler{},
	}, logger)
	if err != nil {
		return nil, nil, err
	}

	router.AddHandler(
		HandlerName,
		TweetCreatedTopic,
		sub,
		FeedUpdatedTopic,
		pub,
		func(msg *message.Message) (messages []*message.Message, err error) {
			return TweetCreatedHandler(msg, feedRepo, logger)
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

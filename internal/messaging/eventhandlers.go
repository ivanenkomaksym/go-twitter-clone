package messaging

import (
	"context"
	"encoding/json"
	"time"
	"twitter-clone/internal/repositories"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/nats-io/stan.go"
)

const (
	TweetCreatedTopic = "tweet-created"
	TweetUpdatedTopic = "tweet-updated"
	FeedUpdatedTopic  = "feed-updated"
)

func SetupMessageRouter(
	feedRepo repositories.FeedRepository,
	logger watermill.LoggerAdapter,
) (message.Publisher, message.Subscriber, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, nil, err
	}
	router.AddMiddleware(middleware.Recoverer)

	natsURL := stan.NatsURL("nats://localhost:4222")
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
		"update-feeds-on-tweet-created",
		TweetCreatedTopic,
		sub,
		FeedUpdatedTopic,
		pub,
		func(msg *message.Message) (messages []*message.Message, err error) {
			defer func() {
				if err == nil {
					logger.Info("Successfully updated feeds on new tweet created", nil)
				} else {
					logger.Error("Error while updating feeds on new twet created", err, nil)
				}
			}()

			event := TweetCreated{}
			err = json.Unmarshal(msg.Payload, &event)
			if err != nil {
				return nil, err
			}

			logger.Info("Adding tweet", watermill.LogFields{"post": event.Tweet})

			if len(event.Tweet.Tags) > 0 {
				for _, tag := range event.Tweet.Tags {
					logger.Info("Adding tag", watermill.LogFields{"tag": tag})
					err = feedRepo.CreateFeed(tag)
					if err != nil {
						return nil, err
					}

				}

				err = feedRepo.AppendTweet(event.Tweet)
				if err != nil {
					return nil, err
				}
			}

			return createFeedUpdatedEvents(event.Tweet.Tags)
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

func createFeedUpdatedEvents(tags []string) ([]*message.Message, error) {
	var messages []*message.Message

	for _, tag := range tags {
		event := FeedUpdated{
			Name:       tag,
			OccurredAt: time.Now().UTC(),
		}

		payload, err := json.Marshal(event)
		if err != nil {
			return nil, err
		}

		msg := message.NewMessage(watermill.NewUUID(), payload)

		messages = append(messages, msg)
	}

	return messages, nil
}

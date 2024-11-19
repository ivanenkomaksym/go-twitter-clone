package messaging

import (
	"encoding/json"
	"time"
	repositories "twitter-clone/internal/repositories/feed"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

func TweetCreatedHandler(
	msg *message.Message,
	feedRepo repositories.FeedRepository,
	logger watermill.LoggerAdapter,
) (messages []*message.Message, err error) {

	defer func() {
		if err == nil {
			logger.Info("Successfully updated feeds on new tweet created", nil)
		} else {
			logger.Error("Error while updating feeds on new tweet created", err, nil)
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

	return CreateFeedUpdatedEvents(event.Tweet.Tags)
}

func TweetDeletedHandler(
	msg *message.Message,
	feedRepo repositories.FeedRepository,
	logger watermill.LoggerAdapter,
) (messages []*message.Message, err error) {

	defer func() {
		if err == nil {
			logger.Info("Successfully updated feeds on tweet deleted", nil)
		} else {
			logger.Error("Error while updating feeds on tweet deleted", err, nil)
		}
	}()

	event := TweetDeleted{}
	err = json.Unmarshal(msg.Payload, &event)
	if err != nil {
		return nil, err
	}

	logger.Info("Deleting tweet", watermill.LogFields{"post": event.DeletedTweet})

	feedRepo.DeleteTweet(event.DeletedTweet)
	// TODO: handle error

	return CreateFeedUpdatedEvents(event.DeletedTweet.Tags)
}

func CreateFeedUpdatedEvents(tags []string) ([]*message.Message, error) {
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

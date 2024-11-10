package repositories

import (
	"errors"
	"twitter-clone/internal/config"
	feedrepo "twitter-clone/internal/repositories/feed"
	tweetrepo "twitter-clone/internal/repositories/tweet"
)

func CreateTweetRepository(configuration config.Configuration) (tweetrepo.TweetRepository, error) {
	switch configuration.Mode {
	case config.InMemory:
		return &tweetrepo.InMemoryTweetRepository{}, nil
	case config.Persistent:
		return tweetrepo.NewPersistentTweetRepository(configuration)
	case config.Cloud:
		return tweetrepo.NewFirestoreTweetRepository(configuration)
	default:
		return nil, errors.New("unknown mode")
	}
}

func CreateFeedRepository(configuration config.Configuration) (feedrepo.FeedRepository, error) {
	switch configuration.Mode {
	case config.InMemory:
		return &feedrepo.InMemoryFeedRepository{}, nil
	case config.Persistent:
		return feedrepo.NewPersistentFeedRepository(configuration)
	case config.Cloud:
		return feedrepo.NewFirestoreFeedRepository(configuration)
	default:
		return nil, errors.New("unknown mode")
	}
}

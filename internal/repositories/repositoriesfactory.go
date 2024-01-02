package repositories

import (
	"errors"
	"twitter-clone/internal/config"
)

func CreateTweetRepository(configuration config.Configuration) (TweetRepository, error) {
	switch configuration.Mode {
	case config.InMemory:
		return &InMemoryTweetRepository{}, nil
	case config.Persistent:
		return NewPersistentTweetRepository(configuration)
	default:
		return nil, errors.New("unknown mode")
	}
}

func CreateFeedRepository(configuration config.Configuration) (FeedRepository, error) {
	switch configuration.Mode {
	case config.InMemory:
		return &InMemoryFeedRepository{}, nil
	case config.Persistent:
		return nil, errors.New("persistent mode not implemented")
	default:
		return nil, errors.New("unknown mode")
	}
}

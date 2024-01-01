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
		return nil, errors.New("persistent mode not implemented")
	default:
		return nil, errors.New("unknown mode")
	}
}

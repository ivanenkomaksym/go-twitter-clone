package repositories

import (
	"errors"
	"twitter-clone/internal/config"
)

func CreateTweetRepository(configuration config.Configuration) (TweetRepository, error) {
	switch configuration.Mode {
	case config.InMemory:
		return nil, errors.New("inmemory mode not implemented")
	case config.Persistent:
		return nil, errors.New("persistent mode not implemented")
	default:
		return nil, errors.New("unknown mode")
	}
}

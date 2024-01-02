package repositories

import "twitter-clone/internal/models"

type FeedRepository interface {
	CreateFeed(name string) error
	GetFeeds() ([]models.Feed, error)
	GetFeedByName(name string) (*models.Feed, error)
	AppendTweet(tweet models.Tweet) error
}

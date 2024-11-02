package repositories

import "twitter-clone/internal/models"

type TweetRepository interface {
	CreateTweet(tweet models.CreateTweetRequest) *models.Tweet
	GetTweets() []models.Tweet
	GetTweetById(id string) *models.Tweet
	DeleteTweet(id string) bool
}

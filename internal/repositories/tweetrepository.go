package repositories

import "twitter-clone/internal/models"

type TweetRepository interface {
	createTweet(tweet models.Tweet) models.Tweet
	getTweets() []models.Tweet
	getTweetById(id string) models.Tweet
	deleteTweet(id string) bool
}

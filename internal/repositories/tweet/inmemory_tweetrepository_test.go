package repositories_test

import (
	"testing"
	repositories "twitter-clone/internal/repositories/tweet"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryTweetRepository(t *testing.T) {
	// Initialize the repository
	repo := repositories.InMemoryTweetRepository{}

	tweet := repositories.TestCreateTweetRequest
	user := repositories.TestUser

	// Test CreateTweet
	createdTweet := repo.CreateTweet(tweet, user)
	assert.NotNil(t, createdTweet, "CreateTweet should return the created tweet")
	tweetID := createdTweet.ID

	// Attempt to create a tweet with the same ID (should return nil)
	duplicateTweet := repo.CreateTweet(tweet, user)
	assert.Nil(t, duplicateTweet, "CreateTweet should return nil for duplicate tweet")

	// Test GetTweets
	allTweets := repo.GetTweets()
	assert.Len(t, allTweets, 1, "GetTweets should return a single tweet")

	// Test GetTweetById
	foundTweet := repo.GetTweetById(tweetID)
	assert.NotNil(t, foundTweet, "GetTweetById should find the tweet")
	assert.Equal(t, tweetID, foundTweet.ID, "Found tweet should have the same ID")

	// Attempt to get a non-existing tweet by ID (should return nil)
	nonExistingTweet := repo.GetTweetById("non-existing-id")
	assert.Nil(t, nonExistingTweet, "GetTweetById should return nil for non-existing tweet")

	// Test DeleteTweet
	deleted := repo.DeleteTweet(tweetID)
	assert.True(t, deleted, "DeleteTweet should return true for successful deletion")

	// Attempt to delete the same tweet again (should return false)
	nonExistingDelete := repo.DeleteTweet(tweetID)
	assert.False(t, nonExistingDelete, "DeleteTweet should return false for non-existing tweet deletion")

	// Ensure no tweets are left after deletion
	remainingTweets := repo.GetTweets()
	assert.Len(t, remainingTweets, 0, "GetTweets should return no tweets after deletion")
}

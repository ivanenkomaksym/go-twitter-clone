package repositories_test

import (
	"fmt"
	"testing"
	"time"
	"twitter-clone/internal/config"
	"twitter-clone/internal/models"
	"twitter-clone/internal/repositories"

	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
)

var repo repositories.TweetRepository

func TestMain(m *testing.M) {
	// Setup the test database
	configuration := config.Configuration{
		Mode: config.Persistent,
		TweetsStorage: config.TweetsStorage{
			ConnectionString: "myuser:mypassword@tcp(127.0.0.1:3306)",
			DatabaseName:     "TweetsDb",
		},
	}

	r, err := repositories.CreateTweetRepository(configuration)
	if err != nil {
		fmt.Println("Failed to create tweet repository: ", err)
		return
	}

	repo = r

	// Run tests
	m.Run()
}

func TestCreateTweet(t *testing.T) {
	// Create a tweet
	tweet := models.Tweet{
		ID:        "abc",
		Title:     "title",
		Content:   "content",
		Author:    "author",
		Tags:      []string{"tag1"},
		CreatedAt: models.MySQLTimestamp{Time: time.Now()},
	}

	// Test the CreateTweet method
	createdTweet := repo.CreateTweet(tweet)
	assert.NotNil(t, createdTweet, "CreateTweet should return the created tweet")
	assert.Equal(t, tweet.ID, createdTweet.ID, "Created tweet should have the same ID")

	// Attempt to create a tweet with the same ID (should return nil)
	duplicateTweet := repo.CreateTweet(tweet)
	assert.Nil(t, duplicateTweet, "CreateTweet should return nil for duplicate tweet")

	// Test GetTweets
	allTweets := repo.GetTweets()
	assert.Len(t, allTweets, 1, "GetTweets should return a single tweet")

	// Test GetTweetById
	foundTweet := repo.GetTweetById(tweet.ID)
	assert.NotNil(t, foundTweet, "GetTweetById should find the tweet")
	assert.Equal(t, tweet.ID, foundTweet.ID, "Found tweet should have the same ID")

	// Attempt to get a non-existing tweet by ID (should return nil)
	nonExistingTweet := repo.GetTweetById("non-existing-id")
	assert.Nil(t, nonExistingTweet, "GetTweetById should return nil for non-existing tweet")

	// Test DeleteTweet
	deleted := repo.DeleteTweet(tweet.ID)
	assert.True(t, deleted, "DeleteTweet should return true for successful deletion")

	// Attempt to delete the same tweet again (should return false)
	nonExistingDelete := repo.DeleteTweet(tweet.ID)
	assert.False(t, nonExistingDelete, "DeleteTweet should return false for non-existing tweet deletion")

	// Ensure no tweets are left after deletion
	remainingTweets := repo.GetTweets()
	assert.Len(t, remainingTweets, 0, "GetTweets should return no tweets after deletion")
}

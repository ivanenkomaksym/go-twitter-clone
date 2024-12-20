package repositories_test

import (
	"fmt"
	"testing"
	"twitter-clone/internal/config"
	repositories "twitter-clone/internal/repositories"
	tweetrepo "twitter-clone/internal/repositories/tweet"

	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
)

func setupTweetRepo() tweetrepo.TweetRepository {
	// Setup the test database
	// TODO: Move out connection string and database name to be read from settings or env vars
	configuration := config.Configuration{
		Mode: config.Persistent,
		TweetsStorage: config.TweetsStorage{
			ConnectionString: "myuser:mypassword@tcp(127.0.0.1:3306)",
			DatabaseName:     "Tests_TweetsDb",
		},
	}

	repo, err := repositories.CreateTweetRepository(configuration)
	if err != nil {
		fmt.Println("Failed to create tweet repository: ", err)
		return nil
	}

	return repo
}

func TestCreateTweet(t *testing.T) {
	repo := setupTweetRepo()

	tweet := tweetrepo.TestCreateTweetRequest
	user := tweetrepo.TestUser

	// Test the CreateTweet method
	createdTweet := repo.CreateTweet(tweet, user)
	assert.NotNil(t, createdTweet, "CreateTweet should return the created tweet")
	tweetId := createdTweet.ID

	// Test GetTweets
	allTweets := repo.GetTweets()
	assert.Len(t, allTweets, 1, "GetTweets should return a single tweet")

	// Test GetTweetById
	foundTweet := repo.GetTweetById(tweetId)
	assert.NotNil(t, foundTweet, "GetTweetById should find the tweet")
	assert.Equal(t, tweetId, foundTweet.ID, "Found tweet should have the same ID")

	// Attempt to get a non-existing tweet by ID (should return nil)
	nonExistingTweet := repo.GetTweetById("non-existing-id")
	assert.Nil(t, nonExistingTweet, "GetTweetById should return nil for non-existing tweet")

	// Test DeleteTweet
	deleted := repo.DeleteTweet(tweetId)
	assert.True(t, deleted, "DeleteTweet should return true for successful deletion")

	// Attempt to delete the same tweet again (should return false)
	nonExistingDelete := repo.DeleteTweet(tweetId)
	assert.False(t, nonExistingDelete, "DeleteTweet should return false for non-existing tweet deletion")

	// Ensure no tweets are left after deletion
	remainingTweets := repo.GetTweets()
	assert.Len(t, remainingTweets, 0, "GetTweets should return no tweets after deletion")
}

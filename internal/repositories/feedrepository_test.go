package repositories_test

import (
	"fmt"
	"testing"
	"time"
	"twitter-clone/internal/config"
	"twitter-clone/internal/models"
	"twitter-clone/internal/repositories"

	"github.com/stretchr/testify/assert"
)

func setupFeedRepo() repositories.FeedRepository {
	// Setup the test database
	configuration := config.Configuration{
		Mode: config.Persistent,
		FeedsStorage: config.FeedsStorage{
			ConnectionString: "mongodb://localhost:27017",
			DatabaseName:     "Tests_FeedsDb",
			CollectionName:   "Feeds",
		},
	}

	repo, err := repositories.CreateFeedRepository(configuration)
	if err != nil {
		fmt.Println("Failed to create feed repository: ", err)
		return nil
	}

	return repo
}

func Test_CreateFeed(t *testing.T) {
	feedRepo := setupFeedRepo()
	name := "testFeed"

	// Create a new feed
	err := feedRepo.CreateFeed(name)
	assert.NoError(t, err, "Error creating feed")

	// Try to create the same feed again, it should not return an error
	err = feedRepo.CreateFeed(name)
	assert.NoError(t, err, "Error creating feed")

	// Verify that the feed exists in the repository
	feed, err := feedRepo.GetFeedByName(name)
	assert.NoError(t, err, "Error getting feed by name")
	assert.NotNil(t, feed, "Expected feed to exist, but it doesn't.")

	// Test DeleteFeed
	deleted := feedRepo.DeleteFeed(name)
	assert.True(t, deleted, "DeleteFeed should return true for successful deletion")
}

func Test_GetFeeds(t *testing.T) {
	feedRepo := setupFeedRepo()
	name := "testFeed"

	// Get feeds from an empty repository
	feeds, err := feedRepo.GetFeeds()
	assert.NoError(t, err, "Error getting feeds")
	assert.Empty(t, feeds, "Expected no feeds, got %d feeds", len(feeds))

	// Create a new feed
	err = feedRepo.CreateFeed(name)
	assert.NoError(t, err, "Error creating feed")

	// Get feeds after creating one
	feeds, err = feedRepo.GetFeeds()
	assert.NoError(t, err, "Error getting feeds")
	assert.Len(t, feeds, 1, "Expected 1 feed, got %d feeds", len(feeds))

	deleted := feedRepo.DeleteFeed(name)
	assert.True(t, deleted, "DeleteFeed should return true for successful deletion")
}

func TestGetFeedByName(t *testing.T) {
	feedRepo := setupFeedRepo()

	feedName := "TechNews"
	feedRepo.CreateFeed(feedName)

	feed, err := feedRepo.GetFeedByName(feedName)

	// Assert that the feed can be retrieved by its name without errors
	assert.NoError(t, err, "Should retrieve the feed without errors")
	assert.Equal(t, feedName, feed.Name, "The feed name should match the expected name")

	deleted := feedRepo.DeleteFeed(feedName)
	assert.True(t, deleted, "DeleteFeed should return true for successful deletion")
}

func TestAppendTweet(t *testing.T) {
	feedRepo := setupFeedRepo()

	feedName := "TechNews"
	feedRepo.CreateFeed(feedName)

	expectedTweet := models.Tweet{
		ID:        "abc",
		Title:     "title",
		Content:   "content",
		Author:    "author",
		Tags:      []string{"TechNews"},
		CreatedAt: models.MySQLTimestamp{Time: time.Now()},
	}

	err := feedRepo.AppendTweet(expectedTweet)

	// Assert that the tweet is appended without errors
	assert.NoError(t, err, "Tweet should be appended without errors")

	// Optionally, verify that the tweet was added to the feed
	feed, err := feedRepo.GetFeedByName(feedName)
	assert.NoError(t, err, "Feed should be retrieved without errors")

	// Assuming Feed has a Tweets field that is a slice of tweets
	assert.Equal(t, 1, len(feed.Tweets))
	actualTweet := feed.Tweets[0]

	assert.Equal(t, expectedTweet.ID, actualTweet.ID)
	assert.Equal(t, expectedTweet.Title, actualTweet.Title)
	assert.Equal(t, expectedTweet.Content, actualTweet.Content)
	assert.Equal(t, expectedTweet.Author, actualTweet.Author)
	assert.Equal(t, expectedTweet.Tags, actualTweet.Tags)

	// Compare only the Unix timestamp part of the time
	assert.Equal(t, expectedTweet.CreatedAt.Unix(), actualTweet.CreatedAt.Unix(), "CreatedAt times should be equal")

	deleted := feedRepo.DeleteFeed(feedName)
	assert.True(t, deleted, "DeleteFeed should return true for successful deletion")
}

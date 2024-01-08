package repositories_test

import (
	"testing"
	"twitter-clone/internal/repositories"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryFeedRepository_CreateFeed(t *testing.T) {
	repo := repositories.InMemoryFeedRepository{}
	name := "testFeed"

	// Create a new feed
	err := repo.CreateFeed(name)
	assert.NoError(t, err, "Error creating feed")

	// Try to create the same feed again, it should not return an error
	err = repo.CreateFeed(name)
	assert.NoError(t, err, "Error creating feed")

	// Verify that the feed exists in the repository
	feed, err := repo.GetFeedByName(name)
	assert.NoError(t, err, "Error getting feed by name")
	assert.NotNil(t, feed, "Expected feed to exist, but it doesn't.")
}

func TestInMemoryFeedRepository_GetFeeds(t *testing.T) {
	repo := repositories.InMemoryFeedRepository{}

	// Get feeds from an empty repository
	feeds, err := repo.GetFeeds()
	assert.NoError(t, err, "Error getting feeds")
	assert.Empty(t, feeds, "Expected no feeds, got %d feeds", len(feeds))

	// Create a new feed
	err = repo.CreateFeed("testFeed")
	assert.NoError(t, err, "Error creating feed")

	// Get feeds after creating one
	feeds, err = repo.GetFeeds()
	assert.NoError(t, err, "Error getting feeds")
	assert.Len(t, feeds, 1, "Expected 1 feed, got %d feeds", len(feeds))
}

// Add more test functions for other methods...

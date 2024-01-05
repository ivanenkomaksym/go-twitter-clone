package repositories

import (
	"slices"
	"twitter-clone/internal/models"
)

type InMemoryFeedRepository struct {
	feeds []models.Feed
}

func (repo *InMemoryFeedRepository) CreateFeed(name string) error {
	idx := slices.IndexFunc(repo.feeds, func(f models.Feed) bool { return f.Name == name })
	if idx != -1 {
		return nil
	}

	feed := models.Feed{
		Name:   name,
		Tweets: []models.Tweet{},
	}

	repo.feeds = append(repo.feeds, feed)
	return nil
}

func (repo *InMemoryFeedRepository) GetFeeds() ([]models.Feed, error) {
	return repo.feeds, nil
}

func (repo *InMemoryFeedRepository) GetFeedByName(name string) (*models.Feed, error) {
	idx := slices.IndexFunc(repo.feeds, func(f models.Feed) bool { return f.Name == name })
	if idx == -1 {
		return nil, nil
	}

	return &repo.feeds[idx], nil
}

func (repo *InMemoryFeedRepository) AppendTweet(tweet models.Tweet) error {
	if len(tweet.Tags) == 0 {
		return nil
	}

	for idx, _ := range repo.feeds {
		// Check if the feed should get updated based on new tweet tags
		if containsTag(tweet.Tags, repo.feeds[idx].Name) {
			// Check if the tweet is not already present in the feed
			if !containsTweet(repo.feeds[idx].Tweets, tweet.ID) {
				// Add the tweet to the feed
				repo.feeds[idx].Tweets = append(repo.feeds[idx].Tweets, tweet)
			}
		}
	}

	return nil
}

// containsTag checks if a given tag is present in the tags slice
func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

// containsTweet checks if a given tweet ID is present in the tweets slice
func containsTweet(tweets []models.Tweet, tweetId string) bool {
	for _, t := range tweets {
		if t.ID == tweetId {
			return true
		}
	}
	return false
}

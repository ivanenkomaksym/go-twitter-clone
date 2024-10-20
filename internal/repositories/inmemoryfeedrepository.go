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
		if ContainsTag(tweet.Tags, repo.feeds[idx].Name) {
			// Check if the tweet is not already present in the feed
			if !ContainsTweet(repo.feeds[idx].Tweets, tweet.ID) {
				// Add the tweet to the feed
				repo.feeds[idx].Tweets = append(repo.feeds[idx].Tweets, tweet)
			}
		}
	}

	return nil
}

func (repo *InMemoryFeedRepository) DeleteFeed(name string) bool {
	removed := false
	for i, feed := range repo.feeds {
		if feed.Name == name {
			repo.feeds = append(repo.feeds[:i], repo.feeds[i+1:]...)
			removed = true
			break
		}
	}

	return removed
}

package repositories

import (
	"slices"
	"twitter-clone/internal/models"
)

type InMemoryTweetRepository struct {
	tweets []models.Tweet
}

func (repo *InMemoryTweetRepository) CreateTweet(createTweetRequest models.CreateTweetRequest) *models.Tweet {
	tweet := CreateNewTweet(createTweetRequest)
	idx := slices.IndexFunc(repo.tweets, func(t models.Tweet) bool { return t.ID == tweet.ID })
	if idx != -1 {
		return nil
	}

	repo.tweets = append(repo.tweets, tweet)
	return &tweet
}

func (repo *InMemoryTweetRepository) GetTweets() []models.Tweet {
	return repo.tweets
}

func (repo *InMemoryTweetRepository) GetTweetById(id string) *models.Tweet {
	idx := slices.IndexFunc(repo.tweets, func(t models.Tweet) bool { return t.ID == id })
	if idx == -1 {
		return nil
	}

	return &repo.tweets[idx]
}

func (repo *InMemoryTweetRepository) DeleteTweet(id string) bool {
	idx := slices.IndexFunc(repo.tweets, func(t models.Tweet) bool { return t.ID == id })
	if idx == -1 {
		return false
	}

	repo.tweets[idx] = repo.tweets[len(repo.tweets)-1]
	repo.tweets = repo.tweets[:len(repo.tweets)-1]

	return true
}

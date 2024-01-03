package api

import (
	"encoding/json"
	"net/http"
	"twitter-clone/internal/messaging"
	"twitter-clone/internal/repositories"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi/v5"
)

type FeedStreamAdapter struct {
	repo   repositories.FeedRepository
	logger watermill.LoggerAdapter
}

func (f FeedStreamAdapter) GetResponse(w http.ResponseWriter, r *http.Request) (response interface{}, ok bool) {
	feedName := chi.URLParam(r, "name")

	feed, err := f.repo.GetFeedByName(feedName)
	if err != nil {
		logAndWriteError(f.logger, w, err)
		return nil, false
	}

	return feed, true
}

func (f FeedStreamAdapter) Validate(r *http.Request, msg *message.Message) (ok bool) {
	feedUpdated := messaging.FeedUpdated{}

	err := json.Unmarshal(msg.Payload, &feedUpdated)
	if err != nil {
		return false
	}

	feedName := chi.URLParam(r, "name")

	return feedUpdated.Name == feedName
}

type TweetStreamAdapter struct {
	repo   repositories.TweetRepository
	logger watermill.LoggerAdapter
}

func (p TweetStreamAdapter) GetResponse(w http.ResponseWriter, r *http.Request) (response interface{}, ok bool) {
	tweetID := chi.URLParam(r, "tweetId")

	tweet := p.repo.GetTweetById(tweetID)
	if tweet == nil {
		return nil, false
	}

	return tweet, true
}

func (p TweetStreamAdapter) Validate(r *http.Request, msg *message.Message) (ok bool) {
	postUpdated := messaging.TweetUpdated{}

	err := json.Unmarshal(msg.Payload, &postUpdated)
	if err != nil {
		return false
	}

	tweetID := chi.URLParam(r, "tweetId")

	return postUpdated.OriginalTweet.ID == tweetID
}

type FeedSummary struct {
	Name   string `json:"name"`
	Tweets int    `json:"tweets"`
}

type AllFeedsResponse struct {
	Feeds []FeedSummary `json:"feeds"`
}

type AllFeedsStreamAdapter struct {
	repo   repositories.FeedRepository
	logger watermill.LoggerAdapter
}

func (f AllFeedsStreamAdapter) GetResponse(w http.ResponseWriter, r *http.Request) (interface{}, bool) {
	feeds, err := f.repo.GetFeeds()
	if err != nil {
		logAndWriteError(f.logger, w, err)
		return nil, false
	}

	response := AllFeedsResponse{
		Feeds: []FeedSummary{},
	}

	for _, f := range feeds {
		response.Feeds = append(response.Feeds, FeedSummary{
			Name:   f.Name,
			Tweets: len(f.Tweets),
		})
	}

	return response, true
}

func (f AllFeedsStreamAdapter) Validate(r *http.Request, msg *message.Message) (ok bool) {
	return true
}

type AllTweetsStreamAdapter struct {
	repo repositories.TweetRepository
}

func (f AllTweetsStreamAdapter) GetResponse(w http.ResponseWriter, r *http.Request) (interface{}, bool) {
	tweets := f.repo.GetTweets()
	return tweets, true
}

func (f AllTweetsStreamAdapter) Validate(r *http.Request, msg *message.Message) (ok bool) {
	return true
}

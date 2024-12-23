package api

import (
	"encoding/json"
	"net/http"
	"twitter-clone/internal/messaging"
	feedrepo "twitter-clone/internal/repositories/feed"
	repositories "twitter-clone/internal/repositories/tweet"
	tweetrepo "twitter-clone/internal/repositories/tweet"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi/v5"
)

type FeedStreamAdapter struct {
	repo   feedrepo.FeedRepository
	logger watermill.LoggerAdapter
}

func (adapter FeedStreamAdapter) GetResponse(w http.ResponseWriter, r *http.Request) (response interface{}, ok bool) {
	feedName := chi.URLParam(r, "name")

	feed, err := adapter.repo.GetFeedByName(feedName)
	if err != nil {
		logAndWriteError(adapter.logger, w, err)
		return nil, false
	}

	if feed == nil {
		w.WriteHeader(404)
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
	repo   tweetrepo.TweetRepository
	logger watermill.LoggerAdapter
}

func (adapter TweetStreamAdapter) GetResponse(w http.ResponseWriter, r *http.Request) (response interface{}, ok bool) {
	tweetID := chi.URLParam(r, "tweetId")

	tweet := adapter.repo.GetTweetById(tweetID)
	if tweet == nil {
		w.WriteHeader(404)
		return nil, false
	}

	return tweet, true
}

func (adapter TweetStreamAdapter) Validate(r *http.Request, msg *message.Message) (ok bool) {
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
	repo   feedrepo.FeedRepository
	logger watermill.LoggerAdapter
}

func (adapter AllFeedsStreamAdapter) GetResponse(w http.ResponseWriter, r *http.Request) (interface{}, bool) {
	feeds, err := adapter.repo.GetFeeds()
	if err != nil {
		logAndWriteError(adapter.logger, w, err)
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

func (adapter AllTweetsStreamAdapter) GetResponse(w http.ResponseWriter, r *http.Request) (interface{}, bool) {
	tweets := adapter.repo.GetTweets()
	return tweets, true
}

func (f AllTweetsStreamAdapter) Validate(r *http.Request, msg *message.Message) (ok bool) {
	return true
}

package messaging

import (
	"time"
	"twitter-clone/internal/models"
)

const (
	UpdateFeedsOnNewTweetCreated = "update-feeds-on-tweet-created"
	UpdateFeedsOnTweetDeleted    = "update-feeds-on-tweet-deleted"
	TweetCreatedTopic            = "tweet-created"
	TweetDeletedTopic            = "tweet-deleted"
	TweetUpdatedTopic            = "tweet-updated"
	FeedUpdatedTopic             = "feed-updated"
)

type TweetCreated struct {
	Tweet models.Tweet `json:"tweet"`

	OccurredAt time.Time `json:"occurred_at"`
}

type TweetDeleted struct {
	DeletedTweet models.Tweet `json:"original_tweet"`

	OccurredAt time.Time `json:"occurred_at"`
}

type TweetUpdated struct {
	OriginalTweet models.Tweet `json:"original_tweet"`
	NewTweet      models.Tweet `json:"new_tweet"`

	OccurredAt time.Time `json:"occurred_at"`
}

type FeedUpdated struct {
	Name string `json:"name"`

	OccurredAt time.Time `json:"occurred_at"`
}

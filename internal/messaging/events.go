package messaging

import (
	"time"
	"twitter-clone/internal/models"
)

type TweetCreated struct {
	Tweet models.Tweet `json:"tweet"`

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

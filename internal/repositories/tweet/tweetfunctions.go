package repositories

import (
	"time"
	"twitter-clone/internal/models"

	"github.com/google/uuid"
)

func CreateNewTweet(createTweetRequest models.CreateTweetRequest) models.Tweet {
	return models.Tweet{
		ID:        uuid.NewString(),
		Title:     createTweetRequest.Title,
		Content:   createTweetRequest.Content,
		Tags:      createTweetRequest.Tags,
		CreatedAt: models.MySQLTimestamp{Time: time.Now()},
		User:      models.User{}, // TODO: pass over user info from authentication
	}
}

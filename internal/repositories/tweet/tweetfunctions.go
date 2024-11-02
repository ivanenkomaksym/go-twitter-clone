package repositories

import (
	"time"
	"twitter-clone/internal/models"

	"github.com/google/uuid"
)

func CreateNewTweet(createTweetRequest models.CreateTweetRequest, user models.User) models.Tweet {
	return models.Tweet{
		ID:        uuid.NewString(),
		Title:     createTweetRequest.Title,
		Content:   createTweetRequest.Content,
		Tags:      createTweetRequest.Tags,
		CreatedAt: models.MySQLTimestamp{Time: time.Now()},
		User:      user,
	}
}

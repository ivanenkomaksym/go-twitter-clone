package repositories

import "twitter-clone/internal/models"

var (
	TestCreateTweetRequest = models.CreateTweetRequest{
		Title:   "title",
		Content: "content",
		Tags:    []string{"tag1"},
	}

	TestUser = models.User{
		FirstName: "Alice",
		LastName:  "Liddell",
		Email:     "alice@gmail.com",
		Picture:   "picture.png",
	}
)

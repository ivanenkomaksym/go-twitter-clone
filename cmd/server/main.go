package main

import (
	"fmt"
	"twitter-clone/internal/api"
	"twitter-clone/internal/config"
	"twitter-clone/internal/messaging"
	"twitter-clone/internal/repositories"

	"github.com/ThreeDotsLabs/watermill"
)

func main() {
	configuration := config.ReadConfiguration()

	tweetRepo, err := repositories.CreateTweetRepository(configuration)
	if err != nil {
		fmt.Println("Failed to create tweet repository: ", err)
		return
	}

	feedRepo, err := repositories.CreateFeedRepository(configuration)
	if err != nil {
		fmt.Println("Failed to create feed repository: ", err)
		return
	}

	logger := watermill.NewStdLogger(false, false)

	_, _, err = messaging.SetupMessageRouter(feedRepo, logger)
	if err != nil {
		panic(err)
	}

	api.StartHttpServer(configuration, tweetRepo)
}

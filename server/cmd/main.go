package main

import (
	"fmt"
	"twitter-clone/internal/api"
	"twitter-clone/internal/authn"
	"twitter-clone/internal/config"
	"twitter-clone/internal/messaging"
	"twitter-clone/internal/repositories"
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

	messageHandler, err := messaging.CreateMessageHandler(configuration)
	if err != nil {
		fmt.Println("Failed to create message handler: ", err)
		return
	}

	authenticationValidator := authn.AuthenticationValidator{
		Authentication: configuration.Authentication,
	}

	api.StartRouter(configuration, tweetRepo, feedRepo, messageHandler, authenticationValidator)
}

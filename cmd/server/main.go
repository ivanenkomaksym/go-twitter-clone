package main

import (
	"context"
	"fmt"
	"twitter-clone/internal/api"
	"twitter-clone/internal/config"
	"twitter-clone/internal/repositories"
)

func main() {
	configuration := config.ReadConfiguration()

	_, err := repositories.CreateTweetRepository(configuration)
	if err != nil {
		fmt.Println("Failed to create repository: ", err)
		return
	}

	api.StartHttpServer(configuration, context.TODO())
}

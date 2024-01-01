package main

import (
	"fmt"
	"twitter-clone/internal/api"
	"twitter-clone/internal/config"
	"twitter-clone/internal/repositories"
)

func main() {
	configuration := config.ReadConfiguration()

	repo, err := repositories.CreateTweetRepository(configuration)
	if err != nil {
		fmt.Println("Failed to create repository: ", err)
		return
	}

	api.StartHttpServer(configuration, repo)
}

package main

import (
	"context"
	"twitter-clone/internal/api"
	"twitter-clone/internal/config"
)

func main() {
	configuration := config.ReadConfiguration()

	api.StartHttpServer(configuration, context.TODO())
}

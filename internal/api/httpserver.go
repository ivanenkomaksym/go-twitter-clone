package api

import (
	"context"
	"twitter-clone/internal/config"

	"github.com/gin-gonic/gin"
)

func StartHttpServer(configuration config.Configuration, ctx context.Context) {
	router := gin.Default()
	router.POST("/api/tweets", createTweet(ctx))
	router.GET("/api/tweets", getTweets(ctx))
	router.GET("/api/tweets/:tweetid", getTweetByTweetId(ctx))
	router.DELETE("/api/tweets/:tweetid", deleteTweet(ctx))

	if err := router.Run(configuration.ApiServer.ApplicationUrl); err != nil {
		panic(err)
	}
}

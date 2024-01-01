package api

import (
	"twitter-clone/internal/config"
	"twitter-clone/internal/repositories"

	"github.com/gin-gonic/gin"
)

func StartHttpServer(configuration config.Configuration, repo repositories.TweetRepository) {
	router := gin.Default()
	router.POST("/api/tweets", createTweet(repo))
	router.GET("/api/tweets", getTweets(repo))
	router.GET("/api/tweets/:tweetId", getTweetByTweetId(repo))
	router.DELETE("/api/tweets/:tweetId", deleteTweet(repo))

	if err := router.Run(configuration.ApiServer.ApplicationUrl); err != nil {
		panic(err)
	}
}

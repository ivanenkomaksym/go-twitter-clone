package api

import (
	"context"
	"net/http"
	"twitter-clone/internal/models"

	"github.com/gin-gonic/gin"
)

func createTweet(ctx context.Context) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var newTweet models.Tweet

		// Call BindJSON to bind the received JSON to UserPromo.
		if err := c.BindJSON(&newTweet); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request. Cannor serialize Tweet. Details: " + err.Error()})
			return
		}

		c.IndentedJSON(http.StatusCreated, newTweet)
	}

	return gin.HandlerFunc(fn)
}

func getTweets(ctx context.Context) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var results []models.Tweet

		c.IndentedJSON(http.StatusOK, results)
	}

	return gin.HandlerFunc(fn)
}

func getTweetByTweetId(ctx context.Context) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		tweetId := c.Param("tweetId")
		_ = tweetId

		var result models.Tweet

		c.IndentedJSON(http.StatusOK, result)
	}

	return gin.HandlerFunc(fn)
}

func deleteTweet(ctx context.Context) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		tweetId := c.Param("tweetId")
		_ = tweetId

		var result models.Tweet

		c.IndentedJSON(http.StatusNoContent, result)
	}

	return gin.HandlerFunc(fn)
}

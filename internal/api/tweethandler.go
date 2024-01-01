package api

import (
	"net/http"
	"twitter-clone/internal/models"
	"twitter-clone/internal/repositories"

	"github.com/gin-gonic/gin"
)

func createTweet(repo repositories.TweetRepository) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var newTweet models.Tweet

		// Call BindJSON to bind the received JSON to UserPromo.
		if err := c.BindJSON(&newTweet); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request. Cannor serialize Tweet. Details: " + err.Error()})
			return
		}

		var created = repo.CreateTweet(newTweet)

		if created == nil {
			c.IndentedJSON(http.StatusConflict, gin.H{"message": "Conflict. Tweet already exists."})
			return
		}

		c.IndentedJSON(http.StatusCreated, created)
	}

	return gin.HandlerFunc(fn)
}

func getTweets(repo repositories.TweetRepository) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var results = repo.GetTweets()

		c.IndentedJSON(http.StatusOK, results)
	}

	return gin.HandlerFunc(fn)
}

func getTweetByTweetId(repo repositories.TweetRepository) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		tweetId := c.Param("tweetId")
		var found = repo.GetTweetById(tweetId)

		if found == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Tweet doesn't exists."})
			return
		}

		c.IndentedJSON(http.StatusOK, found)
	}

	return gin.HandlerFunc(fn)
}

func deleteTweet(repo repositories.TweetRepository) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		tweetId := c.Param("tweetId")
		var deleted = repo.DeleteTweet(tweetId)

		if !deleted {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Tweet doesn't exists."})
			return
		}

		c.IndentedJSON(http.StatusNoContent, deleted)
	}

	return gin.HandlerFunc(fn)
}

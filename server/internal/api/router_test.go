package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	apimock "twitter-clone/internal/__mocks__/api"
	authnmock "twitter-clone/internal/__mocks__/authn"
	tweetmock "twitter-clone/internal/__mocks__/repositories/tweet"
	"twitter-clone/internal/api"
	"twitter-clone/internal/config"
	"twitter-clone/internal/messaging"
	"twitter-clone/internal/models"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateTweet tests the CreateTweet endpoint for successful tweet creation.
func TestCreateTweet(t *testing.T) {
	// Initialize mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthValidator := authnmock.NewMockIAuthenticationValidator(ctrl)
	mockTweetRepo := tweetmock.NewMockTweetRepository(ctrl)
	mockPublisher := apimock.NewMockIPublisher(ctrl)

	// Mock configuration and logger
	config := config.Configuration{AllowOrigin: "*"}
	logger := watermill.NewStdLogger(false, false)

	// Define a test tweet request
	tweetRequest := models.CreateTweetRequest{
		Content: "Hello, world!",
		Tags:    []string{"test"},
	}

	// Set up the authenticated user
	user := &models.User{IsAnonymous: true}

	// Expected created tweet response
	createdTweet := &models.Tweet{
		ID:        "tweet1",
		Content:   tweetRequest.Content,
		Tags:      tweetRequest.Tags,
		CreatedAt: models.MySQLTimestamp{Time: time.Now()},
	}

	// Configure mocks
	mockAuthValidator.EXPECT().ValidateAuthentication(gomock.Any(), gomock.Any()).Return(user)
	mockTweetRepo.EXPECT().CreateTweet(tweetRequest, *user).Return(createdTweet)
	mockPublisher.EXPECT().Publish(messaging.TweetCreatedTopic, gomock.Any()).Return(nil)

	// Set up the router
	router := api.Router{
		Config:                  config,
		AuthenticationValidator: mockAuthValidator,
		TweetRepo:               mockTweetRepo,
		Publisher:               mockPublisher,
		Logger:                  logger,
	}

	// Create the HTTP request
	body, _ := json.Marshal(tweetRequest)
	req := httptest.NewRequest("POST", "/api/tweets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create the HTTP response recorder
	rr := httptest.NewRecorder()

	// Call the CreateTweet endpoint
	router.CreateTweet(rr, req)

	// Validate the response
	require.Equal(t, http.StatusCreated, rr.Code)
	var responseTweet models.Tweet
	err := json.Unmarshal(rr.Body.Bytes(), &responseTweet)
	require.NoError(t, err)
	assert.Equal(t, createdTweet.Content, responseTweet.Content)
	assert.Equal(t, createdTweet.Tags, responseTweet.Tags)
}

// TestCreateTweetPublishError tests the CreateTweet endpoint when the Publish operation fails.
func TestCreateTweetPublishError(t *testing.T) {
	// Initialize mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthValidator := authnmock.NewMockIAuthenticationValidator(ctrl)
	mockTweetRepo := tweetmock.NewMockTweetRepository(ctrl)
	mockPublisher := apimock.NewMockIPublisher(ctrl)

	// Mock configuration and logger
	config := config.Configuration{AllowOrigin: "*"}
	logger := watermill.NewStdLogger(false, false)

	// Define a test tweet request
	tweetRequest := models.CreateTweetRequest{
		Content: "Hello, world!",
		Tags:    []string{"test"},
	}

	// Set up the authenticated user
	user := &models.User{IsAnonymous: true}

	// Expected created tweet response
	createdTweet := &models.Tweet{
		ID:        "tweet1",
		Content:   tweetRequest.Content,
		Tags:      tweetRequest.Tags,
		CreatedAt: models.MySQLTimestamp{Time: time.Now()},
	}

	// Configure mocks
	mockAuthValidator.EXPECT().ValidateAuthentication(gomock.Any(), gomock.Any()).Return(user)
	mockTweetRepo.EXPECT().CreateTweet(tweetRequest, *user).Return(createdTweet)
	mockPublisher.EXPECT().Publish(messaging.TweetCreatedTopic, gomock.Any()).Return(errors.New("publish error"))

	// Set up the router
	router := api.Router{
		Config:                  config,
		AuthenticationValidator: mockAuthValidator,
		TweetRepo:               mockTweetRepo,
		Publisher:               mockPublisher,
		Logger:                  logger,
	}

	// Create the HTTP request
	body, _ := json.Marshal(tweetRequest)
	req := httptest.NewRequest("POST", "/api/tweets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create the HTTP response recorder
	rr := httptest.NewRecorder()

	// Call the CreateTweet endpoint
	router.CreateTweet(rr, req)

	// Validate the response
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

package repositories

import (
	"context"
	"log"
	"twitter-clone/internal/config"
	"twitter-clone/internal/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type FirestoreTweetRepository struct {
	client *firestore.Client
}

func NewFirestoreTweetRepository(configuration config.Configuration) (*FirestoreTweetRepository, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, configuration.ProjectId)
	if err != nil {
		return nil, err
	}

	return &FirestoreTweetRepository{client: client}, nil
}

func (r *FirestoreTweetRepository) CreateTweet(createTweetRequest models.CreateTweetRequest, user models.User) *models.Tweet {
	tweet := CreateNewTweet(createTweetRequest, user)
	_, err := r.client.Collection("tweets").Doc(tweet.ID).Set(context.Background(), tweet)
	if err != nil {
		log.Printf("Failed to create tweet: %v", err)
		return nil
	}
	return &tweet
}

func (r *FirestoreTweetRepository) GetTweets() []models.Tweet {
	var tweets []models.Tweet
	iter := r.client.Collection("tweets").Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to fetch tweet: %v", err)
			return nil
		}
		var tweet models.Tweet
		if err := doc.DataTo(&tweet); err != nil {
			log.Printf("Failed to decode tweet: %v", err)
			return nil
		}
		tweets = append(tweets, tweet)
	}
	return tweets
}

func (r *FirestoreTweetRepository) GetTweetById(id string) *models.Tweet {
	doc, err := r.client.Collection("tweets").Doc(id).Get(context.Background())
	if err != nil {
		log.Printf("Failed to fetch tweet by ID: %v", err)
		return nil
	}
	var tweet models.Tweet
	if err := doc.DataTo(&tweet); err != nil {
		log.Printf("Failed to decode tweet: %v", err)
		return nil
	}
	return &tweet
}

func (r *FirestoreTweetRepository) DeleteTweet(id string) bool {
	_, err := r.client.Collection("tweets").Doc(id).Delete(context.Background())
	return err == nil
}

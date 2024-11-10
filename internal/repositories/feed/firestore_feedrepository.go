package repositories

import (
	"context"
	"time"
	"twitter-clone/internal/config"
	"twitter-clone/internal/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type FirestoreFeedRepository struct {
	client *firestore.Client
}

func NewFirestoreFeedRepository(configuration config.Configuration) (*FirestoreFeedRepository, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, configuration.ProjectId)
	if err != nil {
		return nil, err
	}

	return &FirestoreFeedRepository{client: client}, nil
}

func (r *FirestoreFeedRepository) CreateFeed(name string) error {
	_, err := r.client.Collection("feeds").Doc(name).Set(context.Background(), map[string]interface{}{
		"name":       name,
		"created_at": time.Now(),
	})
	return err
}

func (r *FirestoreFeedRepository) GetFeeds() ([]models.Feed, error) {
	var feeds []models.Feed
	iter := r.client.Collection("feeds").Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var feed models.Feed
		if err := doc.DataTo(&feed); err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}
	return feeds, nil
}

func (r *FirestoreFeedRepository) GetFeedByName(name string) (*models.Feed, error) {
	doc, err := r.client.Collection("feeds").Doc(name).Get(context.Background())
	if err != nil {
		return nil, err
	}
	var feed models.Feed
	if err := doc.DataTo(&feed); err != nil {
		return nil, err
	}
	return &feed, nil
}

func (r *FirestoreFeedRepository) AppendTweet(tweet models.Tweet) error {
	ctx := context.Background()

	// If there are no tags, there's no feed to append the tweet to.
	if len(tweet.Tags) == 0 {
		return nil
	}

	// Loop over each tag and update the corresponding feed document.
	for _, tag := range tweet.Tags {
		feedDocRef := r.client.Collection("feeds").Doc(tag)

		// Run a Firestore transaction to safely append the tweet.
		err := r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			// Check if the tweet already exists in this feed's tweets.
			tweetsCollection := feedDocRef.Collection("tweets")
			tweetDoc := tweetsCollection.Doc(tweet.ID)
			doc, err := tweetDoc.Get(ctx)
			if err == nil && doc.Exists() {
				// Tweet already exists, no need to add it again.
				return nil
			}

			// If tweet doesn't exist, add it at the beginning of the tweets subcollection.
			_, err = tweetsCollection.Doc(tweet.ID).Set(ctx, tweet)
			return err
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *FirestoreFeedRepository) DeleteFeed(name string) bool {
	_, err := r.client.Collection("feeds").Doc(name).Delete(context.Background())
	return err == nil
}

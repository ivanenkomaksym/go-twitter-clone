package repositories

import (
	"context"
	"time"
	"twitter-clone/internal/config"
	"twitter-clone/internal/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	ctx := context.Background()

	// Reference to the feed document for the given name
	feedDocRef := r.client.Collection("feeds").Doc(name)

	// Check if the document already exists
	docSnapshot, err := feedDocRef.Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		// If the error is not 'NotFound', something else went wrong
		return err
	}

	// If the document already exists, we don't create it again
	if docSnapshot.Exists() {
		return nil
	}

	// Otherwise, create the feed document lazily
	_, err = feedDocRef.Set(ctx, map[string]interface{}{
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

		// Run a Firestore transaction to safely check and append the tweet.
		err := r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			// Check if the feed document exists.
			doc, err := tx.Get(feedDocRef)
			if err != nil {
				if status.Code(err) == codes.NotFound {
					// Document does not exist, so we initialize it with an empty tweets array.
					err := tx.Set(feedDocRef, map[string]interface{}{
						"tweets": []models.Tweet{},
					})
					if err != nil {
						return err
					}
				} else {
					return err // Return any other error.
				}
			}

			var feedData struct {
				Tweets []models.Tweet `firestore:"tweets"`
			}

			err = doc.DataTo(&feedData)
			if err != nil {
				return err
			}

			updatedTweets := append(feedData.Tweets, tweet)

			// Append the tweet to the tweets array using ArrayUnion to avoid duplicates.
			err = tx.Update(feedDocRef, []firestore.Update{
				{
					Path:  "tweets",
					Value: updatedTweets,
				},
			})
			return err
		})

		if err != nil {
			return err // Return if there's an error in any transaction.
		}
	}

	return nil
}

func (r *FirestoreFeedRepository) DeleteFeed(name string) bool {
	_, err := r.client.Collection("feeds").Doc(name).Delete(context.Background())
	return err == nil
}

func (r *FirestoreFeedRepository) DeleteTweet(deletedTweet models.Tweet) bool {
	ctx := context.Background()

	// If there are no tags, there's no feed to delete the tweet from.
	if len(deletedTweet.Tags) == 0 {
		return false
	}

	// Loop over each tag and update the corresponding feed document.
	for _, tag := range deletedTweet.Tags {
		feedDocRef := r.client.Collection("feeds").Doc(tag)

		// Run a Firestore transaction to safely check and delete the tweet.
		err := r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			// Check if the feed document exists.
			doc, err := tx.Get(feedDocRef)
			if err != nil {
				if status.Code(err) == codes.NotFound {
					// Document does not exist, so we don't need to delete anything.
					return nil
				} else {
					return err // Return any other error.
				}
			}

			var feedData struct {
				Tweets []models.Tweet `firestore:"tweets"`
			}

			err = doc.DataTo(&feedData)
			if err != nil {
				return err
			}

			// Filter out the deleted tweet from the tweets array.
			var updatedTweets []models.Tweet
			for _, t := range feedData.Tweets {
				if t.ID != deletedTweet.ID {
					updatedTweets = append(updatedTweets, t)
				}
			}

			// Update the tweets array with the filtered tweets.
			err = tx.Update(feedDocRef, []firestore.Update{
				{
					Path:  "tweets",
					Value: updatedTweets,
				},
			})
			return err
		})

		if err != nil {
			return false // Return if there's an error in any transaction.
		}
	}

	return true
}

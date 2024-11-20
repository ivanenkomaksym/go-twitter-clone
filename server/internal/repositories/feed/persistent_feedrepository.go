package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"twitter-clone/internal/config"
	"twitter-clone/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PersistentFeedRepository struct {
	client          *mongo.Client
	feedsCollection *mongo.Collection
}

func NewPersistentFeedRepository(configuration config.Configuration) (*PersistentFeedRepository, error) {
	repo := &PersistentFeedRepository{}
	initComplete := make(chan error)

	// Initialize the database connection asynchronously
	go func() {
		err := repo.init(configuration)
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
		initComplete <- err // Send the result of the init to the channel
	}()

	// Wait for initialization to complete
	err := <-initComplete
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	return repo, nil
}

func (repo *PersistentFeedRepository) init(configuration config.Configuration) error {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(configuration.FeedsStorage.ConnectionString).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	repo.client = client
	repo.feedsCollection = client.Database(configuration.FeedsStorage.DatabaseName).Collection(configuration.FeedsStorage.CollectionName)

	return nil
}

func (repo *PersistentFeedRepository) CreateFeed(name string) error {
	_, err := repo.GetFeedByName(name)
	if err != nil {
		return err
	}

	feed := models.Feed{
		Name:   name,
		Tweets: []models.Tweet{},
	}

	insertOneResult, err := repo.feedsCollection.InsertOne(context.Background(), feed)
	if err != nil {
		if !isDuplicateError(err) {
			return err
		}
		return nil
	}

	if insertOneResult.InsertedID == nil {
		return errors.New("failed to insert object")
	}

	return nil
}

func isDuplicateError(err error) bool {
	mErr, ok := err.(mongo.WriteException)
	if !ok {
		return false
	}

	return mErr.WriteErrors[0].Code == 11000
}

func (repo *PersistentFeedRepository) GetFeeds() ([]models.Feed, error) {
	var feeds []models.Feed

	cursor, err := repo.feedsCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return feeds, err
	}

	if err = cursor.All(context.Background(), &feeds); err != nil {
		panic(err)
	}

	return feeds, nil
}

func (repo *PersistentFeedRepository) GetFeedByName(name string) (*models.Feed, error) {
	filter := bson.D{{Key: "_id", Value: name}}

	var result models.Feed
	var foundResult = repo.feedsCollection.FindOne(context.Background(), filter)
	if foundResult.Err() == nil {
		foundResult.Decode(&result)
		return &result, nil
	}

	return nil, nil
}

func (repo *PersistentFeedRepository) AppendTweet(tweet models.Tweet) error {
	if len(tweet.Tags) == 0 {
		return nil
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": tweet.Tags,
		},
		"tweets.id": bson.M{
			"$ne": tweet.ID,
		},
	}

	update := bson.M{
		"$push": bson.M{
			"tweets": bson.M{
				"$each":     bson.A{tweet},
				"$position": 0,
			},
		},
	}

	_, err := repo.feedsCollection.UpdateMany(context.Background(), filter, update)
	return err
}

func (repo *PersistentFeedRepository) DeleteFeed(name string) bool {
	filter := bson.D{{Key: "_id", Value: name}}

	deleteResult, err := repo.feedsCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("Error deleting feed: %v", err)
		return false
	}

	return deleteResult.DeletedCount > 0
}

func (repo *PersistentFeedRepository) DeleteTweet(deletedTweet models.Tweet) bool {
	log.Printf("Deleting tweet: %v", deletedTweet)

	filter := bson.M{
		"_id": bson.M{
			"$in": deletedTweet.Tags,
		},
	}

	update := bson.M{
		"$pull": bson.M{
			"tweets": bson.M{
				"id": deletedTweet.ID,
			},
		},
	}

	_, err := repo.feedsCollection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error deleting tweet: %v", err)
		return false
	}

	return true
}

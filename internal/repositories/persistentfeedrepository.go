package repositories

import (
	"context"
	"fmt"
	"log"
	"twitter-clone/internal/config"
	"twitter-clone/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PersistentFeedRepository struct {
	client *mongo.Client
}

func NewPersistentFeedRepository(configuration config.Configuration) (*PersistentFeedRepository, error) {
	repo := &PersistentFeedRepository{}

	// Initialize the database connection asynchronously
	go func() {
		err := repo.init(configuration)
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
	}()

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

	return nil
}

func (repo *PersistentFeedRepository) CreateFeed(name string) error {
	return nil
}

func (repo *PersistentFeedRepository) GetFeeds() ([]models.Feed, error) {
	var feeds []models.Feed
	return feeds, nil
}

func (repo *PersistentFeedRepository) GetFeedByName(name string) (*models.Feed, error) {
	return nil, nil
}

func (repo *PersistentFeedRepository) AppendTweet(tweet models.Tweet) error {
	return nil
}

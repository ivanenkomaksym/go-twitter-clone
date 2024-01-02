package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"twitter-clone/internal/config"
	"twitter-clone/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type PersistentTweetRepository struct {
	db *sql.DB
}

func NewPersistentTweetRepository(configuration config.Configuration) (*PersistentTweetRepository, error) {
	repo := &PersistentTweetRepository{}

	// Initialize the database connection asynchronously
	go func() {
		err := repo.init(configuration)
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
	}()

	return repo, nil
}

func (repo *PersistentTweetRepository) init(configuration config.Configuration) error {
	// Construct the full connection string
	connString := fmt.Sprintf("%s/%s", configuration.TweetsStorage.ConnectionString, configuration.TweetsStorage.DatabaseName)
	log.Println(connString)
	// Open the database connection
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return err
	}

	// Ping the database to ensure connectivity
	err = db.Ping()
	if err != nil {
		return err
	}

	repo.db = db
	return nil
}

func (repo *PersistentTweetRepository) CreateTweet(tweet models.Tweet) *models.Tweet {
	// Implement the logic to insert a tweet into the MySQL database
	// ...

	return &tweet
}

func (repo *PersistentTweetRepository) GetTweets() []models.Tweet {
	// Implement the logic to retrieve all tweets from the MySQL database
	// ...

	return []models.Tweet{}
}

func (repo *PersistentTweetRepository) GetTweetById(id string) *models.Tweet {
	// Implement the logic to retrieve a tweet by ID from the MySQL database
	// ...

	return nil
}

func (repo *PersistentTweetRepository) DeleteTweet(id string) bool {
	// Implement the logic to delete a tweet by ID from the MySQL database
	// ...

	return true
}

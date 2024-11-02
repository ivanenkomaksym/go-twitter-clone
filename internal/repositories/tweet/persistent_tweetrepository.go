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
	initComplete := make(chan error)

	// Initialize the database connection asynchronously
	go func() {
		err := repo.init(configuration)
		initComplete <- err // Send the result of the init to the channel
	}()

	// Wait for initialization to complete
	err := <-initComplete
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	return repo, nil
}

func (repo *PersistentTweetRepository) init(configuration config.Configuration) error {
	// Get the connection string without the database name
	connString := fmt.Sprintf("%s/", configuration.TweetsStorage.ConnectionString)
	log.Println("Connecting without database:", connString)

	// Open the database connection (without specifying a database)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return err
	}

	// Check if the database exists
	dbName := configuration.TweetsStorage.DatabaseName
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		return fmt.Errorf("error creating database: %v", err)
	}

	// Now that the database exists, close the connection and reopen with the database name
	err = db.Close()
	if err != nil {
		return fmt.Errorf("error closing connection: %v", err)
	}

	// Reconnect with the database specified
	connStringWithDB := fmt.Sprintf("%s/%s", configuration.TweetsStorage.ConnectionString, dbName)
	db, err = sql.Open("mysql", connStringWithDB)
	if err != nil {
		return fmt.Errorf("error reconnecting to database: %v", err)
	}

	// Ping the database to ensure connectivity
	err = db.Ping()
	if err != nil {
		return err
	}

	repo.db = db

	// SQL to create the 'tweets' table if it does not exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tweets (
		id VARCHAR(36) PRIMARY KEY,
		title VARCHAR(255),
		content TEXT,
		author VARCHAR(255),
		created_at TIMESTAMP
	)
	`
	_, err = repo.db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error creating 'tweets' table: %v", err)
		return err
	}

	return nil
}

func (repo *PersistentTweetRepository) CreateTweet(createTweetRequest models.CreateTweetRequest) *models.Tweet {
	tweet := CreateNewTweet(createTweetRequest)
	// Check if the tweet with the given ID already exists
	if existingTweet := repo.GetTweetById(tweet.ID); existingTweet != nil {
		log.Printf("Tweet with ID '%s' already exists", tweet.ID)
		return nil
	}

	// Perform the actual insertion into the database
	// TODO: insert user into database
	_, err := repo.db.Exec("INSERT INTO tweets (id, title, content, author, created_at) VALUES (?, ?, ?, ?, ?)",
		tweet.ID, tweet.Title, tweet.Content, tweet.User, tweet.CreatedAt.Time)
	if err != nil {
		log.Printf("Error inserting tweet into database: %v", err)
		return nil
	}

	// Return the created tweet
	return &tweet
}

func (repo *PersistentTweetRepository) GetTweets() []models.Tweet {
	rows, err := repo.db.Query("SELECT id, title, content, author, created_at FROM tweets")
	if err != nil {
		log.Printf("Error retrieving tweets from database: %v", err)
		return nil
	}
	defer rows.Close()

	var tweets []models.Tweet
	for rows.Next() {
		var tweet models.Tweet
		// TODO: read user from database
		var author string
		err := rows.Scan(&tweet.ID, &tweet.Title, &tweet.Content, &author, &tweet.CreatedAt)
		if err != nil {
			log.Printf("Error scanning tweet rows: %v", err)
			return nil
		}

		tweets = append(tweets, tweet)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over tweet rows: %v", err)
		return nil
	}

	return tweets
}

func (repo *PersistentTweetRepository) GetTweetById(id string) *models.Tweet {
	row := repo.db.QueryRow("SELECT id, title, content, author, created_at FROM tweets WHERE id = ?", id)

	var tweet models.Tweet
	// TODO: read user from database
	var author string
	err := row.Scan(&tweet.ID, &tweet.Title, &tweet.Content, &author, &tweet.CreatedAt)
	if err == sql.ErrNoRows {
		// No tweet found with the given ID
		return nil
	} else if err != nil {
		log.Printf("Error retrieving tweet from database: %v", err)
		return nil
	}

	return &tweet
}

func (repo *PersistentTweetRepository) DeleteTweet(id string) bool {
	result, err := repo.db.Exec("DELETE FROM tweets WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting tweet from database: %v", err)
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected after tweet deletion: %v", err)
		return false
	}

	if rowsAffected == 0 {
		// No tweet found with the given ID
		return false
	}

	return true
}

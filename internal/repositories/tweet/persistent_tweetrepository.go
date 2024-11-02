package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"twitter-clone/internal/config"
	"twitter-clone/internal/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
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
	createUsersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		email VARCHAR(255) UNIQUE,
		picture TEXT
	)`

	_, err = repo.db.Exec(createUsersTableSQL)
	if err != nil {
		log.Printf("Error creating 'users' table: %v", err)
		return err
	}

	createTweetsTableSQL := `
	CREATE TABLE IF NOT EXISTS tweets (
		id VARCHAR(36) PRIMARY KEY,
		title VARCHAR(255),
		content TEXT,
		created_at TIMESTAMP,
		user_id VARCHAR(36),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`

	_, err = repo.db.Exec(createTweetsTableSQL)
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

	// Insert or update the user

	// Check if user already exists by email
	var userID string
	err := repo.db.QueryRow("SELECT id FROM users WHERE email = ?", tweet.User.Email).Scan(&userID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking user existence in database: %v", err)
		return nil
	}

	// If the user does not exist, insert a new user
	if err == sql.ErrNoRows {
		userID = uuid.NewString()
		_, err := repo.db.Exec(`
        INSERT INTO users (id, first_name, last_name, email, picture) 
        VALUES (?, ?, ?, ?, ?)
    `, userID, tweet.User.FirstName, tweet.User.LastName, tweet.User.Email, tweet.User.Picture)
		if err != nil {
			log.Printf("Error inserting new user into database: %v", err)
			return nil
		}
	}

	// Insert the tweet with a reference to the user_id
	_, err = repo.db.Exec(`
		INSERT INTO tweets (id, title, content, created_at, user_id) 
		VALUES (?, ?, ?, ?, ?)
	`, tweet.ID, tweet.Title, tweet.Content, tweet.CreatedAt.Time, userID)
	if err != nil {
		log.Printf("Error inserting tweet into database: %v", err)
		return nil
	}

	// Return the created tweet
	return &tweet
}

func (repo *PersistentTweetRepository) GetTweets() []models.Tweet {
	// Query to fetch tweets along with user details
	rows, err := repo.db.Query(`
		SELECT t.id, t.title, t.content, t.created_at, 
		       u.id AS user_id, u.first_name, u.last_name, u.email, u.picture
		FROM tweets t
		JOIN users u ON t.user_id = u.id
	`)
	if err != nil {
		log.Printf("Error retrieving tweets from database: %v", err)
		return nil
	}
	defer rows.Close()

	var tweets []models.Tweet
	for rows.Next() {
		var tweet models.Tweet
		var user models.User
		var userID string

		// Scan the values from the row into the tweet and user structs
		err := rows.Scan(
			&tweet.ID,
			&tweet.Title,
			&tweet.Content,
			&tweet.CreatedAt,
			&userID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Picture,
		)
		if err != nil {
			log.Printf("Error scanning tweet row: %v", err)
			return nil
		}

		// Set the user struct in the tweet and add to the tweets slice
		tweet.User = user
		tweets = append(tweets, tweet)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over tweet rows: %v", err)
		return nil
	}

	return tweets
}

func (repo *PersistentTweetRepository) GetTweetById(id string) *models.Tweet {
	// Query to fetch a single tweet along with user details by tweet ID
	row := repo.db.QueryRow(`
		SELECT t.id, t.title, t.content, t.created_at, 
		       u.id AS user_id, u.first_name, u.last_name, u.email, u.picture
		FROM tweets t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = ?
	`, id)

	var tweet models.Tweet
	var user models.User
	var userID string

	// Scan the result into the tweet and user structs
	err := row.Scan(
		&tweet.ID,
		&tweet.Title,
		&tweet.Content,
		&tweet.CreatedAt,
		&userID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Picture,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Printf("Error retrieving tweet by ID from database: %v", err)
		return nil
	}

	// Set the user struct in the tweet
	tweet.User = user

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

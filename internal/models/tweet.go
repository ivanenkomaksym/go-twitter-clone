package models

import (
	"fmt"
	"time"
)

// MySQLTimestamp is a custom type for scanning time.Time from MySQL
type MySQLTimestamp struct {
	time.Time
}

// Scan implements the sql.Scanner interface
func (t *MySQLTimestamp) Scan(value interface{}) error {
	switch v := value.(type) {
	case []uint8:
		// Parse the []uint8 value into a time.Time
		parseTime, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		t.Time = parseTime
	default:
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type *MySQLTimestamp", value)
	}
	return nil
}

type Tweet struct {
	ID        string         `json:"id" bson:"id"`
	Title     string         `json:"title" bson:"title"`
	Content   string         `json:"content" bson:"content"`
	Author    string         `json:"author" bson:"author"`
	Tags      []string       `json:"tags" bson:"tags"`
	CreatedAt MySQLTimestamp `json:"created_at" bson:"created_at"`
	Likes     []Like         `json:"likes" bson:"likes"`
}

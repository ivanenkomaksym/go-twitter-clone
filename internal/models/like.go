package models

import "time"

type Like struct {
	ID        string    `json:"id" bson:"id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

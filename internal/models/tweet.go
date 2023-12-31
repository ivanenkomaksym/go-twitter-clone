package models

import "time"

type Tweet struct {
	ID        string    `json:"id" bson:"id"`
	Title     string    `json:"title" bson:"title"`
	Content   string    `json:"content" bson:"content"`
	Author    string    `json:"author" bson:"author"`
	Tags      []string  `json:"tags" bson:"tags"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Likes     []Like    `json:"likes" bson:"likes"`
}

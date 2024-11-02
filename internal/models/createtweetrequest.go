package models

type CreateTweetRequest struct {
	Title   string   `json:"title" bson:"title"`
	Content string   `json:"content" bson:"content"`
	Tags    []string `json:"tags" bson:"tags"`
}

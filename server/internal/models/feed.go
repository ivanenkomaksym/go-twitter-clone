package models

type Feed struct {
	Name   string  `json:"name" bson:"_id"`
	Tweets []Tweet `json:"tweets" bson:"tweets"`
}

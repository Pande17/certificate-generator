package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Counter struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Month   string             `json:"month" bson:"month"`
	Year    string             `json:"year" bson:"year"`
	Counter int64              `json:"counter" bson:"counter"`
}

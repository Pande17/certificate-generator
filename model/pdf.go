package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PDF struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	DataID string
	Data   CertificateData `json:"data" bson:"data"`
	Model
}

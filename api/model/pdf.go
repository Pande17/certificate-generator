package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PDF struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DataID     string             `json:"data_id" bson:"data_id"`
	SertifName string             `json:"sertif_name" bson:"sertif_name"`
	Data       CertificateData    `json:"data" bson:"data"`
	Model      `bson:",inline"`   // flatten the model fields
}

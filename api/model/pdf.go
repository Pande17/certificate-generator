package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type PDF struct {
	AdminId    primitive.ObjectID `json:"admin_id" bson:"admin_id"`
	DataID     string             `json:"data_id" bson:"data_id"`
	SertifName string             `json:"sertif_name" bson:"sertif_name"`
	Data       CertificateData    `json:"data" bson:"data"`
	Model      `bson:",inline"`   // flatten the model fields
}

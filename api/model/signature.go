package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Signature struct {
	AdminId    primitive.ObjectID `json:"admin_id" bson:"admin_id"`
	ConfigName string             `json:"config_name" bson:"config_name"`
	Stamp      string             `json:"stamp" bson:"stamp"`
	Logo       string             `json:"logo" bson:"logo"`
	Signature  string             `json:"signature" bson:"signature"`
	Name       string             `json:"name" bson:"name"`
	Role       string             `json:"role" bson:"role"`
	Model      `bson:",inline"`
}

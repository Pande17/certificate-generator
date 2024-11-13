package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// struct for Admin Account
type AdminAccount struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AdminName     string             `json:"admin_name" bson:"admin_name"`
	AdminPassword string             `json:"admin_password" bson:"admin_password" `
	Model         `bson:",inline"`   // flatten the model fields
}

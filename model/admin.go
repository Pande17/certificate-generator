package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// struct for Admin Account
type AdminAccount struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" `
	AccID         int64              `json:"acc_id" bson:"acc_id"`
	AdminName     string             `json:"admin_name" bson:"admin_name"`
	AdminPassword string             `json:"admin_password" bson:"admin_password" `
	Model
}

// ADMIN ACCOUNT LIST:
// usrname : password
// pande : pande
// pande2 : pande2

// note: AdminPassword == AdminName

package dbmongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type Admin struct {
	ID      		primitive.ObjectID 	`bson:"_id,omitempty" `
	AdminName 		string				`bson:"admin_name" gorm:"unique"`
	AdminPassword	string				`bson:"admin_password" `
	Model
}
package dbmongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type AdminAccount struct {
	ID      		primitive.ObjectID 	`bson:"_id,omitempty" `
	AdminName 		string				`bson:"admin_name"`
	AdminPassword	string				`bson:"admin_password" `
	Model
}
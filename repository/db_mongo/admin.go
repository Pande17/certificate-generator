package dbmongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type AdminAccount struct {
	ID      		primitive.ObjectID 	`bson:"_id,omitempty" `
	AccID			int64				`bson:"acc_id"`
	AdminName 		string				`bson:"admin_name"`
	AdminPassword	string				`bson:"admin_password" `
	Model
}

// ADMIN ACCOUNT LIST:
// usrname : password
// pande : pande
// pande2 : pande2

// note: AdminPassword == AdminName
package model

// struct for Admin Account
type AdminAccount struct {
	AdminName     string           `json:"admin_name" bson:"admin_name"`
	AdminPassword string           `json:"admin_password" bson:"admin_password" `
	Model         `bson:",inline"` // flatten the model fields
}

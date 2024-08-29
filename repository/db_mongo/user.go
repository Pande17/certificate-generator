package dbmongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User struct untuk tabel user
type User struct {
	ID               primitive.ObjectID 	`bson:"_id,omitempty" json:"id"`
	UserName         string             	`bson:"userName" json:"userName"`
	Email            string             	`bson:"email" json:"email"`
	Password         string             	`bson:"password" json:"password"`
	Files            []PdfData				`bson:"files" json:"files"`
	JumlahSertifikat uint               	`bson:"jumlahSertifikat" json:"jumlahSertifikat"`
	Model
}
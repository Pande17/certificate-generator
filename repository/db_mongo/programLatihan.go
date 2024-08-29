package dbmongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProgramLatihan struct untuk tabel programLatihan
type ProgramLatihan struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NamaProgramLat string             `bson:"namaProgramLat" json:"namaProgramLat"`
	BidangKomId   primitive.ObjectID `bson:"bidangKomId" json:"bidangKomId"`
	Model
}
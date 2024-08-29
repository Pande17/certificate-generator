package dbmongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BidangKompetensi struct untuk tabel bidangKompetensi
type BidangKompetensi struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProgramLatId  primitive.ObjectID `bson:"programLatId" json:"programLatId"`
	NamaBidangKom string             `bson:"namaBidangKom" json:"namaBidangKom"`
	Model
}
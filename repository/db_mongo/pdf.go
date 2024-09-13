package dbmongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type PDF struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	SertifName    string
	SertifVersion string
	SertifID      string
	KodeReferal
	NamaPeserta string
	SKKNI       string
	Kompetensi  string
	ValidDate
	Validation string
	KodeQR
	DataID   string
	TotalJP  uint64
	Meet     uint64
	MeetTime float64
	HardSkills
	SoftSkills
	FinalSkor
	Model
}

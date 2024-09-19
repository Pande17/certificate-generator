package dbmongo

import (
	"encoding/base64"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PDF struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	SertifName string             `bson:"sertif_name"`
	KodeReferal
	NamaPeserta    string `bson:"nama_peserta"`
	SKKNI          string `bson:"skkni"`
	KompetenBidang string `bson:"kompeten_bidang"`
	Kompetensi     string `bson:"kompetensi"`
	ValidDate
	Validation string `bson:"validation"`
	KodeQR     base64.Encoding
	DataID     string `bson:"data_id"`
	TotalJP    uint64 `bson:"total_jp"`
	TotalMeet  uint64 `bson:"total_meet"`
	MeetTime   string `bson:"meet_time"`
	HardSkillPDF
	SoftSkillPDF
	FinalSkor float64
	Model
}

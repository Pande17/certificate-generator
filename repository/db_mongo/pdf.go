package dbmongo

import (
	"encoding/base64"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PDF struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	SertifName     string             `json:"sertif_name" bson:"sertif_name"`
	NamaPeserta    string             `json:"nama_peserta" bson:"nama_peserta"`
	SKKNI          string             `json:"skkni" bson:"skkni"`
	KompetenBidang string             `json:"kompeten_bidang" bson:"kompeten_bidang"`
	Kompetensi     string             `json:"kompetensi" bson:"kompetensi"`
	Validation     string             `json:"validation" bson:"validation"`
	KodeQR         base64.Encoding    `json:"kode_qr" bson:"kode_qr"`
	DataID         string             `json:"data_id" bson:"data_id"`
	TotalJP        uint64             `json:"total_jp" bson:"total_jp"`
	TotalMeet      uint64             `json:"total_meet" bson:"total_meet"`
	MeetTime       string             `json:"meet_time" bson:"meet_time"`
	FinalSkor      float64            `json:"final_skor" bson:"final_skor"`
	KodeReferal
	ValidDate
	HardSkillPDF
	SoftSkillPDF
	Model
}

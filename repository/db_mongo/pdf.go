package dbmongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type PDF struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	SertifName    string             `bson:"sertif_name"`
	SertifVersion string             `bson:"sertif_version"`
	SertifID      string             `bson:"sertif_id"`
	NamaPeserta   string             `bson:"nama_peserta"`
	SKKNI         string             `bson:"skkni"`
	Kompetensi    string             `bson:"kompetensi"`
	Validation    string             `bson:"validation"`
	DataID        string             `bson:"data_id"`
	TotalJP       uint64             `bson:"total_jp"`
	Meet          uint64             `bson:"meet"`
	MeetTime      float64            `bson:"meet_time"`
	HardSkills    []HardSkill        // ini nggak tau cara implementasinya
	SoftSkills    []SoftSkill        // ini juga
	KodeReferal
	ValidDate
	KodeQR
	FinalSkor
	Model
}

package dbmongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DataNilai struct untuk tabel dataNilai
type DataNilai struct {
	ID             primitive.ObjectID 	`bson:"_id,omitempty" json:"id"`
	FileId         primitive.ObjectID 	`bson:"fileId" json:"file_id"`
	ProgramLatihan string             	`bson:"programLatihan" json:"program_latihan"`
	TotalJP        uint               	`bson:"totalJP" json:"total_jp"`
	TotalPertemuan uint               	`bson:"totalPertemuan" json:"total_pertemuan"`
	Kompeten       float64            	`bson:"kompeten" json:"kompeten"`
	HardSkillData  []HardSkillData 		`bson:"hardSkillData" json:"hard_skill_data"`
	SoftSkillData  []SoftSkillData 		`bson:"softSkillData" json:"soft_skill_data"`
	NilaiAkhir     float64            	`bson:"nilaiAkhir" json:"nilai_akhir"`
	Model
}
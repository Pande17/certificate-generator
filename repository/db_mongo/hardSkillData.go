package dbmongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HardSkillData struct untuk tabel hardSkillData
type HardSkillData struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DataNilaiId   primitive.ObjectID `bson:"dataNilaiId" json:"data_nilai_id"`
	ProgramLatId  primitive.ObjectID `bson:"programLatId" json:"program_latihan_id"`
	HardSkills    []interface{}      `bson:"hardSkills" json:"hardSkills"`
	HsJp          []interface{}      `bson:"hsJp" json:"hsJp"`
	HsSkor        []interface{}      `bson:"hsSkor" json:"hsSkor"`
	Model
}
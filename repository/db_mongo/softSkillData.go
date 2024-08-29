package dbmongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SoftSkillData struct untuk tabel softSkillData
type SoftSkillData struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DataNilaiId   primitive.ObjectID `bson:"dataNilaiId" json:"dataNilaiId"`
	SoftSkills    []interface{}      `bson:"softSkills" json:"softSkills"`
	SsJp          []interface{}      `bson:"ssJp" json:"ssJp"`
	SsSkor        []interface{}      `bson:"ssSkor" json:"ssSkor"`
	Model
}
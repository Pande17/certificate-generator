package dbmongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// struct for Kompetensi
type Kompetensi struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	KompetensiID   uint64             `json:"kompetensi_id" bson:"kompetensi_id"`
	NamaKompetensi string             `json:"nama_kompetensi" bson:"nama_kompetensi"`
	HardSkills     []HardSkill        `json:"hard_skills" bson:"hard_skills"`
	SoftSkills     []SoftSkill        `json:"soft_skills" bson:"soft_skills"`
	Model
}

// struct for Hard Skills
type HardSkill struct {
	HardSkillName string        `json:"hardSkill_name" bson:"hardSkill_name"`
	Descriptions  []Description `json:"description" bson:"description"`
}

// struct for Soft Skills
type SoftSkill struct {
	SoftSkillName string        `json:"softSkill_name" bson:"softSkill_name"`
	Descriptions  []Description `json:"description" bson:"description"`
}

// struct for Deskripsi Hard Skill and Soft Skill
type Description struct {
	UnitCode  string `json:"unit_code" bson:"unit_code"`
	UnitTitle string `json:"unit_title" bson:"unit_title"`
}

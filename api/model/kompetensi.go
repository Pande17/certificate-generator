package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Struct for Kompetensi
type Kompetensi struct {
	AdminId        primitive.ObjectID `json:"admin_id" bson:"admin_id"`
	NamaKompetensi string             `json:"nama_kompetensi" bson:"nama_kompetensi"`
	Divisi         string             `json:"divisi" bson:"divisi"`
	SKKNI          string             `json:"skkni" bson:"skkni"`
	HardSkills     []Skill            `json:"hard_skills" bson:"hard_skills"`
	SoftSkills     []Skill            `json:"soft_skills" bson:"soft_skills"`
	Model          `bson:",inline"`   // Flatten the model fields
}

// Struct for Skill
type Skill struct {
	SkillName  string        `json:"skill_name" bson:"skill_name"`
	SkillDescs []Description `json:"description" bson:"description"`
	SkillJP    uint64        `json:"skill_jp" bson:"skill_jp"`
	SkillScore float64       `json:"skill_score" bson:"skill_score"`
}

// Struct for Description of Hard Skill and Soft Skill
type Description struct {
	UnitCode  string `json:"unit_code" bson:"unit_code" valid:"required~Kode Unit tidak boleh kosong!"`
	UnitTitle string `json:"unit_title" bson:"unit_title" valid:"required~Judul Unit tidak boleh kosong!"`
}

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// struct for Kompetensi
type Kompetensi struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	KompetensiID   uint64             `json:"kompetensi_id" bson:"kompetensi_id"`
	NamaKompetensi string             `json:"nama_kompetensi" bson:"nama_kompetensi"`
	HardSkills     []Skill            `json:"hard_skills" bson:"hard_skills"`
	SoftSkills     []Skill            `json:"soft_skills" bson:"soft_skills"`
	Model          struct {
		CreatedAt time.Time  `json:"created_at" bson:"created_at"`
		UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
	} `bson:",inline"`
}

// `bson:",inline"`   // flatten the model fields

type Skill struct {
	SkillName  string        `json:"skill_name" bson:"skill_name"`
	SkillDescs []Description `json:"description" bson:"description"`
	SkillJP    uint64        `json:"skill_jp" bson:"skill_jp"`
	SkillScore float64       `json:"skill_score" bson:"skill_score"`
}

// struct for Deskripsi Hard Skill and Soft Skill
type Description struct {
	UnitCode  string `json:"unit_code" bson:"unit_code"`
	UnitTitle string `json:"unit_title" bson:"unit_title"`
}

// struct for Hard Skills
// type HardSkill struct {
// 	// HardSkillName  string        `json:"hardSkill_name" bson:"hardSkill_name"`
// 	// Descriptions   []Description `json:"description" bson:"description"`
// 	// HardSkillJP    uint64        `json:"hardSkill_jp" bson:"hardSkill_jp"`
// 	// HardSkillScore float64       `json:"hardSkill_score" bson:"hardSkill_score"`
// 	Skill `json:"hard_skill" bson:"hard_skill"`
// }

// // struct for Soft Skills
// type SoftSkill struct {
// 	// SoftSkillName  string        `json:"softSkill_name" bson:"softSkill_name"`
// 	// Descriptions   []Description `json:"description" bson:"description"`
// 	// SoftSkillJP    uint64        `json:"softSkill_jp" bson:"softSkill_jp"`
// 	// SoftSkillScore float64       `json:"softSkill_score" bson:"softSkill_score"`
// 	Skill `json:"soft_skill" bson:"soft_skill"`
// }

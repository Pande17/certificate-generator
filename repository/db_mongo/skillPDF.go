package dbmongo

type HardSkillPDF struct {
	// HardSkillName       string  `json:"hard_skill_name" bson:"hard_skill_name"`
	// HardSkillCode       string  `json:"hard_skill_code" bson:"hard_skill_code"`
	// HardSkillDesc       string  `json:"hard_skill_desc" bson:"hard_skill_desc"`
	HardSkill
	HardSkillJP         uint64  `json:"hardSkill_jp" bson:"hardSkill_jp"`
	HardSkillScore      float64 `json:"hardSkill_score" bson:"hardSkill_score"`
	TotalHardSkillJP    string  `json:"total_hardSkill_jp" bson:"total_hardSkill_jp"`
	TotalHardSkillScore float64 `json:"total_hardSkill_score" bson:"total_hardSkill_score"`
}

type SoftSkillPDF struct {
	// SoftSkillName       string  `json:"soft_skill_name" bson:"soft_skill_name"`
	// SoftSkillCode       string  `json:"soft_skill_code" bson:"soft_skill_code"`
	// SoftSkillDesc       string  `json:"soft_skill_desc" bson:"soft_skill_desc"`
	SoftSkill
	SoftSkillJP         uint64  `json:"softSkill_jp" bson:"softSkill_jp"`
	SoftSkillScore      float64 `json:"softSkill_score" bson:"softSkill_score"`
	TotalSoftSkillJP    string  `json:"total_softSkill_jp" bson:"total_softSkill_jp"`
	TotalSoftSkillScore float64 `json:"total_softSkill_score" bson:"total_softSkill_score"`
}

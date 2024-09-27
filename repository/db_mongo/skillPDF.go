package dbmongo

type HardSkillPDF struct {
	HardSkillName       string  `json:"hard_skill_name" bson:"hard_skill_name"`
	HardSkillCode       string  `json:"hard_skill_code" bson:"hard_skill_code"`
	HardSkillDesc       string  `json:"hard_skill_desc" bson:"hard_skill_desc"`
	HardSkillJP         uint64  `json:"hard_skill_jp" bson:"hard_skill_jp"`
	HardSkillScore      float64 `json:"hard_skill_score" bson:"hard_skill_score"`
	TotalHardSkillJP    string  `json:"total_hard_skill_jp" bson:"total_hard_skill_jp"`
	TotalHardSkillScore float64 `json:"total_hard_skill_score" bson:"total_hard_skill_score"`
}

type SoftSkillPDF struct {
	SoftSkillName       string  `json:"soft_skill_name" bson:"soft_skill_name"`
	SoftSkillCode       string  `json:"soft_skill_code" bson:"soft_skill_code"`
	SoftSkillDesc       string  `json:"soft_skill_desc" bson:"soft_skill_desc"`
	SoftSkillJP         uint64  `json:"soft_skill_jp" bson:"soft_skill_jp"`
	SoftSkillScore      float64 `json:"soft_skill_score" bson:"soft_skill_score"`
	TotalSoftSkillJP    string  `json:"total_soft_skill_jp" bson:"total_soft_skill_jp"`
	TotalSoftSkillScore float64 `json:"total_soft_skill_score" bson:"total_soft_skill_score"`
}

package dbmongo

type HardSkillPDF struct {
	HardSkillName       string
	HardSkillCode       string
	HardSkillDesc       string
	HardSkillJP         uint64
	HardSkillScore      float64
	TotalHardSkillJP    string
	TotalHardSkillScore float64
}

type SoftSkillPDF struct {
	SoftSkillName       string
	SoftSkillCode       string
	SoftSkillDesc       string
	SoftSkillJP         uint64
	SoftSkillScore      float64
	TotalSoftSkillJP    string
	TotalSoftSkillScore float64
}

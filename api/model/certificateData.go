package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type CertificateData struct {
	AdminId        primitive.ObjectID `json:"admin_id" bson:"admin_id"`
	SertifName     string             `json:"sertif_name" bson:"sertif_name"`
	SertifTitle    string             `json:"sertif_title" bson:"sertif_title,omitempty"`
	KodeReferral   KodeReferral       `json:"kode_referral" bson:"kode_referral"`
	NamaPeserta    string             `json:"nama_peserta" bson:"nama_peserta"`
	SKKNI          string             `json:"skkni" bson:"skkni"`
	KompetenBidang string             `json:"kompeten_bidang" bson:"kompeten_bidang"`
	Kompetensi     string             `json:"kompetensi" bson:"kompetensi"`
	Validation     string             `json:"validation" bson:"validation"`
	ValidDate      ValidDate          `json:"valid_date" bson:"valid_date"`
	DataID         string             `json:"data_id" bson:"data_id"`
	QRCode         string             `json:"kode_qr" bson:"kode_qr"`
	TotalJP        uint64             `json:"total_jp" bson:"total_jp"`
	TotalMeet      uint64             `json:"total_meet" bson:"total_meet"`
	MeetTime       string             `json:"meet_time" bson:"meet_time"`
	FinalSkor      float64            `json:"final_skor" bson:"final_skor"`
	HardSkills     SkillPDF           `json:"hard_skills" bson:"hard_skills"`
	SoftSkills     SkillPDF           `json:"soft_skills" bson:"soft_skills"`
	Signature      SignatureData      `json:"signature" bson:"signature"`
}

type KodeReferral struct {
	ReferralID int64  `json:"referral_id" bson:"referral_id"`
	Divisi     string `json:"divisi" bson:"divisi"`
	BulanRilis string `json:"bulan_rilis" bson:"bulan_rilis"`
	TahunRilis int    `json:"tahun_rilis" bson:"tahun_rilis"`
}

type ValidDate struct {
	ValidTotal string `json:"valid_total" bson:"valid_total"`
	ValidStart string `json:"valid_start" bson:"valid_start"`
	ValidEnd   string `json:"valid_end" bson:"valid_end"`
}

type SkillPDF struct {
	Skills          []Skill `json:"skills" bson:"skills"`
	TotalSkillJP    uint64  `json:"total_skill_jp" bson:"total_skill_jp"`
	TotalSkillScore float64 `json:"total_skill_score" bson:"total_skill_score"`
}

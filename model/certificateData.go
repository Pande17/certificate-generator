package model

type CertificateData struct {
	SertifName     string       `json:"sertif_name" bson:"sertif_name"`
	KodeReferral   KodeReferral `json:"kode_referral" bson:"kode_referral"`
	NamaPeserta    string       `json:"nama_peserta" bson:"nama_peserta"`
	SKKNI          string       `json:"skkni" bson:"skkni"`
	KompetenBidang string       `json:"kompeten_bidang" bson:"kompeten_bidang"`
	Kompetensi     string       `json:"kompetensi" bson:"kompetensi"`
	Validation     string       `json:"validation" bson:"validation"`
	ValidDate      ValidDate    `json:"valid_date" bson:"valid_date"`
	DataID         string       `json:"data_id" bson:"data_id"`
	TotalJP        uint64       `json:"total_jp" bson:"total_jp"`
	TotalMeet      uint64       `json:"total_meet" bson:"total_meet"`
	MeetTime       string       `json:"meet_time" bson:"meet_time"`
	FinalSkor      float64      `json:"final_skor" bson:"final_skor"`
	QRCode         `json:"kode_qr" bson:"kode_qr"`
	HardSkillPDF   `json:"hard_skills" bson:"hard_skills"`
	SoftSkillPDF   `json:"soft_skills" bson:"soft_skills"`
}

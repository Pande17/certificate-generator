package model

type PDF struct {
	DataID     string           `json:"data_id" bson:"data_id"`
	SertifName string           `json:"sertif_name" bson:"sertif_name"`
	Data       CertificateData  `json:"data" bson:"data"`
	Model      `bson:",inline"` // flatten the model fields
}

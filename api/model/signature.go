package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Signature struct {
	Model         `bson:",inline"`
	SignatureData `bson:",inline"`
}

type SignatureData struct {
	AdminId    primitive.ObjectID `json:"admin_id" bson:"admin_id,omitempty"`
	ConfigName string             `json:"config_name" bson:"config_name" valid:"required~Nama konfigurasi tidak boleh kosong!, stringlength(1|40)~Nama harus antara 1 hingga 60 karakter!"`
	Stamp      string             `json:"stamp" bson:"stamp" valid:"required~Stamp tidak boleh kosong!, url"`
	Logo       string             `json:"logo" bson:"logo" valid:"required~Logo tidak boleh kosong!, url"`
	Signature  string             `json:"signature" bson:"signature" valid:"required~Signature tidak boleh kosong!, url"`
	Name       string             `json:"name" bson:"name" valid:"required~Nama tidak boleh kosong!, stringlength(1|60)~Nama harus antara 1 hingga 60 karakter!"`
	Role       string             `json:"role" bson:"role" valid:"required~Peran tidak boleh kosong!, stringlength(1|60)~Peran harus antara 1 hingga 60 karakter!"`
}

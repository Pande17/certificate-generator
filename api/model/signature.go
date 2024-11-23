package model

type Signature struct {
	Stamp     string `bson:"stamp" json:"stamp"`
	Signature string `bson:"signature" json:"signature"`
	Name      string `bson:"name" json:"name"`
	Role      string `bson:"role" json:"role"`
	Model     `bson:",inline"`
}

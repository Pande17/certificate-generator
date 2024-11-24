package model

type Signature struct {
	Stamp     string `json:"stamp" bson:"stamp"`
	Signature string `json:"signature" bson:"signature"`
	Name      string `json:"name" bson:"name"`
	Role      string `json:"role" bson:"role"`
	Model     `bson:",inline"`
}

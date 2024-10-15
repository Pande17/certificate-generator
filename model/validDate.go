package model

type ValidDate struct {
	ValidTotal string `json:"valid_total" bson:"valid_total"`
	ValidStart string `json:"valid_start" bson:"valid_start"`
	ValidEnd   string `json:"valid_end" bson:"valid_end"`
}

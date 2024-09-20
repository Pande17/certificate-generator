package dbmongo

type ValidDate struct {
	ValidTotal		string `bson:"valid_total"`
	ValidStart		string `bson:"valid_start"`
	ValidEnd		string `bson:"valid_end"`
}
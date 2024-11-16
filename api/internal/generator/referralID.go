package generator

import (
	"certificate-generator/model"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GenerateReferralID(collection *mongo.Collection, createdAt time.Time) (int64, error) {
	month := fmt.Sprintf("%02d", createdAt.Month())
	year := fmt.Sprintf("%d", createdAt.Year())

	filter := bson.M{"month": month, "year": year}
	update := bson.M{"$inc": bson.M{"counter": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var updatedCounter model.Counter
	err := collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&updatedCounter)
	if err != nil {
		return 0, nil
	}

	return updatedCounter.Counter, err
}

// function for convert month to roman numerals
func MonthToRoman(month int) string {
	romans := []string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII"}
	if month >= 1 && month <=12 {
		return romans[month]
	}
	return ""
}

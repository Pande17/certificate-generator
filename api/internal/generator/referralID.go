package generator

import (
	model "certificate-generator/api/model"
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

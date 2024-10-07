package generator

import (
	"context"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Function for generating an random string ID for data ID
func GetRandomID(length int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rand.Seed(time.Now().UnixNano())
	id := make([]byte, length)
	for i := range id {
		id[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(id)
}

// function to check if ID is unique in MongoDB
func IsIDUnique(collection *mongo.Collection, data_id string) (bool, error) {
	filter := bson.M{"data_id": data_id}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// function to get a unique random ID in MongoDB
func GetUniqueRandomID(collection *mongo.Collection, length int) (string, error) {
	var data_id string
	for {
		data_id = GetRandomID(8)
		unique, err := IsIDUnique(collection, data_id)
		if err != nil {
			return "", err
		}
		if unique {
			break
		}
	}
	return data_id, nil
}

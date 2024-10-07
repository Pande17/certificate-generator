package generator

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Function for generating an incremental ID for any collection
func GetNextIncrementalID(collection *mongo.Collection, fieldName string) (int64, error) {
	// Define a filter to find the maximum value of the specified field (fieldName)
	opts := options.FindOne().SetSort(bson.D{{Key: fieldName, Value: -1}}) // Sort by the specified field descending

	// Create a map to hold the result
	var result bson.M
	var ctx = context.Background() // Define the context

	// Retrieve the last inserted document based on the specified field
	err := collection.FindOne(ctx, bson.M{}, opts).Decode(&result)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return 1, nil

		}
		return 0, fmt.Errorf("failed to find the last document: %v", err)
	}

	// Extract the value of the fieldName
	if lastID, ok := result[fieldName].(int64); ok {
		// Increment the last ID by 1
		return lastID + 1, nil
	}

	// If the field is not found or not the expected type
	return 0, fmt.Errorf("field %s not found or not an int64", fieldName)
}

package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global MongoDB client variable
var MongoClient *mongo.Client

// function to connect database MongoDB
func ConnectMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO"))

	var err error
	MongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Check the connection with a ping
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = MongoClient.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB!")
	return MongoClient, nil
}

func CreateCollectionsAndIndexes(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database("certificate-generator")

	// create collection (if there is none yet)
	collections := []string{"adminAcc", "certificate", "competence", "counters", "auditLog"}
	for _, collections := range collections {
		if err := db.CreateCollection(ctx, collections); err != nil {
			log.Printf("Collection %s already exists or an error occured: %v", collections, err)
		}
	}

	// creating indexes
	adminCollection := db.Collection("adminAcc")
	certificateCollection := db.Collection("certificate")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "admin_name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := adminCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	indexModel = mongo.IndexModel{
		Keys:    bson.D{{Key: "data_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = certificateCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	log.Println("Collections and indexes created successfully")
	return nil
}

// function to get collection from database
func GetCollection(collectionName string) *mongo.Collection {
	if MongoClient == nil {
		log.Fatal("MongoClient is not initialized, call ConnectMongoDB() first")
	}
	return MongoClient.Database("certificate-generator").Collection(collectionName)
}

package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global MongoDB client variable
var MongoClient *mongo.Client

// function to connect database MongoDB
func ConnectMongoDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb+srv://pande08:pande@belajar.ln6gf7x.mongodb.net/")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	MongoClient = client

	return MongoClient
}

// // function to get collection from database
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("certificate-generator").Collection(collectionName)
}

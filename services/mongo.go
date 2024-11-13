package services

import (
	"context"
	"log"
	"os"
	//"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var QuestionsCollection *mongo.Collection

// SetupDB connects the app to the database
func SetupDB() error {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not provided")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Cannot connect to MongoDB:", err)
		return err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Cannot connect to MongoDB:", err)
		return err
	}

	MongoClient = client
	QuestionsCollection = client.Database("test").Collection("questions")
	log.Println("Connected to MongoDB!")

	return nil
}

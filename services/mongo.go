package services

import (
	"context"
	"encoding/json"
	"fmt"

	// "io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var QuestionsCollection *mongo.Collection
var MetaCollection *mongo.Collection

// SetupDB connects the app to the database
func SetupDB() error {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not provided")
	}
	clientOptions := options.Client().ApplyURI(mongoURI)
	var client *mongo.Client
	var err error
	retries := 5
	for i := 0; i < retries; i++ {
		client, err = mongo.Connect(context.TODO(), clientOptions)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to MongoDB: %v. Retrying...", err)
		err = startMongoContainer()
		if err != nil {
			log.Printf("Failed to start MongoDB container: %v", err)
			return err
		}
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatal("Cannot connect to MongoDB:", err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Cannot connect to MongoDB:", err)
	}
	MongoClient = client
	QuestionsCollection = client.Database("test").Collection("questions")
	MetaCollection = client.Database("test").Collection("meta")
	log.Println("Connected to MongoDB!")

	// Initialize questions
	err = initializeQuestionsFromFile("services/questions.json")
	if err != nil {
		log.Printf("Error initializing questions: %v", err)
	}

	return nil
}

func startMongoContainer() error {
	cmd := exec.Command("docker", "start", "my_mongo_db")
	err := cmd.Run()
	if err != nil {
		return err
	}
	log.Println("MongoDB container started successfully")
	return nil
}

// initializeQuestionsFromFile loads the questions from a JSON file and inserts them into the database
func initializeQuestionsFromFile(filePath string) error {
	// Check if questions are already initialized
	var initStatus bson.M
	err := MetaCollection.FindOne(context.TODO(), bson.M{"key": "questions_initialized"}).Decode(&initStatus)
	if err == nil {
		log.Println("Questions already initialized. Skipping.")
		return nil
	}

	// Read the JSON file
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read questions file: %v", err)
	}

	// Parse JSON
	var rawQuestions []map[string]interface{} // Use a map instead of interface{} for better control
	err = json.Unmarshal(fileContent, &rawQuestions)
	if err != nil {
		return fmt.Errorf("failed to parse questions JSON: %v", err)
	}

	var questions []interface{}
	for _, q := range rawQuestions {
		// Convert the string id to primitive.ObjectID
		if idStr, ok := q["id"].(string); ok {
			id, err := primitive.ObjectIDFromHex(idStr)
			if err != nil {
				return fmt.Errorf("failed to convert id to ObjectID: %v", err)
			}
			q["id"] = id // Update the ID in the map to the ObjectID

			// Handle inputs: convert them to the proper format if necessary
			if inputs, ok := q["inputs"].([]interface{}); ok {
				// Verify inputs structure here (optional, depending on the JSON format)
				q["inputs"] = inputs // If needed, transform inputs to match the expected structure
			}
		}
		questions = append(questions, q)
	}

	// Insert questions into the collection
	_, err = QuestionsCollection.InsertMany(context.TODO(), questions)
	if err != nil {
		return fmt.Errorf("failed to insert questions into database: %v", err)
	}
	log.Println("Questions initialized successfully.")

	// Mark questions as initialized
	_, err = MetaCollection.InsertOne(context.TODO(), bson.M{"key": "questions_initialized", "timestamp": time.Now()})
	if err != nil {
		return fmt.Errorf("failed to mark questions as initialized: %v", err)
	}

	return nil
}

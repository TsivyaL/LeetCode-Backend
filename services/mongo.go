// services/db.go
package services

import (
    "context"
    "log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var QuestionsCollection *mongo.Collection

// SetupDB מחבר את האפליקציה למסד הנתונים
func SetupDB() error {
    clientOptions := options.Client().ApplyURI("mongodb://TsivyaL:MyName1sTsivya@localhost:27017")

    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        return err
    }

    // בדיקה שהחיבור הצליח
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        return err
    }

    MongoClient = client
    QuestionsCollection = client.Database("test").Collection("questions") // הגדרת הקולקציה
    log.Println("Connected to MongoDB!")
    return nil
}

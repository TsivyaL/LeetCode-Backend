package services

import (
    "context"
    "log"
    "os"
    "time"
    "os/exec"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var QuestionsCollection *mongo.Collection

// SetupDB מחבר את האפליקציה למסד הנתונים
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
        
        // אם החיבור נכשל, נוודא שהקונטיינר פעיל
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
    log.Println("Connected to MongoDB!")

    return nil
}

// startMongoContainer מוודא שהקונטיינר של MongoDB רץ
func startMongoContainer() error {
    cmd := exec.Command("docker", "start", "my_mongo_db")
    err := cmd.Run()
    if err != nil {
        return err
    }

    log.Println("MongoDB container started successfully")
    return nil
}

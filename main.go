package main

import (
    "context"
    "fmt"
    "log"
    //"time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    // יצירת הקשר והגדרת כתובת השרת
    clientOptions := options.Client().ApplyURI("mongodb://TsivyaL:MyName1sTsivy@@localhost:27017")

    // חיבור ל-MongoDB
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(context.TODO())

    // בדיקה שהחיבור פעיל
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal("לא ניתן להתחבר ל-MongoDB:", err)
    }

    fmt.Println("חיבור ל-MongoDB הצליח!")

    // גישה למסד נתונים ואוסף (Collection) לדוגמה
    collection := client.Database("testdb").Collection("users")

    // דוגמה להוספת מסמך לאוסף
    user := map[string]interface{}{
        "name": "John Doe",
        "email": "johndoe@example.com",
        "age": 29,
    }
    insertResult, err := collection.InsertOne(context.TODO(), user)
    if err != nil {
        log.Fatal("שגיאה בהוספת מסמך:", err)
    }

    fmt.Println("מסמך נוסף בהצלחה עם ID:", insertResult.InsertedID)
}

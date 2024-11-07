// services/services.go
package services

import (
    "context"
    "errors"
    "Backend/models"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

    //"log"
)

func FetchAllQuestions() ([]models.Question, error) {
    var questions []models.Question
    cursor, err := QuestionsCollection.Find(context.TODO(), bson.D{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var question models.Question
        if err := cursor.Decode(&question); err != nil {
            return nil, err
        }
        questions = append(questions, question)
    }

    return questions, nil
}

func FetchQuestionByID(id string) (models.Question, error) {
    var question models.Question
    err := QuestionsCollection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&question)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return question, errors.New("question not found")
        }
        return question, err
    }
    return question, nil
}

func AddQuestion(question models.Question) error {
    _, err := QuestionsCollection.InsertOne(context.TODO(), question)
    return err
}

func UpdateQuestion(id string, updatedQuestion models.Question) error {
    _, err := QuestionsCollection.UpdateOne(
        context.TODO(),
        bson.D{{"_id", id}},
        bson.D{{"$set", updatedQuestion}},
    )
    return err
}

func DeleteQuestion(id string) error {
    _, err := QuestionsCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}})
    return err
}

package services

import (
	"Backend/models"
	"context"
	"errors"

	//"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive" // Import for ObjectId conversion
)

// FetchAllQuestions fetches all questions from the database
func FetchAllQuestions() ([]models.Question, error) {
    var questions []models.Question
    cursor, err := QuestionsCollection.Find(context.TODO(), bson.D{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    // Loop through all questions in the cursor
    for cursor.Next(context.TODO()) {
        var question models.Question
        if err := cursor.Decode(&question); err != nil {
            return nil, err
        }
        questions = append(questions, question)
    }

    return questions, nil
}

// FetchQuestionByID fetches a single question by its ID
func FetchQuestionByID(id string) (models.Question, error) {
    var question models.Question

    // Convert string ID to ObjectId
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return question, errors.New("invalid ObjectId format")
    }

    // Query the question by its ObjectId
    err = QuestionsCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&question)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return question, errors.New("question not found")
        }
        return question, err
    }
    return question, nil
}

// AddQuestion adds a new question to the database
func AddQuestion(question models.Question) error {
    if question.Status == "" {
        question.Status = models.StatusNotStarted
    }
    // Convert the string ID to ObjectId if necessary
    if question.ID.IsZero() {
		question.ID = primitive.NewObjectID() // יצירת ObjectId חדש
	}

	_, err := QuestionsCollection.InsertOne(context.TODO(), question)
	return err
  
}

// UpdateQuestion updates an existing question in the database
func UpdateQuestionStatus(id string, status string, codeInProgress string) error {
    // Convert string ID to ObjectId
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return errors.New("invalid ObjectId format")
    }

    // Prepare update to only modify the "status" field
    update := bson.M{
        "$set": bson.M{
            "status": status, // Update the status field
            "code_in_progress": codeInProgress,
        },
    }

    // Perform the update operation
    _, err = QuestionsCollection.UpdateOne(
        context.TODO(),
        bson.M{"_id": objID}, // Use ObjectId for matching
        update,
    )
    return err
}


// DeleteQuestion deletes a question by its ID
func DeleteQuestion(id string) error {
    // Convert string ID to ObjectId
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return errors.New("invalid ObjectId format")
    }

    _, err = QuestionsCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
    return err
}
func DeleteAllQuestions() error {
	// Delete all documents in the collection
	_, err := QuestionsCollection.DeleteMany(context.TODO(), bson.D{})
	return err
}
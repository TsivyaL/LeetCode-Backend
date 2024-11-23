package controllers

import (
	"Backend/models"
	"Backend/services"
	"context"
	"log"
	"net/http"
	//"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

// SetupDB connects the app to the database
func SetupDB() {
	
    clientOptions := options.Client().ApplyURI("mongodb://TsivyaL:MyName1sTsivya@my_mongo_db:27017")
    client, err := mongo.Connect(context.Background(), clientOptions) // תחבר את האפליקציה למסד נתונים
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(context.Background(), nil) // בדוק אם החיבור הצליח
    if err != nil {
        log.Fatal("Cannot connect to MongoDB:", err)
    }

    MongoClient = client
    log.Println("Connected to MongoDB!")
}

// GetQuestions retrieves all questions from the database
func GetQuestions(c *gin.Context) {
	log.Println("in GetQuestions")
	questions, err := services.FetchAllQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}

// GetQuestion retrieves a single question by ID
func GetQuestion(c *gin.Context) {
    id := c.Param("id")
    question, err := services.FetchQuestionByID(id)
    if err != nil {
        log.Printf("Error fetching question with ID %s: %v", id, err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
        return
    }
    c.JSON(http.StatusOK, question)
}


// CreateQuestion creates a new question in the database
func CreateQuestion(c *gin.Context) {
	var question models.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if Id didn't receive, create a new one
	if question.ID.IsZero() {  // If the ObjectId is empty
		question.ID = primitive.NewObjectID()  // Create a new ID using ObjectId
	}

	if err := services.AddQuestion(question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Question created"})
}

// UpdateQuestion updates an existing question
func UpdateQuestionStatus(c *gin.Context) {
    id := c.Param("id")
    var status struct {
        Status string `json:"status"` 
    }

    if err := c.ShouldBindJSON(&status); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := services.UpdateQuestionStatus(id, status.Status); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Question status updated"})
}


// DeleteQuestion deletes a question by ID
func DeleteQuestion(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteQuestion(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Question deleted"})
}

// DeleteAllQuestions deletes all questions
func DeleteAllQuestions(c *gin.Context) {
	if err := services.DeleteAllQuestions(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All questions deleted"})
}

package controllers

import (
	"Backend/models"
	"Backend/services"
	"context" // ייבוא החבילה חסרה
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"         // ייבוא חבילה לחיבור ל-MongoDB
	"go.mongodb.org/mongo-driver/mongo/options" // ייבוא חבילה לחיבור ל-MongoDB
)

var MongoClient *mongo.Client

// SetupDB מחבר את האפליקציה למסד הנתונים
func SetupDB() {
	
    clientOptions := options.Client().ApplyURI("mongodb://TsivyaL:MyName1sTsivya@localhost:27017")
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
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

	// אם לא מסרת ID, צור אחד אוטומטית
	if question.ID.IsZero() {  // אם ה־ObjectId ריק
		question.ID = primitive.NewObjectID()  // יצירת ID חדש בעזרת ObjectId
	}

	if err := services.AddQuestion(question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Question created"})
}

// UpdateQuestion updates an existing question
func UpdateQuestion(c *gin.Context) {
	id := c.Param("id")
	var updatedQuestion models.Question
	if err := c.ShouldBindJSON(&updatedQuestion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateQuestion(id, updatedQuestion); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question updated"})
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
func DeleteAllQuestions(c *gin.Context) {
	if err := services.DeleteAllQuestions(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All questions deleted"})
}
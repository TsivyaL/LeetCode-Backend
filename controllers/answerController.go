package controllers

import (
	"Backend/models"
	"Backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SubmitAnswer handles answer submission
func SubmitAnswer(c *gin.Context) {
	questionID := c.Param("question_id") // קבל את question_id מה-URL

	var answer models.Answer
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer.QuestionID = questionID // עדכן את question_id בתשובה

	// בצע את הקוד ובדוק אם התשובה נכונה
	isCorrect, err := services.ExecuteAnswer(answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	answer.IsCorrect = isCorrect
	c.JSON(http.StatusOK, gin.H{"is_correct": isCorrect})
}

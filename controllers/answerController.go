package controllers

import (
	"Backend/models"
	"Backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SubmitAnswer handles answer submission
func SubmitAnswer(c *gin.Context) {
	questionID := c.Param("question_id") // Get the question_id from the URL

	var answer models.Answer
	if err := c.ShouldBindJSON(&answer); err != nil {
		// Return a more detailed error message if the input is invalid
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input format. Please check the fields and try again.",
			"details": err.Error(),
		})
		return
	}

	answer.QuestionID = questionID // Update the question_id in the answer

	// Execute the code and check if the answer is correct
	isCorrect, err := services.ExecuteAnswer(answer)
	if err != nil {
		// Return a more detailed error message if there was an issue executing the code
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An error occurred while evaluating the answer.",
			"details": err.Error(),
		})
		return
	}

	answer.IsCorrect = isCorrect
	c.JSON(http.StatusOK, gin.H{
		"is_correct": isCorrect,
		"message":    "Answer evaluated successfully.",
	})
}

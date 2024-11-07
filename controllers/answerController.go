package controllers

import (
	"Backend/models"
	"Backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SubmitAnswer ב-controller
func SubmitAnswer(c *gin.Context) {
	var answer models.Answer
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// הרץ את הקוד של התשובה
	isCorrect, err := services.ExecuteAnswer(answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// עדכן את התשובה אם היא נכונה או לא
	answer.IsCorrect = isCorrect

	// נניח שאנחנו שומרים את התשובה במסד נתונים
	c.JSON(http.StatusOK, gin.H{"is_correct": isCorrect})
}


package routes

import (
	"Backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/questions", controllers.GetQuestions)
	r.GET("/questions/:id", controllers.GetQuestion)
	r.POST("/questions", controllers.CreateQuestion)
	r.PUT("/questions/:id", controllers.UpdateQuestionStatus)
	r.DELETE("/questions/:id", controllers.DeleteQuestion)
	r.DELETE("/questions", controllers.DeleteAllQuestions)
}
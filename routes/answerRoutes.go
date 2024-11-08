package routes

import (
	"Backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetupAnswerRoutes(r *gin.Engine) {
	r.POST("/answers/:question_id", controllers.SubmitAnswer)
}

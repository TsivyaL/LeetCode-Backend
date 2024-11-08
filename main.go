package main

import (
	"fmt"
	"log"
	"Backend/services"
	"Backend/routes"

	"github.com/gin-gonic/gin"
	//"net/http"
)

func main() {
	r := gin.Default()

	// חיבור ל-DB
	// controllers.SetupDB()
	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "Welcome to the Questions API!",
	// 	})
	// })
	services.SetupDB()
	// הגדרת ה-Routes
	routes.SetupRoutes(r)
	routes.SetupAnswerRoutes(r)

	// הרצת השרת
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting server: ", err)
	}

	fmt.Println("Server is running on port 8080")
}
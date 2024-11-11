package main

import (
	"Backend/routes"
	"Backend/services"
	//"fmt"
	"log"
	//"os"
 	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	//"net/http"
)

func main() {

    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
 
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
	port := os.Getenv("PORT")
    if port == "" {
        port = "8080"  // ברירת מחדל
    }
    if err := r.Run(":" + port); err != nil {
        log.Fatal("Error starting server: ", err)
    }
    log.Println("Server is running on port", port)
}

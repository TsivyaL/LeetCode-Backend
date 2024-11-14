package main

import (
	"Backend/routes"
	"Backend/services"
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors" // Use this package for CORS middleware
)

func main() {
	log.Println("main 3")
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file") // Log an error if .env file is not loaded
	}

	r := gin.New()

	// Define CORS middleware using the gin-contrib/cors package
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Define allowed origin (e.g., frontend running on port 3000)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // Allowed HTTP methods for requests
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allowed headers for requests
		AllowCredentials: true, // Allow cookies and credentials
	}))

	r.Use(gin.Logger()) // Add logger middleware for request logging
	r.Use(gin.Recovery()) // Add recovery middleware to recover from panics and continue

	// Setup the database connection
	services.SetupDB()

	// Define the routes for the API
	routes.SetupRoutes(r)
	routes.SetupAnswerRoutes(r)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"  // Default port if not specified in .env file
	}
	// Run the server on the defined port
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Error starting server: ", err) // Log error if server fails to start
	}
	log.Println("Server is running on port", port) // Log the running server's port
}

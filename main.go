package main

import (
	"log"
	"net/http"
	"user-management-system/database"
	"user-management-system/routes"
)

func main() {
	// Connect to the database
	database.Connect()

	// Initialize routes
	router := routes.InitializeRoutes()

	// Start the server
	log.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

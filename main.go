package main

import (
	"log"
	"net/http"
	"user-management-system/database"
	"user-management-system/routes"

	"github.com/gorilla/handlers"
)

func main() {
	// Connect to the database
	database.Connect()

	// Initialize routes
	router := routes.InitializeRoutes()

	// CORS handler to allow requests from different origins
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins (you can specify specific origins like http://localhost:3000)
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Wrap the router with CORS middleware
	http.Handle("/", cors(router))

	// Start the server
	log.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

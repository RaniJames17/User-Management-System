package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"user-management-system/database"
	"user-management-system/models"
	"user-management-system/utils"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Debug: Print incoming request body
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Debug: Log decoded user details
	log.Printf("Decoded user: %+v", user)

	// Validate the input
	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "All fields (name, email, password) are required", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	// Insert the user into the database
	query := "INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)"
	_, err = database.ExecuteQuery(query, user.Name, user.Email, user.Password)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

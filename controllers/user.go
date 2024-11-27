package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"user-management-system/database"
	"user-management-system/models"
	"user-management-system/utils"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate input
	if credentials.Email == "" || credentials.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Fetch user from the database
	var user models.User
	query := "SELECT id, name, email, password_hash FROM users WHERE email = ?"
	row := database.DB.QueryRow(query, credentials.Email)
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Verify the password
	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate a session token (for simplicity, we'll use a simple string here)
	token := utils.GenerateToken(user.ID)

	// Return success response with token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Sign-In successful",
		"token":   token,
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

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

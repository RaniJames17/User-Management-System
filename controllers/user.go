package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
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

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email string `json:"email"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate email
	if request.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Find user by email
	var user models.User
	query := "SELECT id, email FROM users WHERE email = ?"
	err = database.DB.QueryRow(query, request.Email).Scan(&user.ID, &user.Email)
	if err != nil {
		// Don't disclose that the email is incorrect
		http.Error(w, "If the email is valid, you will receive a reset link.", http.StatusOK)
		return
	}

	// Generate password reset token
	resetToken := utils.GenerateToken(user.ID)

	// Store the reset token in the database
	query = "INSERT INTO password_resets (user_id, token, created_at) VALUES (?, ?, ?)"
	_, err = database.ExecuteQuery(query, user.ID, resetToken, time.Now())
	if err != nil {
		http.Error(w, "Failed to create password reset token", http.StatusInternalServerError)
		return
	}

	// For now, log the reset token (for testing purposes)
	log.Printf("Password reset token for user %s: %s", user.Email, resetToken)

	// Send response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "If the email is valid, you will receive a reset link.",
	})
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ResetToken  string `json:"reset_token"`
		NewPassword string `json:"new_password"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate input
	if request.ResetToken == "" || request.NewPassword == "" {
		http.Error(w, "Reset token and new password are required", http.StatusBadRequest)
		return
	}

	// Find the reset token in the database
	var reset models.PasswordReset
	var createdAtString string

	log.Printf("Received reset token: %s", request.ResetToken)

	query := "SELECT user_id, token, created_at FROM password_resets WHERE token = ?"
	err = database.DB.QueryRow(query, request.ResetToken).Scan(&reset.UserID, &reset.Token, &createdAtString)
	if err != nil {
		log.Printf("Error fetching reset token from DB: %v", err)
		http.Error(w, "Invalid or expired reset token", http.StatusBadRequest)
		return
	}

	// Parse the created_at string into a time.Time object
	// Use the correct format for parsing the timestamp
	reset.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtString)
	if err != nil {
		log.Printf("Error parsing created_at: %v", err)
		http.Error(w, "Invalid reset token timestamp", http.StatusBadRequest)
		return
	}

	// Check if the reset token is expired (24 hours)
	currentTime := time.Now().UTC() // Get current time in UTC
	if currentTime.After(reset.CreatedAt.Add(24 * time.Hour)) {
		http.Error(w, "Reset token has expired", http.StatusBadRequest)
		return
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Update the user's password
	updateQuery := "UPDATE users SET password_hash = ? WHERE id = ?"
	_, err = database.ExecuteQuery(updateQuery, hashedPassword, reset.UserID)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	// Optionally, remove the reset token from the database after successful reset
	deleteQuery := "DELETE FROM password_resets WHERE token = ?"
	_, err = database.ExecuteQuery(deleteQuery, request.ResetToken)
	if err != nil {
		log.Println("Failed to delete reset token:", err)
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password has been successfully reset.",
	})
}

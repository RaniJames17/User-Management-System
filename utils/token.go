package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"strings"
	"time"
)

// GenerateToken creates a simple session token with a timestamp (valid for 24 hours)
func GenerateToken(userID int) string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	timestamp := time.Now().Unix() // Current timestamp in seconds
	return fmt.Sprintf("%x-%d-%d", bytes, userID, timestamp)
}

// ValidateToken checks if a token is valid and not expired
func ValidateToken(token string) bool {
	parts := strings.Split(token, "-")
	if len(parts) != 3 {
		log.Println("Invalid token format")
		return false
	}

	// Extract and parse the timestamp part
	timestamp := parts[2]
	var ts int64
	_, err := fmt.Sscanf(timestamp, "%d", &ts)
	if err != nil {
		log.Printf("Error parsing timestamp: %v", err)
		return false
	}

	// Check if the token is expired (valid for 24 hours)
	expiry := time.Unix(ts, 0).Add(24 * time.Hour)
	log.Printf("Token expires at: %v", expiry)

	return time.Now().Before(expiry)
}

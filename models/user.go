package models

import (
	"time"
)

type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	ResetToken  string    `json:"reset_token,omitempty"`
	TokenExpiry time.Time `json:"token_expiry,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

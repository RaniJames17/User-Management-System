package models

import (
	"time"
)

type PasswordReset struct {
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

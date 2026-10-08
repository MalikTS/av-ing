package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Phone        string    `json:"phone,omitempty"`
	City         string    `json:"city,omitempty"`
	IsAdmin      bool      `json:"is_admin"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	City     string `json:"city"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

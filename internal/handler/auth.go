package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"av-ing/internal/config"
	"av-ing/internal/database"
	"av-ing/internal/middleware"
	"av-ing/internal/models"
)

var cfg *config.Config

func InitAuthHandler(c *config.Config) {
	cfg = c
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := decodeJSON(w, r, &req); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	address, err := mail.ParseAddress(req.Email)
	if err != nil || address.Address != req.Email {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}
	if len(req.Username) < 2 || len(req.Username) > 50 {
		http.Error(w, "Username must contain 2 to 50 characters", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 8 || len(req.Password) > 72 {
		http.Error(w, "Password must contain 8 to 72 characters", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO users (username, email, password_hash, phone, city) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	var userID string
	var createdAt time.Time
	err = database.DB.QueryRow(query, req.Username, req.Email, string(hashedPassword), req.Phone, req.City).Scan(&userID, &createdAt)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "User registered successfully",
		"id":      userID,
	})

}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := decodeJSON(w, r, &req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	var storedHash string
	query := `SELECT id, username, email, password_hash, is_admin FROM users WHERE email = $1`
	err := database.DB.QueryRow(query, req.Email).Scan(&user.ID, &user.Username, &user.Email, &storedHash, &user.IsAdmin)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"token":   tokenString,
		"message": "Login successful",
	})
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	var user models.User
	query := `SELECT id, username, email, COALESCE(phone, ''), COALESCE(city, ''), is_admin FROM users WHERE id = $1`
	err := database.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.City, &user.IsAdmin)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"

	"av-ing/internal/database"
	"av-ing/internal/middleware"
)

type ReviewResponse struct {
	ID        string    `json:"id"`
	FromUser  string    `json:"from_user"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateReview(w http.ResponseWriter, r *http.Request) {
	fromUserID := middleware.GetUserIDFromContext(r)
	adID := chi.URLParam(r, "id")

	var req struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}
	req.Comment = strings.TrimSpace(req.Comment)

	var toUserID string
	err := database.DB.QueryRow("SELECT user_id FROM ads WHERE id = $1 AND status = 'active'", adID).Scan(&toUserID)
	if err == sql.ErrNoRows {
		http.Error(w, "Ad not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if fromUserID == toUserID {
		http.Error(w, "Cannot review your own ad", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO reviews (from_user_id, to_user_id, ad_id, rating, comment) VALUES ($1, $2, $3, $4, $5)`
	_, err = database.DB.Exec(query, fromUserID, toUserID, adID, req.Rating, req.Comment)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			http.Error(w, "You have already reviewed this ad", http.StatusConflict)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "Review submitted"})
}

func GetAdReviews(w http.ResponseWriter, r *http.Request) {
	adID := chi.URLParam(r, "id")

	rows, err := database.DB.Query(`
		SELECT r.id, u.username, r.rating, r.comment, r.created_at 
		FROM reviews r
		JOIN users u ON r.from_user_id = u.id
		JOIN ads a ON r.ad_id = a.id
		WHERE r.ad_id = $1 AND a.status = 'active'
		ORDER BY r.created_at DESC`, adID)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reviews []ReviewResponse
	for rows.Next() {
		var rev ReviewResponse
		if err := rows.Scan(&rev.ID, &rev.FromUser, &rev.Rating, &rev.Comment, &rev.CreatedAt); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		reviews = append(reviews, rev)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var avgRating sql.NullFloat64
	err = database.DB.QueryRow(`
		SELECT AVG(r.rating)
		FROM reviews r
		JOIN ads a ON r.ad_id = a.id
		WHERE r.ad_id = $1 AND a.status = 'active'`, adID).Scan(&avgRating)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{
		"reviews":        reviews,
		"average_rating": 0.0,
	}
	if avgRating.Valid {
		result["average_rating"] = avgRating.Float64
	}

	writeJSON(w, http.StatusOK, result)
}

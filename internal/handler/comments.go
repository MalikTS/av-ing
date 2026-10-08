package handler

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"av-ing/internal/database"
	"av-ing/internal/middleware"
)

type Comment struct {
	ID        string    `json:"id"`
	AdID      string    `json:"ad_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	adID := chi.URLParam(r, "id")

	var req struct {
		Content string `json:"content"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO comments (ad_id, user_id, content)
		SELECT $1, $2, $3
		WHERE EXISTS (SELECT 1 FROM ads WHERE id = $1 AND status = 'active')
		RETURNING id, created_at`
	var commentID string
	var createdAt time.Time

	err := database.DB.QueryRow(query, adID, userID, req.Content).Scan(&commentID, &createdAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Ad not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"id":      commentID,
		"message": "Comment added",
	})
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	adID := chi.URLParam(r, "id")

	rows, err := database.DB.Query(`
		SELECT c.id, c.user_id, u.username, c.content, c.created_at 
		FROM comments c
		JOIN users u ON c.user_id = u.id 
		JOIN ads a ON c.ad_id = a.id
		WHERE c.ad_id = $1 AND a.status = 'active' AND c.is_deleted = FALSE
		ORDER BY c.created_at ASC`, adID)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, comments)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	commentID := chi.URLParam(r, "comment_id")

	var ownerID string
	err := database.DB.QueryRow("SELECT user_id FROM comments WHERE id = $1", commentID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if ownerID != userID && !middleware.IsAdminFromContext(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_, err = database.DB.Exec("UPDATE comments SET is_deleted = TRUE WHERE id = $1", commentID)
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Comment deleted"})
}

package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"

	"av-ing/internal/database"
	"av-ing/internal/middleware"
	"av-ing/internal/models"
)

func AddToFavorites(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	adID := chi.URLParam(r, "id")

	result, err := database.DB.Exec(`
		INSERT INTO favorites (user_id, ad_id)
		SELECT $1, $2
		WHERE EXISTS (SELECT 1 FROM ads WHERE id = $2 AND status = 'active')`, userID, adID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			http.Error(w, "Already in favorites", http.StatusConflict)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Ad not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "Added to favorites"})
}

func RemoveFromFavorites(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	adID := chi.URLParam(r, "id")

	_, err := database.DB.Exec("DELETE FROM favorites WHERE user_id = $1 AND ad_id = $2", userID, adID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Removed from favorites"})
}

func GetMyFavorites(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	rows, err := database.DB.Query(`
		SELECT a.id, a.title, a.price, a.currency, a.status, a.created_at 
		FROM ads a 
		JOIN favorites f ON a.id = f.ad_id 
		WHERE f.user_id = $1 AND a.status = 'active'`, userID)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ads []models.Ad
	for rows.Next() {
		var ad models.Ad
		if err := rows.Scan(&ad.ID, &ad.Title, &ad.Price, &ad.Currency, &ad.Status, &ad.CreatedAt); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		ads = append(ads, ad)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, ads)
}

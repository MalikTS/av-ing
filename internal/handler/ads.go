package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"av-ing/internal/database"
	"av-ing/internal/middleware"
	"av-ing/internal/models"
)

func CreateAd(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	var req models.CreateAdRequest
	if err := decodeJSON(w, r, &req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" || strings.TrimSpace(req.CategoryID) == "" {
		http.Error(w, "Title and category_id are required", http.StatusBadRequest)
		return
	}
	if req.Price < 0 {
		http.Error(w, "Price cannot be negative", http.StatusBadRequest)
		return
	}
	if len(req.Title) > 200 {
		http.Error(w, "Title must not exceed 200 characters", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Currency) == "" {
		req.Currency = "RUB"
	}
	if strings.TrimSpace(req.Condition) == "" {
		req.Condition = "used"
	}
	var locationID any
	if strings.TrimSpace(req.LocationID) != "" {
		locationID = req.LocationID
	}

	query := `
		INSERT INTO ads (title, description, price, currency, category_id, user_id, location_id, address, condition, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'active')
		RETURNING id, created_at, updated_at`

	var ad models.Ad
	err := database.DB.QueryRow(query,
		req.Title, req.Description, req.Price, req.Currency,
		req.CategoryID, userID, locationID, req.Address, req.Condition,
	).Scan(&ad.ID, &ad.CreatedAt, &ad.UpdatedAt)

	if err != nil {
		http.Error(w, "Failed to create ad", http.StatusInternalServerError)
		return
	}

	ad.UserID = userID
	ad.Title = req.Title
	ad.Description = req.Description
	ad.Price = req.Price
	ad.Currency = req.Currency
	ad.CategoryID = req.CategoryID
	ad.LocationID = req.LocationID
	ad.Condition = req.Condition
	ad.Status = "active"

	writeJSON(w, http.StatusCreated, ad)
}

func GetAds(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")
	priceMin := r.URL.Query().Get("price_min")
	priceMax := r.URL.Query().Get("price_max")
	categoryID := r.URL.Query().Get("category_id")
	city := r.URL.Query().Get("city")
	sortBy := r.URL.Query().Get("sort_by")

	page, err := parsePositiveInt(pageStr, "page", 1, 1_000_000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit, err := parsePositiveInt(limitStr, "limit", 20, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset := (page - 1) * limit

	query := `
		SELECT a.id, a.title, a.description, a.price, a.currency, a.category_id, c.name as category_name,
		       a.user_id, COALESCE(a.location_id::text, ''), l.city, a.condition, a.status, a.views_count,
		       a.is_promoted, a.created_at, a.updated_at, a.published_at
		FROM ads a
		LEFT JOIN categories c ON a.category_id = c.id
		LEFT JOIN locations l ON a.location_id = l.id
		WHERE a.status = 'active'`

	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += " AND (a.title ILIKE $" + fmt.Sprint(argIndex) + " OR a.description ILIKE $" + fmt.Sprint(argIndex) + ")"
		args = append(args, "%"+search+"%")
		argIndex++
	}
	if categoryID != "" {
		query += " AND a.category_id = $" + fmt.Sprint(argIndex)
		args = append(args, categoryID)
		argIndex++
	}
	if priceMin != "" {
		value, err := parseNonNegativeFloat(priceMin, "price_min")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		query += " AND a.price >= $" + fmt.Sprint(argIndex)
		args = append(args, value)
		argIndex++
	}
	if priceMax != "" {
		value, err := parseNonNegativeFloat(priceMax, "price_max")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		query += " AND a.price <= $" + fmt.Sprint(argIndex)
		args = append(args, value)
		argIndex++
	}
	if city != "" {
		query += " AND l.city = $" + fmt.Sprint(argIndex)
		args = append(args, city)
		argIndex++
	}

	switch sortBy {
	case "cheapest":
		query += " ORDER BY a.price ASC"
	case "expensive":
		query += " ORDER BY a.price DESC"
	default:
		query += " ORDER BY a.created_at DESC"
	}

	query += " LIMIT $" + fmt.Sprint(argIndex) + " OFFSET $" + fmt.Sprint(argIndex+1)
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ads []models.Ad
	for rows.Next() {
		var ad models.Ad
		var category sql.NullString
		var city sql.NullString
		var publishedAt sql.NullTime

		err := rows.Scan(&ad.ID, &ad.Title, &ad.Description, &ad.Price, &ad.Currency,
			&ad.CategoryID, &category, &ad.UserID, &ad.LocationID, &city,
			&ad.Condition, &ad.Status, &ad.ViewsCount, &ad.IsPromoted, &ad.CreatedAt,
			&ad.UpdatedAt, &publishedAt)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if category.Valid {
			ad.CategoryName = category.String
		}
		if city.Valid {
			ad.City = city.String
		}
		if publishedAt.Valid {
			ad.PublishedAt = &publishedAt.Time
		}
		ads = append(ads, ad)
	}

	// countQuery := "SELECT COUNT(*) FROM ads WHERE status = 'active'"

	if err := rows.Err(); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, ads)
}

func GetAd(w http.ResponseWriter, r *http.Request) {
	adID := chi.URLParam(r, "id")

	var ad models.Ad
	var category sql.NullString
	var city sql.NullString
	var publishedAt sql.NullTime

	query := `
		SELECT a.id, a.title, a.description, a.price, a.currency, a.category_id, c.name as category_name,
		       a.user_id, COALESCE(a.location_id::text, ''), l.city, a.condition, a.status, a.views_count,
		       a.is_promoted, a.created_at, a.updated_at, a.published_at
		FROM ads a
		LEFT JOIN categories c ON a.category_id = c.id
		LEFT JOIN locations l ON a.location_id = l.id
		WHERE a.id = $1 AND a.status = 'active'`

	err := database.DB.QueryRow(query, adID).Scan(
		&ad.ID, &ad.Title, &ad.Description, &ad.Price, &ad.Currency,
		&ad.CategoryID, &category, &ad.UserID, &ad.LocationID, &city,
		&ad.Condition, &ad.Status, &ad.ViewsCount, &ad.IsPromoted, &ad.CreatedAt,
		&ad.UpdatedAt, &publishedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Ad not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if category.Valid {
		ad.CategoryName = category.String
	}
	if city.Valid {
		ad.City = city.String
	}
	if publishedAt.Valid {
		ad.PublishedAt = &publishedAt.Time
	}

	_, _ = database.DB.ExecContext(r.Context(), "UPDATE ads SET views_count = views_count + 1 WHERE id = $1", adID)

	writeJSON(w, http.StatusOK, ad)
}

func UpdateAd(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	adID := chi.URLParam(r, "id")

	var ownerID string
	err := database.DB.QueryRow("SELECT user_id FROM ads WHERE id = $1", adID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		http.Error(w, "Ad not found", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("DB Error in UpdateAd: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if ownerID != userID {
		if !middleware.IsAdminFromContext(r) {
			log.Printf("Forbidden: User %s tried to update Ad %s owned by %s", userID, adID, ownerID)
			http.Error(w, "Forbidden: You are not the owner of this ad", http.StatusForbidden)
			return
		}
	}
	var req models.UpdateAdRequest
	if err := decodeJSON(w, r, &req); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Title != nil {
		*req.Title = strings.TrimSpace(*req.Title)
		if *req.Title == "" || len(*req.Title) > 200 {
			http.Error(w, "Title must contain 1 to 200 characters", http.StatusBadRequest)
			return
		}
		updates = append(updates, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *req.Title)
		argIndex++
	}
	if req.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.Price != nil {
		if *req.Price < 0 {
			http.Error(w, "Price cannot be negative", http.StatusBadRequest)
			return
		}
		updates = append(updates, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, *req.Price)
		argIndex++
	}
	if req.CategoryID != nil {
		if strings.TrimSpace(*req.CategoryID) == "" {
			http.Error(w, "Category_id cannot be empty", http.StatusBadRequest)
			return
		}
		updates = append(updates, fmt.Sprintf("category_id = $%d", argIndex))
		args = append(args, *req.CategoryID)
		argIndex++
	}
	if req.Condition != nil {
		updates = append(updates, fmt.Sprintf("condition = $%d", argIndex))
		args = append(args, *req.Condition)
		argIndex++
	}
	if req.Status != nil {
		if *req.Status != "active" && *req.Status != "archived" {
			http.Error(w, "Status must be active or archived", http.StatusBadRequest)
			return
		}
		updates = append(updates, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	if len(updates) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	updates = append(updates, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	args = append(args, adID)

	query := fmt.Sprintf("UPDATE ads SET %s WHERE id = $%d",
		strings.Join(updates, ", "), argIndex+1)

	_, err = database.DB.Exec(query, args...)
	if err != nil {
		http.Error(w, "Failed to update ad", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Ad updated"})
}

func DeleteAd(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	adID := chi.URLParam(r, "id")

	var ownerID string
	err := database.DB.QueryRow("SELECT user_id FROM ads WHERE id = $1", adID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		http.Error(w, "Ad not found", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("DB Error in DeleteAd: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if ownerID != userID && !middleware.IsAdminFromContext(r) {
		log.Printf("Forbidden: User %s tried to delete Ad %s owned by %s", userID, adID, ownerID)
		http.Error(w, "Forbidden: You are not the owner of this ad", http.StatusForbidden)
		return
	}

	_, err = database.DB.Exec("UPDATE ads SET status = 'archived', updated_at = $2 WHERE id = $1", adID, time.Now())
	if err != nil {
		http.Error(w, "Failed to delete ad", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Ad deleted"})
}

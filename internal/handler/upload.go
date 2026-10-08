package handler

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"

	"av-ing/internal/database"
	"av-ing/internal/middleware"
)

const (
	uploadDir      = "web/uploads/ads/"
	maxFileSize    = 5 * 1024 * 1024
	maxImages      = 10
	maxRequestSize = maxImages*maxFileSize + 1024*1024
)

func init() {
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		panic("cannot create upload directory: " + err.Error())
	}
}

func UploadAdImages(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	adID := chi.URLParam(r, "id")

	var ownerID string
	err := database.DB.QueryRow("SELECT user_id FROM ads WHERE id = $1", adID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		http.Error(w, "Ad not found", http.StatusNotFound)
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

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	if err := r.ParseMultipartForm(maxRequestSize); err != nil {
		http.Error(w, "Request is too large or has invalid multipart data", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		http.Error(w, "No files provided", http.StatusBadRequest)
		return
	}
	if len(files) > maxImages {
		http.Error(w, "Too many files; maximum is 10", http.StatusBadRequest)
		return
	}

	extensions := make([]string, len(files))
	allowedTypes := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/webp": ".webp",
	}
	for index, fileHeader := range files {
		if fileHeader.Size <= 0 || fileHeader.Size > maxFileSize {
			http.Error(w, "Each file must be no larger than 5 MB", http.StatusBadRequest)
			return
		}
		src, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		buf := make([]byte, 512)
		n, readErr := src.Read(buf)
		_ = src.Close()
		if readErr != nil && readErr != io.EOF {
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			return
		}
		ext, valid := allowedTypes[http.DetectContentType(buf[:n])]
		if !valid {
			http.Error(w, "Invalid file content type", http.StatusBadRequest)
			return
		}
		extensions[index] = ext
	}

	tx, err := database.DB.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var hasImages bool
	if err = tx.QueryRow("SELECT EXISTS (SELECT 1 FROM ad_images WHERE ad_id = $1)", adID).Scan(&hasImages); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	savedPaths := make([]string, 0, len(files))
	uploadedURLs := make([]string, 0, len(files))
	cleanup := func() {
		for _, path := range savedPaths {
			_ = os.Remove(path)
		}
	}

	for index, fileHeader := range files {
		src, err := fileHeader.Open()
		if err != nil {
			cleanup()
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}

		uniqueName := fmt.Sprintf("%s_%d_%d%s", adID, time.Now().UnixNano(), index, extensions[index])
		dstPath := filepath.Join(uploadDir, uniqueName)
		dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			_ = src.Close()
			cleanup()
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		written, copyErr := io.Copy(dst, io.LimitReader(src, maxFileSize+1))
		closeErr := dst.Close()
		_ = src.Close()
		if copyErr != nil || closeErr != nil || written > maxFileSize {
			_ = os.Remove(dstPath)
			cleanup()
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		savedPaths = append(savedPaths, dstPath)

		imageURL := "/uploads/ads/" + uniqueName
		isMain := !hasImages && index == 0
		_, err = tx.Exec(
			"INSERT INTO ad_images (ad_id, image_url, is_main) VALUES ($1, $2, $3)",
			adID, imageURL, isMain,
		)
		if err != nil {
			cleanup()
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		uploadedURLs = append(uploadedURLs, imageURL)
	}

	if err = tx.Commit(); err != nil {
		cleanup()
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message": "Images uploaded",
		"urls":    uploadedURLs,
	})
}

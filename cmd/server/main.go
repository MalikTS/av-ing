package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"av-ing/internal/config"
	"av-ing/internal/database"
	"av-ing/internal/handler"
	"av-ing/internal/middleware"
)

func main() {
	cfg := config.LoadConfig()

	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer database.DB.Close()

	handler.InitAuthHandler(cfg)
	middleware.InitAuthMiddleware(cfg)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("Welcome to av-ing API!"))
	})

	// Публичные маршруты
	r.Get("/ads", handler.GetAds)
	r.Get("/ads/{id}", handler.GetAd)
	r.Get("/ads/{id}/comments", handler.GetComments)
	r.Get("/ads/{id}/reviews", handler.GetAdReviews)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/me", handler.GetMe)

		// Маршруты для объявлений
		r.Post("/ads", handler.CreateAd)
		r.Put("/ads/{id}", handler.UpdateAd)
		r.Delete("/ads/{id}", handler.DeleteAd)
		r.Post("/ads/{id}/images", handler.UploadAdImages)

		// Комментарии
		r.Post("/ads/{id}/comments", handler.CreateComment)
		r.Delete("/comments/{comment_id}", handler.DeleteComment)

		// Избранное
		r.Post("/ads/{id}/favorite", handler.AddToFavorites)
		r.Delete("/ads/{id}/favorite", handler.RemoveFromFavorites)
		r.Get("/me/favorites", handler.GetMyFavorites)

		// Рейтинг
		r.Post("/ads/{id}/reviews", handler.CreateReview)
	})

	r.Handle("/uploads/*",
		http.StripPrefix("/uploads/",
			http.FileServer(http.Dir("./web/uploads/"))))

	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("server is running on http://localhost:%s", cfg.ServerPort)
	log.Fatal(server.ListenAndServe())
}

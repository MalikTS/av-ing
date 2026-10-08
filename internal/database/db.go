package database

import (
	"database/sql"
	"fmt"
	"time"

	"av-ing/internal/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(cfg *config.Config) error {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err = DB.Ping(); err != nil {
		_ = DB.Close()
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}

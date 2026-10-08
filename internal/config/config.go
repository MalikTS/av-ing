package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	JWTSecret  string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		ServerPort: os.Getenv("SERVER_PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}
	return cfg
}

func (c *Config) Validate() error {
	if c == nil {
		return errors.New("config is nil")
	}

	required := map[string]string{
		"DB_HOST": c.DBHost, "DB_PORT": c.DBPort, "DB_USER": c.DBUser,
		"DB_PASSWORD": c.DBPassword, "DB_NAME": c.DBName, "JWT_SECRET": c.JWTSecret,
	}
	for name, value := range required {
		if value == "" {
			return fmt.Errorf("required environment variable %s is empty", name)
		}
	}
	if len(c.JWTSecret) < 32 {
		return errors.New("JWT_SECRET must contain at least 32 characters")
	}
	return nil
}

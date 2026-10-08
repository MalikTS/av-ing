package models

import "time"

type Ad struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Price        float64    `json:"price"`
	Currency     string     `json:"currency"`
	CategoryID   string     `json:"category_id"`
	CategoryName string     `json:"category_name,omitempty"`
	UserID       string     `json:"user_id"`
	LocationID   string     `json:"location_id,omitempty"`
	City         string     `json:"city,omitempty"`
	Condition    string     `json:"condition"`
	Status       string     `json:"status"`
	ViewsCount   int        `json:"views_count"`
	IsPromoted   bool       `json:"is_promoted"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	PublishedAt  *time.Time `json:"published_at,omitempty"`
}

type CreateAdRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	CategoryID  string  `json:"category_id"`
	LocationID  string  `json:"location_id"`
	Address     string  `json:"address"`
	Condition   string  `json:"condition"`
}

type UpdateAdRequest struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	CategoryID  *string  `json:"category_id,omitempty"`
	Condition   *string  `json:"condition,omitempty"`
	Status      *string  `json:"status,omitempty"`
}

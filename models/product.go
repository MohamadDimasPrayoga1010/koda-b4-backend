package models

import "time"

type Product struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	BasePrice   float64        `json:"base_price"`
	Stock       int            `json:"stock"`
	CategoryID  int64          `json:"category_id"`
	VariantID   int64          `json:"variant_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty"`
	Images      []ProductImage `json:"images,omitempty"`
	Sizes       []Size         `json:"sizes,omitempty"`
}

type ProductResponse struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	BasePrice   float64        `json:"base_price"`
	Stock       int            `json:"stock"`
	Images      []ProductImage `json:"images,omitempty"`
	Sizes       []Size         `json:"sizes,omitempty"`
	CategoryID  int64          `json:"category_id"`
	VariantID   int64          `json:"variant_id"`
	CreatedAt   time.Time      `json:"created_at"`
}

type ProductImage struct {
	ProductID int64      `json:"product_id"`
	Image     string     `json:"image"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}


type Size struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	AdditionalPrice float64 `json:"additional_price"`
}

type ProductRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	BasePrice   float64  `json:"base_price" binding:"required"`
	Stock       int      `json:"stock" binding:"required"`
	CategoryID  int64    `json:"category_id"`
	VariantID   int64    `json:"variant_id"`
	Images      []string `json:"images"` 
	Sizes       []int64  `json:"sizes"`  
}
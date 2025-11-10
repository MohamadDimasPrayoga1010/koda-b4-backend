package models

import (
	"mime/multipart"
	"time"
)

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
	CategoryID  int64          `json:"category_id"`
	VariantID   int64          `json:"variant_id"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
	Images      []ProductImage `json:"images,omitempty"`
	Sizes       []Size         `json:"sizes,omitempty"`
}


type ProductImage struct {
	ProductID int64      `json:"product_id"`
	Image     string     `json:"image"` 
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}


type Size struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name" form:"name"`
	AdditionalPrice float64 `json:"additional_price" form:"additional_price"`
}

type ProductRequest struct {
	Title       string                   `form:"title" binding:"required"`       
	Description string                   `form:"description"`                   
	BasePrice   float64                  `form:"base_price" binding:"required"` 
	Stock       int                      `form:"stock" binding:"required"`      
	CategoryID  int64                    `form:"category_id"`                   
	VariantID   int64                    `form:"variant_id"`                    
	Sizes       []int64                  `form:"sizes"`                         
	Images      []*multipart.FileHeader  `form:"images"`                        
}

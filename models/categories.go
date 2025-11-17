package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}


func CreateCategory(db *pgxpool.Pool, name string) (Category, error) {
	var category Category
	query := `
		INSERT INTO categories (name, created_at, updated_at)
		VALUES ($1, now(), now())
		RETURNING id, name, created_at, updated_at
	`

	err := db.QueryRow(context.Background(), query, name).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return Category{}, err
	}

	return category, nil
}

func GetAllCategories(db *pgxpool.Pool) ([]Category, error) {
	rows, err := db.Query(context.Background(),
		`SELECT id, name, created_at, updated_at FROM categories ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func GetCategoryByID(db *pgxpool.Pool, id int64) (*Category, error) {
	var c Category
	err := db.QueryRow(
		context.Background(),
		`SELECT id, name, created_at, updated_at FROM categories WHERE id=$1`,
		id,
	).Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

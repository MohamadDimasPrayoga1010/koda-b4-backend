package models

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct{
	ID        int64     `json:"id" db:"id"`
	Fullname  string    `json:"fullname" db:"fullname"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Token    string `json:"token,omitempty"`
}


type ProfileRequest struct {
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
}

type ProfileUser struct {
	ID        int64     `json:"id" db:"id"`
	Image     *string   `json:"image,omitempty" db:"image"`
	Phone     *string   `json:"phone,omitempty" db:"phone"`
	Address   *string   `json:"address,omitempty" db:"address"`
	UserID    int64     `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func UpdateProfile(db *pgxpool.Pool, userID int64, phone, address string, fileHeader *multipart.FileHeader) (*ProfileUser, error) {
	ctx := context.Background()

	var imagePath *string
	if fileHeader != nil {
		if fileHeader.Size > 2*1024*1024 {
			return nil, errors.New("file size exceeds 2MB")
		}
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return nil, errors.New("file type must be jpg, jpeg, or png")
		}

		filename := fmt.Sprintf("uploads/%d_%d%s", userID, time.Now().Unix(), ext)
		if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
			return nil, err
		}

		src, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		dst, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return nil, err
		}

		imagePath = &filename
	}

	var existingID int64
	err := db.QueryRow(ctx, `SELECT id FROM profile WHERE user_id=$1`, userID).Scan(&existingID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if existingID > 0 {

		_, err := db.Exec(ctx, `
			UPDATE profile
			SET phone = COALESCE(NULLIF($1, ''), phone),
				address = COALESCE(NULLIF($2, ''), address),
				image = COALESCE($3, image),
				updated_at = NOW()
			WHERE user_id = $4
		`, phone, address, imagePath, userID)
		if err != nil {
			return nil, err
		}
	} else {

		_, err := db.Exec(ctx, `
			INSERT INTO profile (phone, address, image, user_id)
			VALUES ($1,$2,$3,$4)
		`, phone, address, imagePath, userID)
		if err != nil {
			return nil, err
		}
	}

	var profile ProfileUser
	err = db.QueryRow(ctx, `
		SELECT id, phone, address, image, user_id, created_at, updated_at
		FROM profile
		WHERE user_id=$1
	`, userID).Scan(&profile.ID, &profile.Phone, &profile.Address, &profile.Image, &profile.UserID, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}


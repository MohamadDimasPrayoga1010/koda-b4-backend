package models

import (
	"coffeeder-backend/libs"
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




type ProfileRequest struct {
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
}

type ProfileUser struct {
	ID        int64      `json:"id" db:"id"`
	Image     *string    `json:"image,omitempty" db:"image"`
	Phone     *string    `json:"phone,omitempty" db:"phone"`
	Address   *string    `json:"address,omitempty" db:"address"`
	UserID    int64      `json:"userId" db:"user_id"`
	Fullname  *string    `json:"fullname,omitempty" db:"fullname"`
	Email     *string    `json:"email,omitempty" db:"email"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
}

type ProfileResponse struct {
	ID        int64      `json:"id"`
	Fullname  string     `json:"fullname"`
	Email     string     `json:"email"`
	Image     *string    `json:"image,omitempty"`
	Phone     *string    `json:"phone,omitempty"`
	Address   *string    `json:"address,omitempty"`
	UserID    int64      `json:"userId"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}


func UpdateProfile(db *pgxpool.Pool, userID int64, phone, address, fullname, email string, fileHeader *multipart.FileHeader) (ProfileResponse, error) {
	ctx := context.Background()
	var imagePath *string

	if fileHeader != nil {
		if fileHeader.Size > 2*1024*1024 {
			return ProfileResponse{}, errors.New("file size exceeds 2MB")
		}

		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return ProfileResponse{}, errors.New("file type must be jpg, jpeg, or png")
		}

		if os.Getenv("CLOUDINARY_API_KEY") != "" {
			url, err := libs.UploadFile(fileHeader, "profile_images")
			if err != nil {
				return ProfileResponse{}, err
			}
			imagePath = &url
		} else {
			uploadDir := "./uploads/profile"
			if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
				return ProfileResponse{}, err
			}

			filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fileHeader.Filename))
			filePath := filepath.Join(uploadDir, filename)
			imagePath = &filePath

			src, err := fileHeader.Open()
			if err != nil {
				return ProfileResponse{}, err
			}
			defer src.Close()

			dst, err := os.Create(filePath)
			if err != nil {
				return ProfileResponse{}, err
			}
			defer dst.Close()

			if _, err := io.Copy(dst, src); err != nil {
				return ProfileResponse{}, err
			}
		}
	}

	var profileID int64
	err := db.QueryRow(ctx, `SELECT id FROM profile WHERE user_id=$1`, userID).Scan(&profileID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return ProfileResponse{}, err
	}

	if profileID > 0 {
		_, err = db.Exec(ctx, `
			UPDATE profile
			SET phone = COALESCE(NULLIF($1, ''), phone),
				address = COALESCE(NULLIF($2, ''), address),
				image = COALESCE($3, image),
				updated_at = NOW()
			WHERE user_id = $4
		`, phone, address, imagePath, userID)
	} else {
		_, err = db.Exec(ctx, `
			INSERT INTO profile (user_id, phone, address, image, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW())
		`, userID, phone, address, imagePath)
	}
	if err != nil {
		return ProfileResponse{}, err
	}

	if fullname != "" || email != "" {
		_, err := db.Exec(ctx, `
			UPDATE users
			SET fullname = COALESCE(NULLIF($1, ''), fullname),
				email = COALESCE(NULLIF($2, ''), email),
				updated_at = NOW()
			WHERE id = $3
		`, fullname, email, userID)
		if err != nil {
			return ProfileResponse{}, err
		}
	}

	var resp ProfileResponse
	var dbFullname, dbEmail *string
	err = db.QueryRow(ctx, `
		SELECT 
			COALESCE(p.id,0), u.fullname, u.email, p.image, p.phone, p.address, u.id, 
			COALESCE(p.created_at, NOW()), COALESCE(p.updated_at, NOW())
		FROM users u
		LEFT JOIN profile p ON p.user_id = u.id
		WHERE u.id=$1
	`, userID).Scan(
		&resp.ID,
		&dbFullname,
		&dbEmail,
		&resp.Image,
		&resp.Phone,
		&resp.Address,
		&resp.UserID,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)
	if err != nil {
		return ProfileResponse{}, err
	}

	if dbFullname != nil {
		resp.Fullname = *dbFullname
	}
	if dbEmail != nil {
		resp.Email = *dbEmail
	}

	return resp, nil
}




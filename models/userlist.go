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

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserList struct {
	ID        int64     `json:"id"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Profile   *Profile `json:"profile,omitempty"` 
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Profile struct {
	Image   *string `json:"image,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Address *string `json:"address,omitempty"`
}

type AdminUserRequest struct {
    Fullname string                  `form:"fullname" binding:"required"`
    Email    string                  `form:"email" binding:"required,email"`
    Password string                  `form:"password" binding:"required,min=6"`
    Phone    string                  `form:"phone"`
    Address  string                  `form:"address"`
    Image    *multipart.FileHeader   `form:"image"`
    Role     string                  `form:"role" binding:"required"`
}


func AddUser(db *pgxpool.Pool, fullname, email, password, role string, phone, address string, fileHeader *multipart.FileHeader) (*UserList, *Profile, error) {
	ctx := context.Background()

	var imagePath *string
	if fileHeader != nil {
		if fileHeader.Size > 2*1024*1024 {
			return nil, nil, errors.New("file size exceeds 2MB")
		}

		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return nil, nil, errors.New("file type must be jpg, jpeg, or png")
		}

		if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
			return nil, nil, err
		}

		filename := fmt.Sprintf("uploads/%d_%d%s", time.Now().UnixNano(), time.Now().UnixNano(), ext) // unique filename
		src, err := fileHeader.Open()
		if err != nil {
			return nil, nil, err
		}
		defer src.Close()

		dst, err := os.Create(filename)
		if err != nil {
			return nil, nil, err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return nil, nil, err
		}

		imagePath = &filename
	}

	hashed, err := libs.HashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	var userID int64
	var u UserList
	err = db.QueryRow(ctx, `
		INSERT INTO users (fullname, email, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, fullname, email, role, created_at, updated_at
	`, fullname, email, hashed, role).Scan(&userID, &u.Fullname, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, nil, errors.New("email already exists")
		}
		return nil, nil, err
	}
	u.ID = userID

	_, err = db.Exec(ctx, `
		INSERT INTO profile (user_id, image, phone, address)
		VALUES ($1, $2, $3, $4)
	`, userID, imagePath, &phone, &address)
	if err != nil {
		return nil, nil, err
	}

	p := &Profile{
		Image:   imagePath,
		Phone:   &phone,
		Address: &address,
	}

	u.Profile = p

	return &u, p, nil
}

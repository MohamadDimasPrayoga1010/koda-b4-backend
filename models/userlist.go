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
	Profile   *Profile  `json:"profile,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Profile struct {
	Image   *string `json:"image,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Address *string `json:"address,omitempty"`
}

type AdminUserRequest struct {
    Fullname string                `form:"fullname" binding:"required" validate:"required,min=3"`
    Email    string                `form:"email" binding:"required,email" validate:"required,email"`
    Password string                `form:"password" binding:"required,min=6" validate:"required,min=6"`
    Phone    string                `form:"phone" validate:"omitempty,min=10"`
    Address  string                `form:"address" validate:"omitempty"`
    Role     string                `form:"role" binding:"required" validate:"required,oneof=admin staff user"`
    Image    *multipart.FileHeader `form:"image" validate:"omitempty"`
}

func AddUser(db *pgxpool.Pool, fullname, email, password, role string, phone, address string, imagePath string,) (*UserList, *Profile, error) {

	ctx := context.Background()

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
	`, fullname, email, hashed, role).Scan(
		&userID, &u.Fullname, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, nil, errors.New("email already exists")
		}
		return nil, nil, err
	}

	u.ID = userID

	var imgPtr *string
	if imagePath != "" {
		imgPtr = &imagePath
	}

	_, err = db.Exec(ctx, `
		INSERT INTO profile (user_id, image, phone, address)
		VALUES ($1, $2, $3, $4)
	`, userID, imgPtr, phone, address)

	if err != nil {
		return nil, nil, err
	}

	p := &Profile{
		Image:   imgPtr,
		Phone:   &phone,
		Address: &address,
	}

	u.Profile = p

	return &u, p, nil
}



func UpdateUser(db *pgxpool.Pool, userID int64, fullname, email, password, phone, address string, fileHeader *multipart.FileHeader) (*UserList, *Profile, error) {
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

		if os.Getenv("CLOUDINARY_API_KEY") != "" {
			url, err := libs.UploadFile(fileHeader, "profile_images")
			if err != nil {
				return nil, nil, err
			}
			imagePath = &url
		} else {
			filename := fmt.Sprintf("uploads/%d_%d%s", userID, time.Now().UnixNano(), ext)
			if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
				return nil, nil, err
			}
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
	}

	args := []interface{}{}
	fields := []string{}
	argIdx := 1

	if fullname != "" {
		fields = append(fields, fmt.Sprintf("fullname=$%d", argIdx))
		args = append(args, fullname)
		argIdx++
	}
	if email != "" {
		fields = append(fields, fmt.Sprintf("email=$%d", argIdx))
		args = append(args, email)
		argIdx++
	}
	if password != "" {
		hashed, err := libs.HashPassword(password)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, fmt.Sprintf("password=$%d", argIdx))
		args = append(args, hashed)
		argIdx++
	}

	if len(fields) > 0 {
		query := fmt.Sprintf("UPDATE users SET %s, updated_at=NOW() WHERE id=$%d", strings.Join(fields, ","), argIdx)
		args = append(args, userID)
		if _, err := db.Exec(ctx, query, args...); err != nil {
			return nil, nil, err
		}
	}

	profileFields := []string{}
	profileArgs := []interface{}{}
	pIdx := 1

	if imagePath != nil {
		profileFields = append(profileFields, fmt.Sprintf("image=$%d", pIdx))
		profileArgs = append(profileArgs, *imagePath)
		pIdx++
	}
	if phone != "" {
		profileFields = append(profileFields, fmt.Sprintf("phone=$%d", pIdx))
		profileArgs = append(profileArgs, phone)
		pIdx++
	}
	if address != "" {
		profileFields = append(profileFields, fmt.Sprintf("address=$%d", pIdx))
		profileArgs = append(profileArgs, address)
		pIdx++
	}

	if len(profileFields) > 0 {
		query := fmt.Sprintf("UPDATE profile SET %s, updated_at=NOW() WHERE user_id=$%d", strings.Join(profileFields, ","), pIdx)
		profileArgs = append(profileArgs, userID)
		if _, err := db.Exec(ctx, query, profileArgs...); err != nil {
			return nil, nil, err
		}
	}

	var u UserList
	var p Profile
	err := db.QueryRow(ctx, `
	SELECT u.id, u.fullname, u.email, u.role, p.image, p.phone, p.address, u.created_at, u.updated_at
	FROM users u
	LEFT JOIN profile p ON p.user_id=u.id
	WHERE u.id=$1
	`, userID).Scan(
		&u.ID, &u.Fullname, &u.Email, &u.Role, &p.Image, &p.Phone, &p.Address, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, nil, err
	}
	u.Profile = &p
	return &u, &p, nil
}
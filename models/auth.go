package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)


type UserRegister struct {
    Fullname string `json:"fullname" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Role     string `json:"role"`
}


type UserResponse struct {
    ID        int64     `json:"id"`
    Fullname  string    `json:"fullname"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    Token     string    `json:"token,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type UserLogin struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}


type ForgotPasswordRequest struct {
    Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
    Email string `json:"email" binding:"required,email"`
    OTP   string `json:"otp" binding:"required,len=6"`
}

type ResetPasswordRequest struct {
    Token    string `json:"token" binding:"required"`
    Password string `json:"password" binding:"required,min=6"`
}

type ForgotPassword struct {
	ID        int       `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}


func RegisterUser(db *pgxpool.Pool, user UserRegister, hashedPassword string) (UserResponse, error) {
    var resp UserResponse

    role := user.Role
    if role == "" {
        role = "user"
    }

    query := `
        INSERT INTO users (fullname, email, password, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id, fullname, email, role, created_at, updated_at
    `

    err := db.QueryRow(
        context.Background(),
        query,
        user.Fullname,
        user.Email,
        hashedPassword,
        role,
    ).Scan(
        &resp.ID,
        &resp.Fullname,
        &resp.Email,
        &resp.Role,
        &resp.CreatedAt,
        &resp.UpdatedAt,
    )

    if err != nil {
        return UserResponse{}, err
    }

    return resp, nil
}


func LoginUser(db *pgxpool.Pool, email string) (*UserResponse, string, string, error) {
	var user UserResponse
	var hashedPassword string

	query := `
		SELECT id, fullname, email, password, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := db.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&hashedPassword,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, "", "", err
	}

	return &user, hashedPassword, user.Role, nil
}

func CreateForgotPassword(db *pgxpool.Pool, userID int64, token string, duration time.Duration) error {
	expires := time.Now().Add(duration)
	_, err := db.Exec(context.Background(),
		`INSERT INTO forgot_password (user_id, token, expires_at) VALUES ($1, $2, $3)`,
		userID, token, expires,
	)
	return err
}

func GetForgotPasswordByToken(db *pgxpool.Pool, token string) (*ForgotPassword, error) {
	var fp ForgotPassword
	err := db.QueryRow(context.Background(),
		`SELECT id, user_id, token, expires_at, created_at FROM forgot_password WHERE token=$1`,
		token,
	).Scan(&fp.ID, &fp.UserID, &fp.Token, &fp.ExpiresAt, &fp.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &fp, nil
}

func DeleteForgotPassword(db *pgxpool.Pool, id int) error {
	_, err := db.Exec(context.Background(), `DELETE FROM forgot_password WHERE id=$1`, id)
	return err
}
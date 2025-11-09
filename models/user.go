package models

import "time"

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
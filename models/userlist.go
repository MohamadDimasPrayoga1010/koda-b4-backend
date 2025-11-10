package models

import "time"

type UserList struct {
	ID        int64     `json:"id"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Profile   *Profile  `json:"profile,omitempty"` 
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Profile struct {
	Image   string `json:"image"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type AdminUserRequest struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password,omitempty"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Image    string `json:"image"`
	Role     string `json:"role" binding:"required"` 
}
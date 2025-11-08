package controllers

import (
	"context"
	"strings"

	"main/libs"
	"main/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthController struct {
	DB *pgxpool.Pool
}

func (ac *AuthController) Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	hashedPassword, err := libs.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to hash password",
			Data:    nil,
		})
		return
	}

	query := `
		INSERT INTO users (fullname, email, password, role)
		VALUES ($1, $2, $3, 'user')
		RETURNING id, fullname, email, role
	`
	var userResp models.UserResponse
	err = ac.DB.QueryRow(context.Background(), query, user.Fullname, user.Email, hashedPassword).
		Scan(&userResp.ID, &userResp.Fullname, &userResp.Email, &userResp.Role)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			ctx.JSON(409, models.Response{
				Success: false,
				Message: "Email already registered",
				Data:    nil,
			})
			return
		}
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to register user",
			Data:    nil,
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "User registered successfully",
		Data:    userResp,
	})
}


func (ac *AuthController) Login(ctx *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	var user models.User
	query := `
		SELECT id, fullname, email, password, role
		FROM users
		WHERE email = $1
	`
	err := ac.DB.QueryRow(context.Background(), query, input.Email).Scan(
		&user.ID, &user.Fullname, &user.Email, &user.Password, &user.Role,
	)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Email or password incorrect",
			Data:    nil,
		})
		return
	}

	ok, err := libs.VerifyPassword(input.Password, user.Password)
	if err != nil || !ok {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Email or password incorrect",
			Data:    nil,
		})
		return
	}

	
	token, err := libs.GenerateToken(int(user.ID), user.Email, user.Role)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to generate token",
			Data:    nil,
		})
		return
	}

	userResp := models.UserResponse{
		ID:       user.ID,
		Fullname: user.Fullname,
		Email:    user.Email,
		Role:     user.Role,
		Token:    token,
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Login success",
		Data:    userResp,
	})
}

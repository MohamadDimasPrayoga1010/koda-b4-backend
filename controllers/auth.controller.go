package controllers

import (
	"context"
	"fmt"
	"strings"

	"main/libs"
	"main/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthController struct {
	DB *pgxpool.Pool
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User Register"
// @Success      201  {object}  models.Response{data=models.UserResponse}
// @Failure      400  {object}  models.Response
// @Failure      409  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /auth/register [post]
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
	fmt.Println("inputan password: ",user.Password)

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
		VALUES ($1, $2, $3, $4)
		RETURNING id, fullname, email, role
	`
	var userResp models.UserResponse
	err = ac.DB.QueryRow(context.Background(), query, user.Fullname, user.Email, hashedPassword, user.Role).
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

// Login godoc
// @Summary      User login
// @Description  Login with email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        login  body      map[string]string  true  "Login Payload"
// @Success      200  {object}  models.Response{data=models.UserResponse}
// @Failure      400  {object}  models.Response
// @Failure      401  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /auth/login [post]
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
	if err != nil {
		fmt.Println(err)
		return
	}
	if !ok {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Email or password salah",
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

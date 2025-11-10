package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

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
// @Param        user  body      models.UserRegister  true  "User Register Payload"
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
	fmt.Println("inputan password: ", user.Password)

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
// @Param        login  body      models.UserLogin  true  "Login Payload"
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

// ForgotPassword godoc
// @Summary Request OTP for forgot password
// @Description Send OTP to user's email
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.ForgotPasswordRequest true "Email for forgot password"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /auth/forgot-password [post]
func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var req models.ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}
	var exists bool
	err := ac.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", req.Email).Scan(&exists)
	if err != nil || !exists {
		ctx.JSON(404, models.Response{Success: false, Message: "Email not found"})
		return
	}

	otp := libs.GenerateOTP(6)
	expire := time.Now().Add(5 * time.Minute)

	_, err = ac.DB.Exec(context.Background(),
		"UPDATE users SET reset_otp=$1, reset_expires=$2 WHERE email=$3",
		otp, expire, req.Email,
	)

	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to generate token",
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success generate token send to email",
	})

}

// VerifyOTP godoc
// @Summary Verify OTP
// @Description Verify OTP for password reset
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.VerifyOTPRequest true "Verify OTP payload"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /auth/verify-otp [post]
func (ac *AuthController) VerifyOTP(ctx *gin.Context) {
	var req models.VerifyOTPRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	var dbOTP string
	var expire time.Time
	err := ac.DB.QueryRow(context.Background(),
		"SELECT reset_otp, reset_expires FROM users WHERE email=$1", req.Email,
	).Scan(&dbOTP, &expire)
	if err != nil || dbOTP == "" {
		ctx.JSON(401, models.Response{
			Success: false, 
			Message: "Invalid OTP"})
		return
	}

	if time.Now().After(expire) {
		ctx.JSON(401, models.Response{
			Success: false, 
			Message: "OTP expired"})
		return
	}

	if dbOTP != req.OTP {
		ctx.JSON(401, models.Response{
			Success: false, 
			Message: "OTP incorrect"})
		return
	}

	token, _ := libs.GenerateTokenForReset(req.Email)

	ac.DB.Exec(context.Background(),
		"UPDATE users SET reset_token=$1 WHERE email=$2", token, req.Email,
	)

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "OTP verified successfully",
		Data:    map[string]string{"token": token},
	})

}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using verified OTP token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.ResetPasswordRequest true "Reset password payload"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /auth/reset-password [patch]
func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	var req models.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false, 
			Message: "Invalid request body"})
		return
	}

	var email string
	err := ac.DB.QueryRow(context.Background(),
		"SELECT email FROM users WHERE reset_token=$1", req.Token,
	).Scan(&email)
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false, 
			Message: "Invalid or expired token"})
		return
	}

	hashed, err := libs.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false, 
			Message: "Failed to hash password"})
		return
	}

	_, err = ac.DB.Exec(context.Background(),
		"UPDATE users SET password=$1, reset_token=NULL, reset_otp=NULL, reset_expires=NULL WHERE email=$2",
		hashed, email,
	)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false, 
			Message: "Failed to reset password"})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true, 
		Message: "Password has been reset successfully"})
}

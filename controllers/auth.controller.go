package controllers

import (
	"coffeeder-backend/libs"
	"coffeeder-backend/models"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthController struct {
	DB *pgxpool.Pool
}

// Register godoc
// @Summary      Register a new user or admin
// @Description  Create a new user account. If role is "admin", response message will indicate admin registration success.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param login body models.UserLogin true "Login Payload"
// @Success      201  {object}  models.Response{data=models.UserResponse} "User or Admin registered successfully"
// @Failure      400  {object}  models.Response "Invalid request body or password too short"
// @Failure      409  {object}  models.Response "Email already registered"
// @Failure      500  {object}  models.Response "Internal server error"
// @Router       /auth/register [post]
func (ac *AuthController) Register(ctx *gin.Context) {
	var req models.UserRegister

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	hashed, err := libs.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to hash password",
		})
		return
	}

	user, err := models.RegisterUser(ac.DB, req, hashed)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			ctx.JSON(409, models.Response{
				Success: false,
				Message: "Email already registered",
			})
			return
		}
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to register user",
		})
		return
	}

	if req.Role == "admin" {
		ctx.JSON(201, models.Response{
			Success: true,
			Message: "Admin registered successfully",
			Data:    user,
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

// Login godoc
// @Summary      User login
// @Description  Login using email and password
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
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	user, hashedPassword, _, err := models.LoginUser(ac.DB, input.Email)
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Email or password incorrect",
		})
		return
	}

	ok, err := libs.VerifyPassword(input.Password, hashedPassword)
	if err != nil || !ok {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Email or password incorrect",
		})
		return
	}

	token, err := libs.GenerateToken(int(user.ID), user.Email, user.Role)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to generate token",
		})
		return
	}

	user.Token = token

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Login success",
		Data:    user,
	})
}

// ForgotPassword godoc
// @Summary      Request password
// @Description  Generate OTP for user to forgot password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.ForgotPasswordRequest  true  "Email for forgot password"
// @Success      200   {object}  models.Response
// @Failure      400   {object}  models.Response
// @Failure      404   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /auth/forgot-password [post]

func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("Invalid request body:", err)
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	fmt.Println("ForgotPassword request for email:", req.Email)

	var userID int64
	err := ac.DB.QueryRow(context.Background(),
		"SELECT id FROM users WHERE email=$1", req.Email,
	).Scan(&userID)
	if err != nil {
		fmt.Println("User not found:", err)
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Email not found",
		})
		return
	}

	otp := libs.GenerateOTP(6)
	fmt.Println("Generated OTP:", otp)

	err = models.CreateForgotPassword(ac.DB, userID, otp, 2*time.Minute)
	if err != nil {
		fmt.Println("Insert failed, trying to update existing record:", err)
		_, err2 := ac.DB.Exec(context.Background(),
			`UPDATE forgot_password 
			 SET token=$1, expires_at=$2 
			 WHERE user_id=$3`,
			otp, time.Now().Add(2*time.Minute), userID,
		)
		if err2 != nil {
			fmt.Println("Failed to update existing ForgotPassword record:", err2)
			ctx.JSON(500, models.Response{
				Success: false,
				Message: "Failed to generate OTP",
			})
			return
		}
	}

	fmt.Println("ForgotPassword record created/updated successfully")

	err = libs.SendOTPEmail(libs.SendOptions{
		To:         []string{req.Email},
		Subject:    "OTP Reset Password",
		Body:       fmt.Sprintf("Your OTP is: %s. It will expire in 2 minutes.", otp),
		BodyIsHTML: false,
	})
	if err != nil {
		fmt.Println("Error sending OTP email:", err)
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to send OTP email, please check SMTP configuration",
		})
		return
	}

	fmt.Println("OTP sent successfully to email:", req.Email)
	ctx.JSON(200, models.Response{
		Success: true,
		Message: "OTP has been sent to your email",
	})
}


// VerifyOTP godoc
// @Summary      Verify OTP
// @Description  Verify OTP sent to user's email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.VerifyOTPRequest  true  "OTP verification payload"
// @Success      200   {object}  models.Response
// @Failure      400   {object}  models.Response
// @Failure      401   {object}  models.Response
// @Router       /auth/verify-otp [post]
func (ac *AuthController) VerifyOTP(ctx *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		OTP   string `json:"otp" binding:"required,len=6"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	var userID int64
	err := ac.DB.QueryRow(context.Background(), "SELECT id FROM users WHERE email=$1", req.Email).Scan(&userID)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Email not found"})
		return
	}

	fp, err := models.GetForgotPasswordByToken(ac.DB, req.OTP)
	if err != nil || fp.UserID != userID {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Invalid OTP"})
		return
	}

	if time.Now().After(fp.ExpiresAt) {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "OTP expired"})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "OTP verified successfully",
		Data:    map[string]string{"token": req.OTP},
	})
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Reset user password using OTP token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.ResetPasswordRequest  true  "Reset password payload"
// @Success      200   {object}  models.Response
// @Failure      400   {object}  models.Response
// @Failure      401   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /auth/reset-password [patch]
func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	fp, err := models.GetForgotPasswordByToken(ac.DB, req.Token)
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
		"UPDATE users SET password=$1 WHERE id=$2",
		hashed, fp.UserID,
	)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to reset password"})
		return
	}

	models.DeleteForgotPassword(ac.DB, fp.ID)

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Password has been reset successfully",
	})
}

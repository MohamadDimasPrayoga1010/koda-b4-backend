package routers

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthRoutes(r *gin.Engine, pg *pgxpool.Pool) {
	authController := controllers.AuthController{DB: pg}

	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/forgot-password", authController.ForgotPassword)
		auth.POST("/verify-otp", authController.VerifyOTP)
		auth.PATCH("/reset-password", authController.ResetPassword)
	}
}

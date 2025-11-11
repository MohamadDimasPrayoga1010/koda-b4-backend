package routers

import (
	"main/controllers"
	"main/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AdminUserRoutes(r *gin.Engine, pg *pgxpool.Pool) {
	uc := controllers.UserController{DB: pg}

	admin := r.Group("/admin")
	admin.Use(middlewares.AuthMiddleware("admin"))
	{
		admin.GET("/users", uc.GetUsersList)
		admin.GET("/users/:id", uc.GetUserByID)
		admin.POST("/users", uc.AddUser)
		admin.PATCH("/users/:id", uc.EditUser)
		admin.DELETE("/users/:id", uc.DeleteUser)
	}
}

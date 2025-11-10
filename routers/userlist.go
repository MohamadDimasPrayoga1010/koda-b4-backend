package routers

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AdminUserRoutes(r *gin.Engine, pg *pgxpool.Pool) {
	uc := controllers.UserController{DB: pg}

	admin := r.Group("/admin")
	{
		admin.GET("/users", uc.GetUsersList)
		admin.GET("/users/:id", uc.GetUserByID)
		admin.POST("/users", uc.AddUser)
		admin.PATCH("/users/:id", uc.EditUser)
		admin.DELETE("/users/:id", uc.DeleteUser)
	}
}

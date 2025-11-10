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
		admin.GET("/userslist", uc.GetUsersList)      
		admin.POST("/userslist", uc.AddUser)          
		admin.PATCH("/userslist/:id", uc.EditUser)  
		admin.DELETE("/userslist/:id", uc.DeleteUser)    
	}
}

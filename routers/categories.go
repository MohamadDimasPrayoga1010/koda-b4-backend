package routers

import (
	"main/controllers"
	"main/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CategoryRoutes(r *gin.Engine, pg *pgxpool.Pool) {
	cc := controllers.CategoryController{DB: pg}

	admin := r.Group("/admin")
	admin.Use(middlewares.AuthMiddleware("admin"))
	{
		admin.GET("/categories", cc.GetCategories)
		admin.GET("/categories/:id", cc.GetCategoryByID)
		admin.POST("/categories", cc.CreateCategory)
		admin.PATCH("/categories/:id", cc.UpdateCategory)
		admin.DELETE("/categories/:id", cc.DeleteCategory)
	}
}

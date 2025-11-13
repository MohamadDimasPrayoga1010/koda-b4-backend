package routers

import (
	"main/controllers"
	"main/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TransactionRoutes(r *gin.Engine, pg *pgxpool.Pool) {
	tc := controllers.TransactionController{DB: pg}

	admin := r.Group("/admin")
	admin.Use(middlewares.AuthMiddleware("admin"))
	{
		admin.GET("/transactions", tc.GetTransactions)
		admin.GET("/transactions/:id", tc.GetTransactionByID)
		admin.PATCH("/transactions/:id/status", tc.UpdateTransactionStatus)
		admin.DELETE("/transactions/:id", tc.DeleteTransaction)
	}
	r.GET("/history", middlewares.AuthMiddleware(""), tc.GetHistoryTransactions)
	r.GET("/history/:id", middlewares.AuthMiddleware(""), tc.GetHistoryDetailById)
}

package routers

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func OrderRoutes(r *gin.Engine, pg *pgxpool.Pool) {
	oc := controllers.OrderController{DB: pg}

	admin := r.Group("/admin")
	{
		admin.GET("/orders", oc.GetOrders)
		admin.PATCH("/orders/:id/status", oc.UpdateOrderStatus)
		admin.DELETE("/orders/:id", oc.DeleteOrder)
	}
}

package routers

import (
	"coffeeder-backend/libs"
	"coffeeder-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter(pg *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, models.Response{
			Success: true,
			Message: "Backend is running well",
		})
	})

	r.Use(libs.SetupCORS())
	AuthRoutes(r, pg)
	ProductRoutes(r, pg)
	TransactionRoutes(r, pg)
	AdminUserRoutes(r, pg)
	CategoryRoutes(r, pg)
	return r
}

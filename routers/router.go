package routers

import (
	"main/libs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter (pg *pgxpool.Pool) *gin.Engine{
	r := gin.Default()

	r.Use(libs.SetupCORS())
	AuthRoutes(r, pg)
	ProductRoutes(r, pg)
	TransactionRoutes(r, pg)
	AdminUserRoutes(r, pg)
	CategoryRoutes(r, pg)
	return r
}
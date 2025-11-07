package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter (pg *pgxpool.Pool) *gin.Engine{
	r := gin.Default()

	return r
}
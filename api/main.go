package handler

import (
	"main/configs"
	"main/docs"
	"main/libs"
	"main/routers"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var router *gin.Engine

func init() {
	pg := configs.InitDbConfig()
	libs.InitRedis()
	router = routers.InitRouter(pg)

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}

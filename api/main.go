package handler

import (
	"coffeeder-backend/configs"
	_ "coffeeder-backend/docs"
	"coffeeder-backend/libs"
	"coffeeder-backend/routers"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var router *gin.Engine

func initRouter() *gin.Engine {
	if router != nil {
		return router
	}

	pg := configs.InitDbConfig()
	libs.InitRedis()

	router = routers.InitRouter(pg)
	router.Use(gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func Handler(w http.ResponseWriter, r *http.Request) {
	routerEngine := initRouter()
	routerEngine.ServeHTTP(w, r)
}

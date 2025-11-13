package handler

import (
	"main/configs"
	"main/docs"
	"main/libs"
	"main/routers"
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
	router.Use(gin.Recovery())


	router = routers.InitRouter(pg)


	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func Handler(w http.ResponseWriter, r *http.Request) {
	routerEngine := initRouter()
	routerEngine.ServeHTTP(w, r)
}

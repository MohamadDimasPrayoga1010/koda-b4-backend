package main

import (
	"main/configs"
	"main/docs"
	"main/libs"
	"main/routers"
	"net/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	pg := configs.InitDbConfig()
	rtr := routers.InitRouter(pg)
	libs.InitRedis()

	docs.SwaggerInfo.BasePath = "/"
	rtr.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	rtr.ServeHTTP(w, r)
}

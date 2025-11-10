package main

import (

	"main/configs"
	"main/docs"
	"main/routers"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Coffee Shop API
// @version 1.0
// @description API for Coffee Shop Application
// @host localhost:8085
// @BasePath /
func main() {
	pg := configs.InitDbConfig()
	r := routers.InitRouter(pg)

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// fmt.Println(libs.HashPassword("admin123"))
	r.Run(":8085")
}

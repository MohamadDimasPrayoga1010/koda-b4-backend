package main

import (
	"coffeeder-backend/configs"
	_ "coffeeder-backend/docs" 
	"coffeeder-backend/libs"
	"coffeeder-backend/routers"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Coffee Shop API
// @version 1.0
// @description API for Coffee Shop Application
// @host coffeeder-backend.vercel.app
// @schemes https
// @BasePath /
func main() {
	godotenv.Load()
	pg := configs.InitDbConfig()
	r := routers.InitRouter(pg)
	libs.InitRedis()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8085")
}

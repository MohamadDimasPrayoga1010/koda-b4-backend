package main

import (
	"main/configs"
	"main/routers"
)

func main() {
	pg := configs.InitDbConfig()
	r := routers.InitRouter(pg)
	
	r.Run(":8085")
}

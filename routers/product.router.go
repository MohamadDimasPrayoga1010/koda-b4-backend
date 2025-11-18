package routers

import (
	"coffeeder-backend/controllers"
	"coffeeder-backend/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProductRoutes(r *gin.Engine, pg *pgxpool.Pool) {
	pc := controllers.ProductController{DB: pg}

	admin := r.Group("/admin")
	admin.Use(middlewares.AuthMiddleware("admin"))
	{
		admin.POST("/products", pc.CreateProduct)
		admin.GET("/products", pc.GetProducts)
		admin.GET("/products/:id", pc.GetProductByID)
		admin.PATCH("/products/:id", pc.UpdateProduct)
		admin.DELETE("/products/:id", pc.DeleteProduct)
		admin.GET("/products/:id/images", pc.GetProductImages)             
		admin.GET("/products/:id/images/:image_id", pc.GetProductImageByID) 
		admin.PATCH("/products/:id/images/:image_id", pc.UpdateProductImage) 
		admin.DELETE("/products/:id/images/:image_id", pc.DeleteProductImage) 
	}
	r.GET("/favorite-products", middlewares.AuthMiddleware(""),pc.GetFavoriteProducts)
	r.GET("/products",pc.FilterProducts)
	r.GET("/products/:id",middlewares.AuthMiddleware("") ,pc.GetProductDetail)
	r.POST("/cart",middlewares.AuthMiddleware(""), pc.AddToCart)
	r.GET("/cart", middlewares.AuthMiddleware(""), pc.GetCart)
	r.DELETE("/deletecart", middlewares.AuthMiddleware(""), pc.DeleteCart)
	r.POST("/transactions", middlewares.AuthMiddleware(""), pc.CreateTransaction)
}

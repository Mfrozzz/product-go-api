package main

import (
	"product-go-api/controller"
	"product-go-api/db"
	"product-go-api/middleware"
	"product-go-api/repository"
	"product-go-api/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.Use(middleware.RateLimiter())

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	ProductRepository := repository.NewProductRepository(dbConnection)

	ProductUseCase := usecase.NewProductUsecase(ProductRepository)

	ProductController := controller.NewProductController(ProductUseCase)

	server.GET("/products", ProductController.GetProducts)
	server.POST("/product", ProductController.CreateProduct)
	server.GET("/product/:id_product", ProductController.GetProductById)
	server.DELETE("/product/:id_product", ProductController.DeleteProduct)
	server.PUT("/product/:id_product", ProductController.UpdateProduct)

	server.Run(":8000")
}

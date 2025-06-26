package main

import (
	"os"
	"product-go-api/controller"
	"product-go-api/db"
	"product-go-api/middleware"
	"product-go-api/repository"
	"product-go-api/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	server := gin.Default()
	server.Use(cors.Default())
	server.Use(middleware.RateLimiter())

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	UserRepository := repository.NewUserRepository(dbConnection)
	UserUseCase := usecase.NewUserUsecase(UserRepository)
	UserController := controller.NewUserController(UserUseCase)

	ProductRepository := repository.NewProductRepository(dbConnection)
	ProductUseCase := usecase.NewProductUsecase(ProductRepository)
	ProductController := controller.NewProductController(ProductUseCase)

	server.POST("/register", UserController.CreateUser)
	server.POST("/login", UserController.GetUserByEmail)

	protectedRoutes := server.Group("/api")
	protectedRoutes.Use(middleware.AuthMiddleware())

	protectedRoutes.GET("/user/info", UserController.GetUserInfo)
	protectedRoutes.GET("/users/:id_user", UserController.GetUserById)
	protectedRoutes.PUT("/users/:id_user", UserController.UpdateUser)

	protectedRoutes.GET("/products", ProductController.GetProducts)
	protectedRoutes.POST("/products", ProductController.CreateProduct)
	protectedRoutes.GET("/products/:id_product", ProductController.GetProductById)
	protectedRoutes.PUT("/products/:id_product", ProductController.UpdateProduct)

	adminRoutes := protectedRoutes.Group("/admin")
	adminRoutes.Use(middleware.RequireAdmin())
	adminRoutes.DELETE("/products/:id_product", ProductController.DeleteProduct)
	adminRoutes.DELETE("/users/:id_user", UserController.DeleteUser)

	server.Run(os.Getenv("PORT"))
}

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tricong1998/go-ecom/cmd/product/internal/api/handlers"
	"github.com/tricong1998/go-ecom/cmd/product/internal/repository"
	"github.com/tricong1998/go-ecom/cmd/product/internal/services"
	"gorm.io/gorm"
)

func SetupRoutes(routes *gin.Engine, db *gorm.DB) {
	userRepo := repository.NewProductRepository(db)
	userService := services.NewProductService(userRepo)
	userHandler := handlers.NewProductHandler(userService)

	userGroup := routes.Group("products")
	{
		userGroup.POST("", userHandler.CreateProduct)
		userGroup.GET("/:id", userHandler.ReadProduct)
		userGroup.GET("", userHandler.ListProducts)
		userGroup.PUT("/:id", userHandler.UpdateProduct)
		userGroup.DELETE("/:id", userHandler.DeleteProduct)
	}
}

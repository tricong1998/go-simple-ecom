package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/tricong1998/go-ecom/cmd/product/internal/api/handlers"
	"github.com/tricong1998/go-ecom/cmd/product/internal/config"
	"github.com/tricong1998/go-ecom/cmd/product/internal/repository"
	"github.com/tricong1998/go-ecom/cmd/product/internal/services"
	"github.com/tricong1998/go-ecom/pkg/gin/middleware"
	"github.com/tricong1998/go-ecom/pkg/token"
	"gorm.io/gorm"
)

func SetupRoutes(routes *gin.Engine, db *gorm.DB, config *config.Config, log *zerolog.Logger) {
	tokenMaker, err := token.NewJWTMaker(config.Auth.AccessTokenSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create token maker")
		return
	}
	userRepo := repository.NewProductRepository(db)
	userService := services.NewProductService(userRepo)
	userHandler := handlers.NewProductHandler(userService)

	userGroup := routes.Group("products")
	authRoutes := userGroup.Group("/").Use(middleware.AuthMiddleware(tokenMaker, []string{}))
	{
		authRoutes.GET("/:id", userHandler.ReadProduct)
		authRoutes.GET("", userHandler.ListProducts)
	}
	adminRoutes := userGroup.Group("/").Use(middleware.AuthMiddleware(tokenMaker, []string{"admin"}))
	{
		adminRoutes.POST("", userHandler.CreateProduct)
		adminRoutes.PUT("/:id", userHandler.UpdateProduct)
		adminRoutes.DELETE("/:id", userHandler.DeleteProduct)
	}
}

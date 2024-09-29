package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/tricong1998/go-ecom/cmd/user/internal/api/handlers"
	"github.com/tricong1998/go-ecom/cmd/user/internal/config"
	"github.com/tricong1998/go-ecom/cmd/user/internal/repository"
	"github.com/tricong1998/go-ecom/cmd/user/internal/services"
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
	jwtService := services.NewJwtService(tokenMaker, config.Auth)
	userRepo := repository.NewUserRepository(db)
	userPointRepo := repository.NewUserPointRepository(db)
	userPointService := services.NewUserPointService(userPointRepo)
	userService := services.NewUserService(userRepo, userPointService)
	userHandler := handlers.NewUserHandler(userService, jwtService)

	userGroup := routes.Group("users")
	{
		userGroup.POST("", userHandler.CreateUser)
		userGroup.POST("/login", userHandler.Login)
	}
	authRoutes := userGroup.Group("/").Use(middleware.AuthMiddleware(tokenMaker, []string{}))
	{
		authRoutes.GET("/me", userHandler.ReadMe)
		authRoutes.GET("/:id", userHandler.ReadUser)
		authRoutes.PUT("/update-me", userHandler.UpdateMe)
		authRoutes.DELETE("/:id", userHandler.DeleteUser)
	}

	adminRoutes := userGroup.Group("/").Use(middleware.AuthMiddleware(tokenMaker, []string{"admin"}))
	{
		adminRoutes.GET("", userHandler.ListUsers)
	}
}

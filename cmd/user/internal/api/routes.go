package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tricong1998/go-ecom/cmd/user/internal/api/handlers"
	"github.com/tricong1998/go-ecom/cmd/user/internal/repository"
	"github.com/tricong1998/go-ecom/cmd/user/internal/services"
	"gorm.io/gorm"
)

func SetupRoutes(routes *gin.Engine, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userPointRepo := repository.NewUserPointRepository(db)
	userPointService := services.NewUserPointService(userPointRepo)
	userService := services.NewUserService(userRepo, userPointService)
	userHandler := handlers.NewUserHandler(userService)

	userGroup := routes.Group("users")
	{
		userGroup.POST("", userHandler.CreateUser)
		userGroup.GET("/:id", userHandler.ReadUser)
		userGroup.GET("", userHandler.ListUsers)
		userGroup.PUT("/:id", userHandler.UpdateUser)
		userGroup.DELETE("/:id", userHandler.DeleteUser)
	}
}

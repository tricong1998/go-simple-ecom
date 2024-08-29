package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tricong1998/go-ecom/internal/api/handlers"
	"github.com/tricong1998/go-ecom/internal/repository"
	"github.com/tricong1998/go-ecom/internal/services"
	"gorm.io/gorm"
)

func SetupRoutes(routes *gin.Engine, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
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

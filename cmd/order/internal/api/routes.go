package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tricong1998/go-ecom/cmd/order/internal/api/handlers"
	"github.com/tricong1998/go-ecom/cmd/order/internal/config"
	"github.com/tricong1998/go-ecom/cmd/order/internal/gateway/user/grpc"
	"github.com/tricong1998/go-ecom/cmd/order/internal/repository"
	"github.com/tricong1998/go-ecom/cmd/order/internal/services"
	"gorm.io/gorm"
)

func SetupRoutes(routes *gin.Engine, db *gorm.DB, cfg *config.Config) {
	userRepo := repository.NewOrderRepository(db)
	userGateway := grpc.New(cfg.UserServer.Host, cfg.UserServer.Port)
	userService := services.NewOrderService(userRepo, userGateway)
	userHandler := handlers.NewOrderHandler(userService)

	userGroup := routes.Group("orders")
	{
		userGroup.POST("", userHandler.CreateOrder)
		userGroup.GET("/:id", userHandler.ReadOrder)
		userGroup.GET("", userHandler.ListOrders)
		userGroup.PUT("/:id", userHandler.UpdateOrder)
		userGroup.DELETE("/:id", userHandler.DeleteOrder)
	}
}

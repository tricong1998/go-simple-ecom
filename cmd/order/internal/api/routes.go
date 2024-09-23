package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"github.com/tricong1998/go-ecom/cmd/order/internal/api/handlers"
	"github.com/tricong1998/go-ecom/cmd/order/internal/config"
	"github.com/tricong1998/go-ecom/cmd/order/internal/gateway/user/grpc"
	"github.com/tricong1998/go-ecom/cmd/order/internal/repository"
	"github.com/tricong1998/go-ecom/cmd/order/internal/services"
	"github.com/tricong1998/go-ecom/pkg/rabbitmq"
	"gorm.io/gorm"
)

func SetupRoutes(
	routes *gin.Engine,
	db *gorm.DB,
	cfg *config.Config,
	rabbitCfg *rabbitmq.RabbitMQConfig,
	rabbitConn *amqp.Connection,
	log zerolog.Logger,
) {
	userRepo := repository.NewOrderRepository(db)
	userGateway := grpc.New(cfg.UserServer.Host, cfg.UserServer.Port)
	createOrderPublisher := rabbitmq.NewPublisher(
		context.Background(),
		rabbitCfg,
		rabbitConn,
		log,
		rabbitmq.E_COM_EXCHANGE,
		"direct",
		rabbitmq.PAYMENT_ORDER_COMPLETED_QUEUE,
	)
	userService := services.NewOrderService(userRepo, userGateway, createOrderPublisher)
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

package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"github.com/tricong1998/go-ecom/cmd/order/internal/api/handlers"
	"github.com/tricong1998/go-ecom/cmd/order/internal/config"
	paymentGrpc "github.com/tricong1998/go-ecom/cmd/order/internal/gateway/payment/grpc"
	productGrpc "github.com/tricong1998/go-ecom/cmd/order/internal/gateway/product/grpc"
	userGrpc "github.com/tricong1998/go-ecom/cmd/order/internal/gateway/user/grpc"
	"github.com/tricong1998/go-ecom/cmd/order/internal/repository"
	"github.com/tricong1998/go-ecom/cmd/order/internal/services"
	"github.com/tricong1998/go-ecom/pkg/gin/middleware"
	"github.com/tricong1998/go-ecom/pkg/rabbitmq"
	"github.com/tricong1998/go-ecom/pkg/token"
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
	userGateway := userGrpc.New(cfg.UserServer.Host, cfg.UserServer.Port)
	paymentGateway := paymentGrpc.New(cfg.PaymentServer.Host, cfg.PaymentServer.Port)
	productGateway := productGrpc.New(cfg.ProductServer.Host, cfg.ProductServer.Port)
	createOrderPublisher := rabbitmq.NewPublisher(
		context.Background(),
		rabbitCfg,
		rabbitConn,
		log,
		rabbitmq.E_COM_EXCHANGE,
		"direct",
		rabbitmq.PAYMENT_ORDER_COMPLETED_QUEUE,
	)
	userService := services.NewOrderService(userRepo, userGateway, paymentGateway, createOrderPublisher, productGateway)
	userHandler := handlers.NewOrderHandler(userService)

	userGroup := routes.Group("orders")
	tokenMaker, err := token.NewJWTMaker(cfg.Auth.AccessTokenSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create token maker")
		return
	}
	authRoutes := userGroup.Group("/").Use(middleware.AuthMiddleware(tokenMaker, []string{}))
	{
		authRoutes.POST("", userHandler.CreateOrder)
		authRoutes.GET("/:id", userHandler.ReadOrder)
		authRoutes.GET("", userHandler.ListOrders)
		authRoutes.PUT("/:id", userHandler.UpdateOrder)
	}
}

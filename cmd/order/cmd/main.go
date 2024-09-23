package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"github.com/tricong1998/go-ecom/cmd/order/internal/api"
	"github.com/tricong1998/go-ecom/cmd/order/internal/config"
	"github.com/tricong1998/go-ecom/cmd/order/internal/database"
	"github.com/tricong1998/go-ecom/cmd/order/pkg/logger"
	"github.com/tricong1998/go-ecom/pkg/rabbitmq"
	"gorm.io/gorm"
)

func main() {
	log := logger.NewLogger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
	}

	// Initialize db
	db, err := database.Initialize(&cfg.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize database")
	}

	// Migrate db
	err = database.Migrate(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot migrate database")
	}

	rabbitConfig := rabbitmq.RabbitMQConfig{
		Host:     cfg.RabbitMQConfig.Host,
		Port:     cfg.RabbitMQConfig.Port,
		User:     cfg.RabbitMQConfig.User,
		Password: cfg.RabbitMQConfig.Password,
	}
	rabbitConn, err := rabbitmq.NewRabbitMQConn(&rabbitConfig, context.Background())

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect rabbit")
	}

	runGinServer(cfg, db, log, &rabbitConfig, rabbitConn)
}

func runGinServer(
	cfg *config.Config,
	db *gorm.DB,
	log zerolog.Logger,
	rabbitConfig *rabbitmq.RabbitMQConfig,
	conn *amqp.Connection,
) {
	// Initialize router
	routes := gin.Default()
	api.SetupRoutes(routes, db, cfg, rabbitConfig, conn, log)

	// Start server
	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	err := routes.Run(address)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot run server")
	}
}

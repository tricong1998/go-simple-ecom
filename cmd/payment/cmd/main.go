package main

import (
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/tricong1998/go-ecom/cmd/payment/internal/api"
	"github.com/tricong1998/go-ecom/cmd/payment/internal/config"
	"github.com/tricong1998/go-ecom/cmd/payment/internal/database"
	"github.com/tricong1998/go-ecom/cmd/payment/internal/grpc_handler"
	"github.com/tricong1998/go-ecom/cmd/payment/internal/repository"
	"github.com/tricong1998/go-ecom/cmd/payment/internal/services"
	"github.com/tricong1998/go-ecom/cmd/payment/pkg/logger"
	"github.com/tricong1998/go-ecom/cmd/payment/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	go runGrpcServer(cfg, db, log)
	runGinServer(cfg, db, log)
}

func runGinServer(cfg *config.Config, db *gorm.DB, log zerolog.Logger) {
	// Initialize router
	routes := gin.Default()
	api.SetupRoutes(routes, db)

	// Start server
	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	err := routes.Run(address)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot run server")
	}
}

func runGrpcServer(cfg *config.Config, db *gorm.DB, log zerolog.Logger) {
	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := services.NewPaymentService(paymentRepo)
	server := grpc_handler.NewServer(paymentService)

	grpcServer := grpc.NewServer()
	pb.RegisterPaymentGrpcServer(grpcServer, server)
	reflection.Register(grpcServer)

	grpcServerAddress := fmt.Sprintf("%s:%s", cfg.GrpcServer.Host, cfg.GrpcServer.Port)
	listener, err := net.Listen("tcp", grpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot listen grpc server")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start grpc server")
	}
}

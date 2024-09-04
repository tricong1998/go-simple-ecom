package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tricong1998/go-ecom/internal/api"
	"github.com/tricong1998/go-ecom/internal/config"
	"github.com/tricong1998/go-ecom/internal/database"
	"github.com/tricong1998/go-ecom/pkg/logger"
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

	// Initialize router
	routes := gin.Default()
	api.SetupRoutes(routes, db)

	// Start server
	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	err = routes.Run(address)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot run server")
	}
}

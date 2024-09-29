package main

import (
	"fmt"
	"os"

	"github.com/tricong1998/go-ecom/cmd/user/internal/config"
	"github.com/tricong1998/go-ecom/cmd/user/internal/database"
	"github.com/tricong1998/go-ecom/cmd/user/internal/repository"
	"github.com/tricong1998/go-ecom/cmd/user/internal/util"
	"github.com/tricong1998/go-ecom/cmd/user/pkg/logger"
	"github.com/tricong1998/go-ecom/cmd/user/pkg/models"
)

func main() {
	args := os.Args[1:]

	// Check if arguments were provided
	if len(args) != 2 {
		fmt.Println("Invalid number of arguments")
		return
	}

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

	username := args[0]
	password := args[1]
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot hash password")
	}

	user := models.User{
		Username: username,
		Password: hashedPassword,
		Role:     string(models.AdminRole),
	}
	userRepo := repository.NewUserRepository(db)
	userRepo.CreateUser(&user)
	log.Info().Msg("Admin account created successfully")
}

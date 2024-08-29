package database

import (
	"fmt"

	"github.com/tricong1998/go-ecom/internal/config"
	"github.com/tricong1998/go-ecom/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Initialize(dbConfig *config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbConfig.DBHost,
		dbConfig.DBUser,
		dbConfig.DBPassword,
		dbConfig.DBName,
		dbConfig.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Product{},
		&models.User{},
		&models.Order{},
		// Add other models here as needed
	)
}

package config

import (
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

type ServerConfig struct {
	Host string
	Port string
}

type Config struct {
	Server ServerConfig
	DB     DBConfig
}

func Load() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
			Host: os.Getenv("SERVER_HOST"),
		},
		DB: DBConfig{
			DBHost:     os.Getenv("DB_HOST"),
			DBPort:     os.Getenv("DB_PORT"),
			DBUser:     os.Getenv("DB_USER"),
			DBPassword: os.Getenv("DB_PASSWORD"),
			DBName:     os.Getenv("DB_NAME"),
		},
	}

	if config.Server.Port != "" {
		config.Server.Port = "3333"
	}

	if config.Server.Host != "" {
		config.Server.Host = "localhost"
	}

	// "postgresql://user:password@localhost:5432/simple-go-ecom?sslmode=disable"
	if config.DB.DBHost != "" {
		config.DB.DBHost = "localhost"
	}

	if config.DB.DBPort != "" {
		config.DB.DBPort = "5432"
	}

	if config.DB.DBName != "" {
		config.DB.DBName = "simple-ecom"
	}

	return config, nil
}

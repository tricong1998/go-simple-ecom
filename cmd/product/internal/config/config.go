package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tricong1998/go-ecom/pkg/util"
)

type DBConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

type HttpServerConfig struct {
	Host string
	Port string
}

type GrpcServerConfig struct {
	Host string
	Port string
}

type AuthConfig struct {
	AccessTokenDuration  time.Duration
	AccessTokenSecret    string
	RefreshTokenSecret   string
	RefreshTokenDuration time.Duration
}

type Config struct {
	Server     HttpServerConfig
	GrpcServer GrpcServerConfig
	DB         DBConfig
	Auth       AuthConfig
}

func Load() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{
		Server: HttpServerConfig{
			Port: os.Getenv("PRODUCT_SERVER_PORT"),
			Host: os.Getenv("PRODUCT_SERVER_HOST"),
		},
		GrpcServer: GrpcServerConfig{
			Host: os.Getenv("PRODUCT_GRPC_SERVER_HOST"),
			Port: os.Getenv("PRODUCT_GRPC_SERVER_PORT"),
		},
		DB: DBConfig{
			DBHost:     os.Getenv("DB_HOST"),
			DBPort:     os.Getenv("DB_PORT"),
			DBUser:     os.Getenv("DB_USER"),
			DBPassword: os.Getenv("DB_PASSWORD"),
			DBName:     os.Getenv("PRODUCT_DB_NAME"),
		},
		Auth: AuthConfig{
			AccessTokenSecret:    os.Getenv("ACCESS_TOKEN_SECRET"),
			RefreshTokenSecret:   os.Getenv("REFRESH_TOKEN_SECRET"),
			AccessTokenDuration:  util.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"), 15*time.Minute),
			RefreshTokenDuration: util.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"), 24*time.Hour),
		},
	}

	if config.Server.Port == "" {
		config.Server.Port = "3333"
	}

	if config.Server.Host == "" {
		config.Server.Host = "localhost"
	}

	// "postgresql://user:password@localhost:5432/simple-go-ecom?sslmode=disable"
	if config.DB.DBHost == "" {
		config.DB.DBHost = "localhost"
	}

	if config.DB.DBPort == "" {
		config.DB.DBPort = "5432"
	}

	if config.DB.DBName == "" {
		config.DB.DBName = "simple-ecom"
	}

	return config, nil
}

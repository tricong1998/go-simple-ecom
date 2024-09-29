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

type ServerConfig struct {
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
	Server         ServerConfig
	UserServer     ServerConfig
	PaymentServer  ServerConfig
	ProductServer  ServerConfig
	DB             DBConfig
	RabbitMQConfig RabbitMQConfig
	Auth           AuthConfig
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

func Load() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{
		Server: ServerConfig{
			Port: os.Getenv("ORDER_SERVER_PORT"),
			Host: os.Getenv("ORDER_SERVER_HOST"),
		},
		DB: DBConfig{
			DBHost:     os.Getenv("DB_HOST"),
			DBPort:     os.Getenv("DB_PORT"),
			DBUser:     os.Getenv("DB_USER"),
			DBPassword: os.Getenv("DB_PASSWORD"),
			DBName:     os.Getenv("DB_NAME"),
		},
		UserServer: ServerConfig{
			Port: os.Getenv("USER_GRPC_SERVER_PORT"),
			Host: os.Getenv("ORDER_USER_GRPC_SERVER_HOST"),
		},
		PaymentServer: ServerConfig{
			Port: os.Getenv("PAYMENT_GRPC_SERVER_PORT"),
			Host: os.Getenv("ORDER_PAYMENT_GRPC_SERVER_HOST"),
		},
		ProductServer: ServerConfig{
			Port: os.Getenv("PRODUCT_GRPC_SERVER_PORT"),
			Host: os.Getenv("ORDER_PRODUCT_GRPC_SERVER_HOST"),
		},
		RabbitMQConfig: RabbitMQConfig{
			Port:     os.Getenv("AMQP_SERVER_PORT"),
			Host:     os.Getenv("AMQP_SERVER_HOST"),
			User:     os.Getenv("AMQP_SERVER_USER"),
			Password: os.Getenv("AMQP_SERVER_PASSWORD"),
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

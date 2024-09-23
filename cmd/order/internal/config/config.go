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
	Server         ServerConfig
	UserServer     ServerConfig
	DB             DBConfig
	RabbitMQConfig RabbitMQConfig
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
			Host: os.Getenv("USER_GRPC_SERVER_HOST"),
		},
		RabbitMQConfig: RabbitMQConfig{
			Port:     os.Getenv("AMQP_SERVER_PORT"),
			Host:     os.Getenv("AMQP_SERVER_HOST"),
			User:     os.Getenv("AMQP_SERVER_USER"),
			Password: os.Getenv("AMQP_SERVER_PASSWORD"),
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

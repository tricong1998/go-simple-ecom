package rabbit_handler

import (
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"github.com/tricong1998/go-ecom/cmd/user/internal/services"
	"github.com/tricong1998/go-ecom/cmd/user/pkg/dto"
)

type CreateUserPointDependencies struct {
	UserService *services.UserService
	Logger      zerolog.Logger
}

func CreateUserPoint(queue string, msg amqp.Delivery, dependencies *CreateUserPointDependencies) error {
	dependencies.Logger.Info().Msgf("Message received on queue: %s with message: %s", queue, string(msg.Body))

	var productCreated dto.CreateUserPoint

	err := json.Unmarshal(msg.Body, &productCreated)
	if err != nil {
		return err
	}

	err = dependencies.UserService.CreateUserPoint(productCreated)
	return err
}

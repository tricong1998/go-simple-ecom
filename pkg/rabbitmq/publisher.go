package rabbitmq

import (
	"context"
	"reflect"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/iancoleman/strcase"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

//go:generate mockery --name IPublisher
type IPublisher interface {
	PublishMessage(msg interface{}) error
	IsPublished(msg interface{}) bool
}

var publishedMessages []string

type Publisher struct {
	cfg          *RabbitMQConfig
	conn         *amqp.Connection
	log          zerolog.Logger
	ctx          context.Context
	exchangeName string
	exchangeType string
	queueName    string
}

func (p Publisher) PublishMessage(msg interface{}) error {

	data, err := jsoniter.Marshal(msg)

	if err != nil {
		p.log.Fatal().Err(err).Msg("Error in marshalling message to publish message")
		return err
	}

	// Inject the context in the headers
	// headers := otel.InjectAMQPHeaders(ctx)

	channel, err := p.conn.Channel()
	if err != nil {
		p.log.Fatal().Err(err).Msg("Error in opening channel to consume message")
		return err
	}

	defer channel.Close()

	err = channel.ExchangeDeclare(
		p.exchangeName, // name
		p.exchangeType, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)

	if err != nil {
		p.log.Fatal().Err(err).Msg("Error in declaring exchange to publish message")
		return err
	}

	correlationId := ""

	publishingMsg := amqp.Publishing{
		Body:          data,
		ContentType:   "application/json",
		DeliveryMode:  amqp.Persistent,
		MessageId:     uuid.NewV4().String(),
		Timestamp:     time.Now(),
		CorrelationId: correlationId,
		// Headers:       headers,
	}

	err = channel.Publish(p.exchangeName, p.queueName, false, false, publishingMsg)

	if err != nil {
		p.log.Fatal().Err(err).Msg("Error in publishing message")
		return err
	}

	publishedMessages = append(publishedMessages, p.queueName)

	// h, err := jsoniter.Marshal(headers)

	if err != nil {
		p.log.Fatal().Err(err).Msg("Error in marshalling headers to publish message")
		return err
	}

	p.log.Info().Msgf("Published message: %s", publishingMsg.Body)

	return nil
}

func (p Publisher) IsPublished(msg interface{}) bool {

	typeName := reflect.TypeOf(msg).Name()
	snakeTypeName := strcase.ToSnake(typeName)
	isPublished := linq.From(publishedMessages).Contains(snakeTypeName)

	return isPublished
}

func NewPublisher(
	ctx context.Context,
	cfg *RabbitMQConfig,
	conn *amqp.Connection,
	log zerolog.Logger,
	exchangeName string,
	exchangeType string,
	queueName string,
) IPublisher {
	return &Publisher{
		ctx:          ctx,
		cfg:          cfg,
		conn:         conn,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		queueName:    queueName,
	}
}

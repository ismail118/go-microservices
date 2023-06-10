package event

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitEmitter struct {
	connection *amqp.Connection
}

func NewRabbitEmitter(conn *amqp.Connection) *RabbitEmitter {
	emitter := &RabbitEmitter{
		connection: conn,
	}

	return emitter
}

func (e *RabbitEmitter) Push(event string, severity string) error {
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	log.Println("Push to channel")

	err = ch.PublishWithContext(context.Background(),
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (e *RabbitEmitter) SetupEvenEmitter() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	err = declareExchange(channel)
	if err != nil {
		return err
	}

	return nil
}

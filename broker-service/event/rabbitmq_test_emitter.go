package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitTestEmitter struct {
	connection *amqp.Connection
}

func (e *RabbitTestEmitter) Push(event string, severity string) error {
	return nil
}

func (e *RabbitTestEmitter) SetupEvenEmitter() error {
	return nil
}

func NewRabbitTestEmitter(conn *amqp.Connection) *RabbitTestEmitter {
	emitter := &RabbitTestEmitter{
		connection: conn,
	}

	return emitter
}

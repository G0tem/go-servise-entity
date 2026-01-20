package queue

import (
	"github.com/G0tem/go-servise-entity/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func autocreateQueue(conn *amqp.Connection, cfg *config.Config) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	// Always close channel after queue redeclare (otherwise first call delivery/consume fail)
	defer ch.Close()

	_, err = ch.QueueDeclare(
		cfg.RMQConsumeQ, // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	return err
}

func autocreateExchange(conn *amqp.Connection, exchangeName string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	return err
}

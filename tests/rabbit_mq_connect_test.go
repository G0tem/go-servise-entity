package tests

import (
	"testing"

	"github.com/G0tem/go-servise-entity/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TestRabbitMqConnect(t *testing.T) {
	cfg := config.LoadConfig()

	conn, err := amqp.Dial(cfg.RMQConnUrl)
	failOnError(t, err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(t, err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		cfg.RMQConsumeQ, // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(t, err, "Failed to declare a queue")

	// Create exchange

	err = ch.ExchangeDeclare(
		cfg.RMQExchange, // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(t, err, "Failed to create exchange")

}

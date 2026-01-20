package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/G0tem/go-servise-entity/internal/config"
	"github.com/G0tem/go-servise-entity/internal/dto"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TestRabbitSendSpy(t *testing.T) {
	cfg := config.LoadConfig()

	conn, err := amqp.Dial(cfg.RMQConnUrl)
	failOnError(t, err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(t, err, "Failed to open a channel")
	defer ch.Close()

	// Create exchange first
	err = ch.ExchangeDeclare(
		cfg.RMQExchange, // name
		"direct",        // type direct fanout
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(t, err, "Failed to create exchange")

	// Начало недели: 2025-06-16 00:00:00 UTC
	DateFrom := time.Now()
	fmt.Println(DateFrom)

	body, err := json.Marshal(dto.EntityMessage{
		UserId:        "test-user-id-123",
		Timestamp:     DateFrom.UnixNano() / 1e6,
		MessageParams: "update",
	})
	failOnError(t, err, "Failed to json serialization")

	err = ch.Publish(
		cfg.RMQExchange, // exchange
		cfg.RMQConsumeB, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
			Timestamp:   time.Now(),
		})
	failOnError(t, err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

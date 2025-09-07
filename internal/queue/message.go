package queue

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/G0tem/go-servise-entity/internal/dto"
	"github.com/rs/zerolog/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ListenStreamingQueue(rabbitURL string) error {
	log.Info().Msgf("Setup ListenStreamingQueue ResultSFMessage")
	for {
		conn, err := amqp.Dial(rabbitURL)
		if err != nil {
			log.Error().Msgf("Failed to connect to RabbitMQ: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			log.Error().Msgf("Failed to open channel: %v. Reconnecting...", err)
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		// Подключаемся к очереди
		msgs, err := ch.Consume(
			"queue_name",    // queue
			"consumer_name", // consumer
			true,            // auto-ack
			false,           // exclusive
			false,           // no-local
			false,           // no-wait
			nil,             // args
		)
		if err != nil {
			log.Error().Msgf("Failed to consume: %v. Reconnecting...", err)
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		log.Info().Msgf("Connected to queue 'streaming_incoming', waiting for messages...")

		// Флаг: закрыта ли соединение
		connected := true

		// Читаем сообщения
		for msg := range msgs {
			var result dto.NotifyMessage
			if err := json.Unmarshal(msg.Body, &result); err != nil {
				log.Error().Msgf("Failed to unmarshal message: %v", err)
				continue
			}

			log.Printf("Received: UserID=%s, Email=%v, Status=%s",
				result.UserId, result.NotifyTimestamp, result.MessageParams)

			fmt.Println("тут будет запуск логики из в эндпоинте")

		}

		if connected {
			log.Warn().Msgf("Connection lost. Reconnecting...")
			connected = false
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
		}
	}
}

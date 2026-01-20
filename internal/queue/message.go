package queue

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/G0tem/go-servise-entity/internal/config"
	"github.com/G0tem/go-servise-entity/internal/dto"
	"github.com/rs/zerolog/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func printLogo() {
	fmt.Printf("go-servise-entity SERVICE V1.0")
}

func ListenRabbitQueue(cfg *config.Config) error {

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(
		signalChannel,
		syscall.SIGUSR1, // Use for restart listening queue
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGSEGV,
	)

	log.Info().Msgf("Setup listening queue %v, exchange: %v, binding: %v", cfg.RMQConsumeQ, cfg.RMQExchange, cfg.RMQConsumeB)

	go processQueue(cfg)

	for {
		signalEvent := <-signalChannel
		switch signalEvent {
		case syscall.SIGUSR1:
			time.Sleep(5 * time.Second)
			go processQueue(cfg)
		case syscall.SIGQUIT,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGKILL:
			log.Error().Msgf("Signal event %q", signalEvent)
			return nil
		case syscall.SIGHUP:
			log.Error().Msgf("Signal event %q", signalEvent)
			return fmt.Errorf("signal hang up")
		case syscall.SIGSEGV:
			log.Error().Msgf("Signal event %q", signalEvent)
			return fmt.Errorf("segmentation violation")
		default:
			log.Error().Msgf("Unexpected signal %q", signalEvent)
		}
	}
}

func processQueue(cfg *config.Config) {
	for {
		conn, err := amqp.Dial(cfg.RMQConnUrl)
		if err != nil {
			log.Error().Msgf("Failed to connect to RabbitMQ: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		notifyClose := make(chan *amqp.Error)
		conn.NotifyClose(notifyClose)

		chProfileSenders := make(chan *amqp.Delivery, 11)

		activeGoroutines := 0
		var mu sync.Mutex

		go func() {
			for message := range chProfileSenders {
				mu.Lock()
				activeGoroutines++
				currentCount := activeGoroutines
				mu.Unlock()

				log.Info().Msgf("Active goroutines: %d", currentCount)

				go func(message *amqp.Delivery) {
					defer func() {
						mu.Lock()
						activeGoroutines--
						mu.Unlock()
					}()

					var entity_message dto.EntityMessage

					err := json.Unmarshal(message.Body, &entity_message)
					if err != nil {
						log.Error().Msgf("Can't unmarshal body: %v", err)
						message.Nack(false, false)
						log.Warn().Msgf("Error parsing message, skipped. Message user: %v Message date: %v, Message: %v, to message.Nack(false, false)", entity_message.UserId, entity_message.Timestamp, entity_message.MessageParams)
						return
					}
					// message.Nack(false, false)
					// Вызов обработчика сообщения
					log.Debug().Msgf("get message package queue: %s", entity_message.UserId)

					err = message.Ack(false)
					if err != nil {
						log.Error().Msgf("Acknowledge of message fail with error %v", err)
						return
					}
				}(message)
			}
		}()

		messages, err := getMessagesChannel(conn, cfg)
		if err != nil {
			log.Error().Msgf("Failed to get messages channel: %v", err)
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		printLogo()

		for {
			select {
			case err := <-notifyClose:
				log.Error().Msgf("RabbitMQ connection closed: %v", err)
				conn.Close()
				time.Sleep(5 * time.Second)
				goto reconnect
			case message, ok := <-messages:
				if !ok {
					log.Error().Msg("Messages channel closed")
					conn.Close()
					time.Sleep(5 * time.Second)
					goto reconnect
				}

				log.Debug().Msgf("Received new message from RabbitMQ")

				select {
				case chProfileSenders <- &message:
					log.Debug().Msgf("Message sent to processing channel")
				default:
					log.Warn().Msgf("Profile senders channel is full, message will be requeued")
					message.Nack(false, true)
				}
			}
		}

	reconnect:
		log.Info().Msg("Attempting to reconnect to RabbitMQ...")
	}
}

func getMessagesChannel(conn *amqp.Connection, cfg *config.Config) (<-chan amqp.Delivery, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	notifyClose := make(chan *amqp.Error)
	ch.NotifyClose(notifyClose)

	go func() {
		err := <-notifyClose
		log.Error().Msgf("Channel closed: %v", err)
	}()

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		ch.Close()
		return nil, fmt.Errorf("failed to set QoS: %v", err)
	}

	log.Info().Msg("Channel QoS set successfully")

	queue, err := ch.QueueDeclarePassive(
		cfg.RMQConsumeQ, //name
		true,            // durable
		false,           // autoDelete
		false,           // exclusive
		false,           // noWait
		nil,             // arguments
	)
	if err != nil && cfg.RMQConsumeQAutocreate {
		log.Warn().Msgf("Queue %q doesn't exist, autocreate enabled, so attempting to create queue", cfg.RMQConsumeQ)

		err = autocreateQueue(conn, cfg)
		if err != nil {
			ch.Close()
			return nil, err
		}
		ch, err = conn.Channel()
		if err != nil {
			return nil, err
		}
		err = ch.Qos(1, 0, false)
		if err != nil {
			ch.Close()
			return nil, fmt.Errorf("failed to set QoS: %v", err)
		}
	}

	log.Info().Msgf("Queue %s exists with %d messages", queue.Name, queue.Messages)

	err = ch.ExchangeDeclarePassive(
		cfg.RMQExchange, // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Warn().Msgf("Exchange %q doesn't exist (%v)", cfg.RMQExchange, err)
		autocreateExchange(conn, cfg.RMQExchange)
		log.Info().Msgf("Exchange %q create", cfg.RMQExchange)
	}

	err = ch.QueueBind(
		cfg.RMQConsumeQ,
		cfg.RMQConsumeB,
		cfg.RMQExchange,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		return nil, err
	}

	log.Info().Msg("Queue bound to exchange successfully")

	msgs, err := ch.Consume(
		cfg.RMQConsumeQ,
		"servise-entity-backend",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		ch.Close()
		return nil, err
	}

	log.Info().Msg("Successfully started consuming messages")
	return msgs, nil
}

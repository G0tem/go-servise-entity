package queue

import (
	"encoding/json"
	"time"

	"github.com/G0tem/go-servise-entity/internal/config"
	"github.com/G0tem/go-servise-entity/internal/dto"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

func PushMessage(userId string, message interface{}, targetUserId string, cfg *config.Config) error {

	body, err := json.Marshal(dto.NotifyMessage{
		UserId:          userId,
		NotifyTimestamp: time.Now().UnixNano() / 1e6,
		MessageParams:   message,
		TargetUserId:    targetUserId,
	})
	if err != nil {
		return err
	}
	conn, err := amqp.Dial(cfg.RMQConnUrl)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := getChannel(conn, cfg)
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Publish(
		cfg.RMQNotifyExchange,   // exchange
		cfg.RMQNotifyRoutingKey, // routing key
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
			Timestamp:   time.Now(),
		})
	if err != nil {
		return err
	}
	log.Info().Msgf("Message output to: %s", cfg.RMQNotifyExchange)
	return nil
}

func getChannel(conn *amqp.Connection, cfg *config.Config) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	// No defer, channel will automatically closed on error

	err = ch.ExchangeDeclarePassive(
		cfg.RMQNotifyExchange, // name
		"fanout",              // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil && cfg.RMQNotifyExchangeAutocreate {
		log.Warn().Msgf("Exchange %q doesn't exist, autocreate enabled, so attempting to create exchange", cfg.RMQNotifyExchange)
		autocreateExchange(conn, cfg.RMQNotifyExchange)

		ch, err = conn.Channel()
		if err != nil {
			return nil, err
		}
	}

	return ch, nil
}

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
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	return err
}

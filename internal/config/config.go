package config

import (
	"os"
	"strconv"
	"time"

	"github.com/G0tem/go-servise-entity/internal"
)

type Config struct {
	LogLevel  int    `default:"4" envconfig:"LOG_LEVEL"`
	HttpPort  uint16 `default:"8007" envconfig:"HTTP_PORT"`
	SecretKey string `binding:"required" envconfig:"SECRET_KEY"`

	AuthGrpcAddress string `default:"localhost:50051" envconfig:"AUTH_GRPC_ADDRESS"`

	RMQConnUrl            string `binding:"required" envconfig:"RMQ_CONN_URL"`
	RMQConsumeQ           string `binding:"required" envconfig:"RMQ_QUEUE_CONSUME"`
	RMQConsumeB           string `binding:"required" envconfig:"RMQ_BINDING_CONSUME"`
	RMQExchange           string `binding:"required" envconfig:"RMQ_EXCHANGE"`
	RMQConsumeQAutocreate bool   `binding:"required" envconfig:"RMQ_QUEUE_CONSUME_AUTOCREATE_DISABLED"`

	RMQNotifyExchange           string `binding:"required" envconfig:"RMQ_NOTIFY_EXCHANGE"`
	RMQNotifyRoutingKey         string `binding:"required" envconfig:"RMQ_NOTIFY_ROUTING_KEY"`
	RMQNotifyExchangeAutocreate bool   `binding:"required" envconfig:"RMQ_NOTIFY_EXCHANGE_AUTOCREATE_ENABLED"`

	PostgresHost            string        `binding:"required" envconfig:"POSTGRES_HOST"`
	PostgresPort            string        `binding:"required" envconfig:"POSTGRES_PORT"`
	PostgresDb              string        `binding:"required" envconfig:"POSTGRES_DB"`
	PostgresUser            string        `binding:"required" envconfig:"POSTGRES_USER"`
	PostgresPassword        string        `binding:"required" envconfig:"POSTGRES_PASSWORD"`
	PostgresMaxIdleConns    int           `default:"10" envconfig:"POSTGRES_MAX_IDLE_CONNS"`
	PostgresMaxOpenConns    int           `default:"100" envconfig:"POSTGRES_MAX_OPEN_CONNS"`
	PostgresConnMaxLifetime time.Duration `default:"1h" envconfig:"POSTGRES_CONN_MAX_LIFETIME"`

	// RedisAddr string `binding:"required" envconfig:"REDIS_ADDR"`
}

func LoadConfig() Config {

	logLevel, _ := strconv.Atoi(os.Getenv("LOG_LEVEL"))

	return Config{
		LogLevel:        logLevel,
		HttpPort:        internal.ParseUint16(os.Getenv("HTTP_PORT"), 8010),
		SecretKey:       os.Getenv("SECRET_KEY"),
		AuthGrpcAddress: internal.Getenv("AUTH_GRPC_ADDRESS", "localhost:50051"),

		RMQConnUrl:            os.Getenv("RMQ_CONN_URL"),
		RMQConsumeQ:           os.Getenv("RMQ_QUEUE_CONSUME"),
		RMQConsumeB:           internal.Getenv("RMQ_BINDING_CONSUME", "entity.*"),
		RMQExchange:           internal.Getenv("RMQ_EXCHANGE", "entity_exchange"),
		RMQConsumeQAutocreate: internal.ParseBool(os.Getenv("RMQ_QUEUE_CONSUME_AUTOCREATE_DISABLED")),

		RMQNotifyExchange:           os.Getenv("RMQ_NOTIFY_EXCHANGE"),
		RMQNotifyRoutingKey:         internal.Getenv("RMQ_NOTIFY_ROUTING_KEY", "notify.*"),
		RMQNotifyExchangeAutocreate: internal.ParseBool(internal.Getenv("RMQ_NOTIFY_EXCHANGE_AUTOCREATE_ENABLED", "Y")),

		PostgresHost:            os.Getenv("POSTGRES_HOST"),
		PostgresPort:            os.Getenv("POSTGRES_PORT"),
		PostgresDb:              os.Getenv("POSTGRES_DB"),
		PostgresUser:            os.Getenv("POSTGRES_USER"),
		PostgresPassword:        os.Getenv("POSTGRES_PASSWORD"),
		PostgresMaxIdleConns:    internal.ParseInt(os.Getenv("POSTGRES_MAX_IDLE_CONNS"), 10),
		PostgresMaxOpenConns:    internal.ParseInt(os.Getenv("POSTGRES_MAX_OPEN_CONNS"), 100),
		PostgresConnMaxLifetime: internal.ParseDuration(os.Getenv("POSTGRES_CONN_MAX_LIFETIME"), 1*time.Hour),

		// RedisAddr: os.Getenv("REDIS_ADDR"),
	}
}

func GetenvDef(envVariable string, defaultValue string) string {
	result := os.Getenv(envVariable)
	if result != "" {
		return result
	}
	return defaultValue
}

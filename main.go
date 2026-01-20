package main

import (
	"os"
	"runtime"

	_ "github.com/lib/pq"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/G0tem/go-servise-entity/internal/config"
	"github.com/G0tem/go-servise-entity/internal/queue"
	"github.com/G0tem/go-servise-entity/internal/service/factory"
)

// @title Local-Template-Entity Swagger
// @version 1.0
// @description This is an API of payment-service for Socialweb application

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Description for what is this security definition being used

// @BasePath /api/v1
func main() {
	// Initialize Zerolog logger with output to stdout
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := config.LoadConfig()
	zerolog.SetGlobalLevel(zerolog.Level(cfg.LogLevel))

	go queue.ListenRabbitQueue(&cfg)

	// We are need >= 2 threads
	moreThenTwoThreadsRuntime()

	err := factory.StartHttpService(&cfg)

	if err != nil {
		log.Error().Msgf("Attempt to start application fail with error %v", err)
	}
}

func moreThenTwoThreadsRuntime() {
	currentThreadsCount := runtime.GOMAXPROCS(2)
	if currentThreadsCount > 2 {
		runtime.GOMAXPROCS(currentThreadsCount)
	}
}

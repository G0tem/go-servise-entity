package handler

import (
	"net/http"
	"time"

	"github.com/G0tem/go-servise-entity/internal/config"
	grpcClient "github.com/G0tem/go-servise-entity/internal/grpc"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gorm.io/gorm"
)

func NewHandler(db *gorm.DB, cfg *config.Config) *Handler {
	// Инициализируем gRPC клиент
	authClient, err := grpcClient.NewAuthClient(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create auth gRPC client")
		// Продолжаем работу без gRPC клиента, но логируем ошибку
	}

	return &Handler{
		db:         db,
		cfg:        cfg,
		authClient: authClient,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type Handler struct {
	db         *gorm.DB
	cfg        *config.Config
	authClient *grpcClient.AuthClient
	client     *http.Client
}

func (h *Handler) SetupRoutes(app *fiber.App) {
	cfg := config.LoadConfig()

	api := app.Group("/api")
	v1 := api.Group("/v1")

	entity := v1.Group("/entity")

	entity.Use(JWTMiddleware(cfg.SecretKey))

	entity.Get("user_info", h.CheckUser)
	entity.Post("create", h.CreateEntity)
	entity.Get("test_grpc", h.TestGrpc)
}

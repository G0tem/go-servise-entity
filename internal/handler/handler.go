package handler

import (
	"net/http"
	"time"

	"github.com/G0tem/go-servise-entity/internal/config"
	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

func NewHandler(db *gorm.DB, cfg *config.Config) *Handler {
	return &Handler{
		db:  db,
		cfg: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type Handler struct {
	db     *gorm.DB
	cfg    *config.Config
	client *http.Client
}

func (h *Handler) SetupRoutes(app *fiber.App) {
	cfg := config.LoadConfig()

	api := app.Group("/api")
	v1 := api.Group("/v1")

	entity := v1.Group("/entity")

	entity.Use(JWTMiddleware(cfg.SecretKey))

	entity.Get("user_info", h.CheckUser)
	entity.Post("create", h.CreateEntity)
}

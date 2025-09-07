package handler

import (
	"github.com/G0tem/go-servise-entity/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// CheckUser godoc
// @Summary Get user info
// @Description Test endpoint
// @Tags info
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} types.GetMeResponse
// @Failure 400 {object} types.FailureResponse
// @Failure 500 {object} types.FailureErrorResponse
// @Router /entity/user_info [get]
func (h *Handler) CheckUser(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*types.JwtClaims)
	log.Debug().
		Str("email", claims.Email).
		Time("exp", claims.Exp).
		Msg("Attempting to get user")

	return c.Status(fiber.StatusOK).JSON(types.GetMeResponse{
		ID:    claims.UserID,
		Email: claims.Email,
	})
}

// CreateEntity godoc
// @Summary CreateEntity info
// @Description Test endpoint
// @Tags info
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} types.GetMeResponse
// @Failure 400 {object} types.FailureResponse
// @Failure 500 {object} types.FailureErrorResponse
// @Router /entity/create [post]
func (h *Handler) CreateEntity(c *fiber.Ctx) error {
	// Логика Изменения записи в кеше
	log.Info().Msg("Start CreateEntity")

	claims := c.Locals("claims").(*types.JwtClaims)
	log.Debug().
		Str("email", claims.Email).
		Time("exp", claims.Exp).
		Msg("Attempting to get user")

	return c.Status(fiber.StatusOK).JSON(types.GetMeResponse{
		ID:    claims.UserID,
		Email: claims.Email,
	})
}

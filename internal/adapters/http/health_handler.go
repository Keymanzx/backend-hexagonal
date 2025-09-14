package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "backend-hexagonal",
		"version":   "1.0.0",
	})
}

func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	// Add any readiness checks here (database connectivity, etc.)
	return c.JSON(fiber.Map{
		"status": "ready",
		"checks": fiber.Map{
			"database": "ok",
			"auth":     "ok",
		},
	})
}

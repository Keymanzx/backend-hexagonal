package http

import (
	"backend-hexagonal/internal/adapters/http/middleware"
	"backend-hexagonal/internal/config"
	"backend-hexagonal/internal/service"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, userHandler *UserHandler, authHandler *AuthHandler, authService *service.AuthService) {
	// Apply logging middleware to all routes
	if config.IsJSONLogging() {
		app.Use(middleware.JSONLoggingMiddleware())
	} else if config.IsDetailedLogging() {
		app.Use(middleware.DetailedLoggingMiddleware())
	} else {
		app.Use(middleware.LoggingMiddleware())
	}

	// Health check endpoints (no auth required)
	healthHandler := NewHealthHandler()
	app.Get("/health", healthHandler.Health)
	app.Get("/ready", healthHandler.Ready)

	api := app.Group("/api/v1")

	// Public auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected user routes
	users := api.Group("/users")
	users.Use(middleware.JWTMiddleware(authService)) // Apply JWT middleware to all user routes
	users.Get("/me", userHandler.GetMe)
	users.Post("/", userHandler.Create)
	users.Get("/", userHandler.List)
	users.Get("/:id", userHandler.Get)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)
}

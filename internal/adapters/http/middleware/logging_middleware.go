package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoggingMiddleware logs HTTP method, path, and execution time
func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Record start time
		start := time.Now()

		// Get request info
		method := c.Method()
		path := c.Path()
		ip := c.IP()

		// Process request
		err := c.Next()

		// Calculate execution time
		duration := time.Since(start)

		// Get response status
		status := c.Response().StatusCode()

		// Log the request
		log.Printf("[%s] %s %s - %d - %v - %s",
			method,
			path,
			ip,
			status,
			duration,
			time.Now().Format("2006-01-02 15:04:05"),
		)

		return err
	}
}

// DetailedLoggingMiddleware provides more detailed logging with request/response info
func DetailedLoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Request info
		method := c.Method()
		path := c.Path()
		ip := c.IP()
		userAgent := c.Get("User-Agent")
		contentLength := len(c.Body())

		// Process request
		err := c.Next()

		// Response info
		duration := time.Since(start)
		status := c.Response().StatusCode()
		responseSize := len(c.Response().Body())

		// Determine log level based on status code
		logLevel := "INFO"
		if status >= 400 && status < 500 {
			logLevel = "WARN"
		} else if status >= 500 {
			logLevel = "ERROR"
		}

		// Log with detailed information
		log.Printf("[%s] %s %s %s - Status: %d - Duration: %v - ReqSize: %d bytes - RespSize: %d bytes - UA: %s - Time: %s",
			logLevel,
			method,
			path,
			ip,
			status,
			duration,
			contentLength,
			responseSize,
			userAgent,
			time.Now().Format("2006-01-02 15:04:05"),
		)

		return err
	}
}

// JSONLoggingMiddleware logs in JSON format for structured logging
func JSONLoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate execution time
		duration := time.Since(start)

		// Create structured log entry
		logEntry := map[string]interface{}{
			"timestamp":     time.Now().Format(time.RFC3339),
			"method":        c.Method(),
			"path":          c.Path(),
			"status":        c.Response().StatusCode(),
			"duration_ms":   duration.Milliseconds(),
			"ip":            c.IP(),
			"user_agent":    c.Get("User-Agent"),
			"request_size":  len(c.Body()),
			"response_size": len(c.Response().Body()),
			"query_params":  c.Request().URI().QueryArgs().String(),
		}

		// Add user info if available (from JWT middleware)
		if userID := c.Locals("user_id"); userID != nil {
			logEntry["user_id"] = userID
		}
		if email := c.Locals("email"); email != nil {
			logEntry["email"] = email
		}

		// Convert to JSON-like log format
		log.Printf(`{"level":"%s","msg":"HTTP Request","data":%+v}`,
			getLogLevel(c.Response().StatusCode()),
			logEntry,
		)

		return err
	}
}

// getLogLevel determines log level based on HTTP status code
func getLogLevel(status int) string {
	switch {
	case status >= 500:
		return "ERROR"
	case status >= 400:
		return "WARN"
	default:
		return "INFO"
	}
}

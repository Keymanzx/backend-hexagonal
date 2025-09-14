package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func MongoURI() string {
	if v := os.Getenv("MONGO_URI"); v != "" {
		return v
	}
	return "mongodb://localhost:27017"
}

func Port() string {
	if v := os.Getenv("PORT"); v != "" {
		return ":" + v
	}
	return ":3000"
}

func DBName() string {
	if v := os.Getenv("DB_NAME"); v != "" {
		return v
	}
	return "appdb"
}

func JWTSecret() string {
	if v := os.Getenv("JWT_SECRET"); v != "" {
		return v
	}
	return "default-secret-change-this"
}

func LogLevel() string {
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		return v
	}
	return "INFO"
}

func IsDetailedLogging() bool {
	return os.Getenv("DETAILED_LOGGING") == "true"
}

func IsJSONLogging() bool {
	return os.Getenv("JSON_LOGGING") == "true"
}

func GRPCPort() string {
	if v := os.Getenv("GRPC_PORT"); v != "" {
		return ":" + v
	}
	return ":9000"
}

// LoadEnv loads environment variables from .env file
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}
}

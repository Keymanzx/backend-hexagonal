package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"backend-hexagonal/internal/adapters/http"
	mongoadapter "backend-hexagonal/internal/adapters/mongo"
	"backend-hexagonal/internal/config"
	"backend-hexagonal/internal/service"
)

func main() {
	// Load environment variables from .env file
	config.LoadEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI()))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(config.DBName())

	// setup repository -> service -> handler
	userRepo := mongoadapter.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(userRepo)

	userHandler := http.NewUserHandler(userSvc)
	authHandler := http.NewAuthHandler(authSvc)

	app := fiber.New()
	http.RegisterRoutes(app, userHandler, authHandler, authSvc)

	// optional: background goroutine example: log user count every 10s
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			users, err := userSvc.GetAllUsers(context.Background())
			if err != nil {
				log.Println("background user list error:", err)
				continue
			}
			log.Printf("users count: %d\n", len(users))
		}
	}()

	port := config.Port()
	log.Printf("server running on %s", port)
	if err := app.Listen(port); err != nil {
		log.Fatal(err)
	}
}

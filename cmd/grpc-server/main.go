package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"backend-hexagonal/internal/adapters/grpc"
	mongoadapter "backend-hexagonal/internal/adapters/mongo"
	"backend-hexagonal/internal/config"
	"backend-hexagonal/internal/service"
)

func main() {
	// Load environment variables from .env file
	config.LoadEnv()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI()))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(config.DBName())

	// Setup repository -> service -> server
	userRepo := mongoadapter.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(userRepo)

	// Create gRPC server
	grpcServer := grpc.NewServer(userSvc, authSvc, config.GRPCPort())

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		grpcServer.Stop()
	}()

	// Start gRPC server
	log.Printf("Starting gRPC server on port %s", config.GRPCPort())
	if err := grpcServer.Start(); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}

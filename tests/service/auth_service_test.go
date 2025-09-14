package service

import (
	"backend-hexagonal/internal/domain"
	"backend-hexagonal/internal/service"
	"context"
	"testing"
)

func TestAuthService_Register(t *testing.T) {
	repo := newMockUserRepository()
	authService := service.NewAuthService(repo)

	ctx := context.Background()
	req := &domain.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	response, err := authService.Register(ctx, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.Token == "" {
		t.Error("Expected token to be generated")
	}

	if response.User.Name != req.Name {
		t.Errorf("Expected name %s, got %s", req.Name, response.User.Name)
	}

	if response.User.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, response.User.Email)
	}

	if response.User.Password != "" {
		t.Error("Expected password to be empty in response")
	}
}

func TestAuthService_Login(t *testing.T) {
	repo := newMockUserRepository()
	authService := service.NewAuthService(repo)

	ctx := context.Background()

	// First register a user
	registerReq := &domain.RegisterRequest{
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: "password123",
	}

	_, err := authService.Register(ctx, registerReq)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Now try to login
	loginReq := &domain.AuthRequest{
		Email:    "jane@example.com",
		Password: "password123",
	}

	response, err := authService.Login(ctx, loginReq)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.Token == "" {
		t.Error("Expected token to be generated")
	}

	if response.User.Email != loginReq.Email {
		t.Errorf("Expected email %s, got %s", loginReq.Email, response.User.Email)
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	repo := newMockUserRepository()
	authService := service.NewAuthService(repo)

	ctx := context.Background()

	loginReq := &domain.AuthRequest{
		Email:    "nonexistent@example.com",
		Password: "wrongpassword",
	}

	_, err := authService.Login(ctx, loginReq)

	if err == nil {
		t.Error("Expected error for invalid credentials")
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	repo := newMockUserRepository()
	authService := service.NewAuthService(repo)

	ctx := context.Background()

	// Register and get token
	registerReq := &domain.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	response, err := authService.Register(ctx, registerReq)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Validate the token
	claims, err := authService.ValidateToken(response.Token)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if claims.Email != registerReq.Email {
		t.Errorf("Expected email %s, got %s", registerReq.Email, claims.Email)
	}
}

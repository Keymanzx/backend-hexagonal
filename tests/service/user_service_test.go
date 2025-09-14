package service

import (
	"backend-hexagonal/internal/domain"
	"backend-hexagonal/internal/service"
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock repository for testing
type mockUserRepository struct {
	users map[primitive.ObjectID]*domain.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[primitive.ObjectID]*domain.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	user.ID = primitive.NewObjectID()
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *mockUserRepository) Update(ctx context.Context, id primitive.ObjectID, user *domain.User) error {
	existing, exists := m.users[id]
	if !exists {
		return nil
	}
	existing.Name = user.Name
	existing.Email = user.Email
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	delete(m.users, id)
	return nil
}

func TestUserService_CreateUser(t *testing.T) {
	repo := newMockUserRepository()
	userService := service.NewUserService(repo)

	ctx := context.Background()
	name := "John Doe"
	email := "john@example.com"
	password := "password123"

	user, err := userService.CreateUser(ctx, name, email, password)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Name != name {
		t.Errorf("Expected name %s, got %s", name, user.Name)
	}

	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}

	if user.Password != password {
		t.Errorf("Expected password %s, got %s", password, user.Password)
	}

	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	repo := newMockUserRepository()
	userService := service.NewUserService(repo)

	ctx := context.Background()

	// Create a user first
	createdUser, _ := userService.CreateUser(ctx, "Jane Doe", "jane@example.com", "password123")

	// Get the user by ID
	user, err := userService.GetUserByID(ctx, createdUser.ID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.ID != createdUser.ID {
		t.Errorf("Expected ID %s, got %s", createdUser.ID.Hex(), user.ID.Hex())
	}

	if user.Name != "Jane Doe" {
		t.Errorf("Expected name Jane Doe, got %s", user.Name)
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	repo := newMockUserRepository()
	userService := service.NewUserService(repo)

	ctx := context.Background()

	// Create multiple users
	userService.CreateUser(ctx, "User 1", "user1@example.com", "password123")
	userService.CreateUser(ctx, "User 2", "user2@example.com", "password123")

	users, err := userService.GetAllUsers(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	repo := newMockUserRepository()
	userService := service.NewUserService(repo)

	ctx := context.Background()

	// Create a user first
	createdUser, _ := userService.CreateUser(ctx, "Old Name", "old@example.com", "password123")

	// Update the user
	updatedUser, err := userService.UpdateUser(ctx, createdUser.ID, "New Name", "new@example.com")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedUser.Name != "New Name" {
		t.Errorf("Expected name New Name, got %s", updatedUser.Name)
	}

	if updatedUser.Email != "new@example.com" {
		t.Errorf("Expected email new@example.com, got %s", updatedUser.Email)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	repo := newMockUserRepository()
	userService := service.NewUserService(repo)

	ctx := context.Background()

	// Create a user first
	createdUser, _ := userService.CreateUser(ctx, "To Delete", "delete@example.com", "password123")

	// Delete the user
	err := userService.DeleteUser(ctx, createdUser.ID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Try to get the deleted user
	deletedUser, _ := userService.GetUserByID(ctx, createdUser.ID)

	if deletedUser != nil {
		t.Error("Expected user to be deleted, but it still exists")
	}
}

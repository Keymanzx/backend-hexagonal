package service

import (
	"backend-hexagonal/internal/domain"
	"backend-hexagonal/internal/ports"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	user := &domain.User{
		Name:      name,
		Email:     email,
		Password:  password, // In production, hash this password
		CreatedAt: time.Now(),
	}

	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, id primitive.ObjectID, name, email string) (*domain.User, error) {
	user := &domain.User{
		Name:  name,
		Email: email,
	}

	err := s.userRepo.Update(ctx, id, user)
	if err != nil {
		return nil, err
	}

	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	return s.userRepo.Delete(ctx, id)
}

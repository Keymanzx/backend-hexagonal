package ports

import (
	"backend-hexagonal/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, id primitive.ObjectID, user *domain.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

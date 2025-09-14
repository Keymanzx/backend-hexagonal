package grpc

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"backend-hexagonal/internal/domain"
	"backend-hexagonal/internal/service"
)

// Simple gRPC message types (instead of generated proto)
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	User    *User  `json:"user"`
	Message string `json:"message"`
}

type GetUserRequest struct {
	ID string `json:"id"`
}

type GetUserResponse struct {
	User *User `json:"user"`
}

type ListUsersRequest struct {
	Page  int32 `json:"page"`
	Limit int32 `json:"limit"`
}

type ListUsersResponse struct {
	Users []*User `json:"users"`
	Total int32   `json:"total"`
	Page  int32   `json:"page"`
	Limit int32   `json:"limit"`
}

type UserServer struct {
	userService *service.UserService
	authService *service.AuthService
}

func NewUserServer(userService *service.UserService, authService *service.AuthService) *UserServer {
	return &UserServer{
		userService: userService,
		authService: authService,
	}
}

func (s *UserServer) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	// Validate input
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "name, email, and password are required")
	}

	// Create user using auth service (which handles password hashing)
	registerReq := &domain.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	authResponse, err := s.authService.Register(ctx, registerReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert domain user to gRPC user
	grpcUser := &User{
		ID:        authResponse.User.ID.Hex(),
		Name:      authResponse.User.Name,
		Email:     authResponse.User.Email,
		CreatedAt: authResponse.User.CreatedAt,
	}

	return &CreateUserResponse{
		User:    grpcUser,
		Message: "User created successfully",
	}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	// Validate input
	if req.ID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	// Get user from service
	domainUser, err := s.userService.GetUserByID(ctx, objectID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Convert domain user to gRPC user
	grpcUser := &User{
		ID:        domainUser.ID.Hex(),
		Name:      domainUser.Name,
		Email:     domainUser.Email,
		CreatedAt: domainUser.CreatedAt,
	}

	return &GetUserResponse{
		User: grpcUser,
	}, nil
}

func (s *UserServer) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	// Get all users from service
	domainUsers, err := s.userService.GetAllUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch users")
	}

	// Convert domain users to gRPC users
	var grpcUsers []*User
	for _, domainUser := range domainUsers {
		grpcUser := &User{
			ID:        domainUser.ID.Hex(),
			Name:      domainUser.Name,
			Email:     domainUser.Email,
			CreatedAt: domainUser.CreatedAt,
		}
		grpcUsers = append(grpcUsers, grpcUser)
	}

	return &ListUsersResponse{
		Users: grpcUsers,
		Total: int32(len(grpcUsers)),
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

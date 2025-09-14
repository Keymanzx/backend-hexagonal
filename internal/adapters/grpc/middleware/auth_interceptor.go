package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"backend-hexagonal/internal/service"
)

// AuthInterceptor provides JWT authentication for gRPC
type AuthInterceptor struct {
	authService   *service.AuthService
	publicMethods map[string]bool
}

func NewAuthInterceptor(authService *service.AuthService) *AuthInterceptor {
	// Define methods that don't require authentication
	publicMethods := map[string]bool{
		"/user.UserService/CreateUser": true, // Allow user creation without auth
	}

	return &AuthInterceptor{
		authService:   authService,
		publicMethods: publicMethods,
	}
}

// UnaryInterceptor intercepts unary gRPC calls for authentication
func (interceptor *AuthInterceptor) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Check if method requires authentication
	if interceptor.publicMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	// Extract token from metadata
	token, err := interceptor.extractToken(ctx)
	if err != nil {
		return nil, err
	}

	// Validate token
	claims, err := interceptor.authService.ValidateToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// Add user info to context
	ctx = context.WithValue(ctx, "user_id", claims.UserID)
	ctx = context.WithValue(ctx, "email", claims.Email)

	return handler(ctx, req)
}

// StreamInterceptor intercepts streaming gRPC calls for authentication
func (interceptor *AuthInterceptor) StreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	// Check if method requires authentication
	if interceptor.publicMethods[info.FullMethod] {
		return handler(srv, ss)
	}

	// Extract token from metadata
	token, err := interceptor.extractToken(ss.Context())
	if err != nil {
		return err
	}

	// Validate token
	claims, err := interceptor.authService.ValidateToken(token)
	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	// Create new context with user info
	ctx := context.WithValue(ss.Context(), "user_id", claims.UserID)
	ctx = context.WithValue(ctx, "email", claims.Email)

	// Wrap the stream with new context
	wrappedStream := &wrappedServerStream{
		ServerStream: ss,
		ctx:          ctx,
	}

	return handler(srv, wrappedStream)
}

// extractToken extracts JWT token from gRPC metadata
func (interceptor *AuthInterceptor) extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	// Check for authorization header
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization header")
	}

	authHeader := authHeaders[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", status.Error(codes.Unauthenticated, "invalid authorization header format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", status.Error(codes.Unauthenticated, "missing token")
	}

	return token, nil
}

// wrappedServerStream wraps grpc.ServerStream with custom context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

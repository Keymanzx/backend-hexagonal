package grpc

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"backend-hexagonal/internal/adapters/grpc/middleware"
	"backend-hexagonal/internal/service"
)

type Server struct {
	grpcServer  *grpc.Server
	userServer  *UserServer
	authService *service.AuthService
	port        string
}

func NewServer(userService *service.UserService, authService *service.AuthService, port string) *Server {
	// Create auth interceptor
	authInterceptor := middleware.NewAuthInterceptor(authService)

	// Create gRPC server with interceptors
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.UnaryInterceptor),
		grpc.StreamInterceptor(authInterceptor.StreamInterceptor),
	)

	// Create user server
	userServer := NewUserServer(userService, authService)

	// Enable reflection for testing with tools like grpcurl
	reflection.Register(grpcServer)

	return &Server{
		grpcServer:  grpcServer,
		userServer:  userServer,
		authService: authService,
		port:        port,
	}
}

func (s *Server) Start() error {
	// Start HTTP server for gRPC-Web and REST endpoints
	go s.startHTTPServer()

	// Start gRPC server
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", s.port, err)
	}

	log.Printf("gRPC server starting on %s", s.port)
	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop() {
	log.Println("Stopping gRPC server...")
	s.grpcServer.GracefulStop()
}

// startHTTPServer starts an HTTP server that provides REST-like endpoints for gRPC services
func (s *Server) startHTTPServer() {
	httpPort := ":8081" // Different port for HTTP

	mux := http.NewServeMux()

	// Add REST-like endpoints that call gRPC methods
	mux.HandleFunc("/grpc/users", s.handleUsers)
	mux.HandleFunc("/grpc/users/", s.handleUserByID)

	log.Printf("gRPC HTTP gateway starting on %s", httpPort)
	if err := http.ListenAndServe(httpPort, mux); err != nil {
		log.Printf("HTTP server error: %v", err)
	}
}

// handleUsers handles GET (list) and POST (create) for users
func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// List users
		req := &ListUsersRequest{Page: 1, Limit: 100}
		resp, err := s.userServer.ListUsers(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resp)

	case http.MethodPost:
		// Create user
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		resp, err := s.userServer.CreateUser(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleUserByID handles GET for specific user by ID
func (s *Server) handleUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from URL path
	userID := r.URL.Path[len("/grpc/users/"):]
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	req := &GetUserRequest{ID: userID}
	resp, err := s.userServer.GetUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

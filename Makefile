# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=server
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker parameters
DOCKER_IMAGE=backend-hexagonal
DOCKER_TAG=latest

.PHONY: all build clean test coverage deps run dev docker-build docker-run docker-stop help

# Default target
all: test build

# Build the application
build:
	@mkdir -p tmp
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/server

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	$(GOTEST) -race -v ./...

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application locally
run: build
	./$(BINARY_NAME)

# Run without building (direct go run)
run-direct:
	$(GOCMD) run ./cmd/server

# Run gRPC server
run-grpc:
	$(GOCMD) run ./cmd/grpc-server

# Build gRPC server
build-grpc:
	@mkdir -p tmp
	$(GOBUILD) -o grpc-server -v ./cmd/grpc-server

# Test gRPC HTTP gateway
test-grpc:
	$(GOCMD) run ./examples/grpc-client

# Run in development mode (with auto-reload using air if available)
dev: setup
	@echo "Starting development server..."
	@if command -v air > /dev/null; then \
		echo "Using air for auto-reload..."; \
		air; \
	else \
		echo "Air not found. Install with: make install-tools"; \
		echo "Running without auto-reload..."; \
		$(GOCMD) run ./cmd/server; \
	fi

# Format code
fmt:
	$(GOCMD) fmt ./...

# Lint code (requires golangci-lint)
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Vet code
vet:
	$(GOCMD) vet ./...

# Security check (requires gosec)
security:
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Docker commands
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Database commands
db-up:
	docker-compose up -d mongo

db-down:
	docker-compose stop mongo

# Install development tools
install-tools:
	$(GOCMD) install github.com/cosmtrek/air@latest
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOCMD) install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Check if all required tools are installed
check-tools:
	@echo "Checking required tools..."
	@command -v air >/dev/null 2>&1 || echo "❌ air not found"
	@command -v golangci-lint >/dev/null 2>&1 || echo "❌ golangci-lint not found"
	@command -v gosec >/dev/null 2>&1 || echo "❌ gosec not found"
	@command -v docker >/dev/null 2>&1 || echo "❌ docker not found"
	@command -v docker-compose >/dev/null 2>&1 || echo "❌ docker-compose not found"
	@echo "✅ Tool check complete"

# Full CI pipeline
ci: deps fmt vet lint test coverage

# Production build
prod-build: clean deps test build-linux

# Setup project directories
setup:
	@mkdir -p tmp
	@echo "Project directories created"

# Test build (useful for debugging)
build-test:
	@mkdir -p tmp
	@echo "Testing build..."
	$(GOBUILD) -o ./tmp/main -v ./cmd/server
	@echo "Build successful!"

# Generate protobuf files
proto:
	@chmod +x scripts/generate-proto.sh
	@./scripts/generate-proto.sh

# Install protobuf tools
install-proto-tools:
	@echo "Installing protobuf tools..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Protobuf tools installed. Make sure 'protoc' is installed on your system."

# Help
help:
	@echo "Available commands:"
	@echo "  setup         - Create project directories"
	@echo "  build         - Build the application"
	@echo "  build-test    - Test build (for debugging)"
	@echo "  build-linux   - Build for Linux"
	@echo "  clean         - Clean build files"
	@echo "  test          - Run tests"
	@echo "  test-race     - Run tests with race detection"
	@echo "  coverage      - Run tests with coverage report"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  run           - Build and run the application"
	@echo "  run-direct    - Run without building (go run)"
	@echo "  run-grpc      - Run gRPC server"
	@echo "  build-grpc    - Build gRPC server"
	@echo "  test-grpc     - Test gRPC HTTP gateway"
	@echo "  dev           - Run in development mode with auto-reload"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  vet           - Vet code"
	@echo "  security      - Run security checks"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  docker-stop   - Stop Docker Compose"
	@echo "  docker-logs   - View Docker logs"
	@echo "  db-up         - Start MongoDB only"
	@echo "  db-down       - Stop MongoDB"
	@echo "  install-tools - Install development tools"
	@echo "  install-proto-tools - Install protobuf tools"
	@echo "  proto         - Generate protobuf files"
	@echo "  check-tools   - Check if tools are installed"
	@echo "  ci            - Run full CI pipeline"
	@echo "  prod-build    - Production build"
	@echo "  help          - Show this help message"
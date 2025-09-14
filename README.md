# Backend Hexagonal Architecture

A Go backend service implementing hexagonal architecture with user management functionality.

## Features

### Authentication
- **User Registration**: Register new users with secure password hashing
- **User Login**: Authenticate users and receive JWT tokens
- **JWT Protection**: All user endpoints protected with JWT middleware
- **Token Validation**: HMAC (HS256) signed tokens with configurable secret

### User Management
- **Get Current User**: Fetch authenticated user's profile
- **Get User by ID**: Fetch a specific user by their ID
- **List All Users**: Retrieve all users in the system
- **Update User**: Modify user's name and email
- **Delete User**: Remove a user from the system

### Logging & Monitoring
- **HTTP Request Logging**: Logs all HTTP requests with method, path, and execution time
- **Configurable Log Levels**: Support for INFO, WARN, ERROR levels
- **Multiple Log Formats**: Simple, detailed, and JSON structured logging
- **Health Check Endpoints**: `/health` and `/ready` for monitoring

### gRPC API
- **Protocol Buffers**: Defined .proto files for type-safe communication
- **gRPC Server**: High-performance gRPC server with JWT authentication
- **HTTP Gateway**: REST-like HTTP endpoints that proxy to gRPC methods
- **Token Security**: JWT token validation via gRPC metadata

## API Endpoints

### Authentication

#### Register User
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

#### Login
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

### Protected User Endpoints
**Note: All user endpoints require JWT token in Authorization header: `Bearer <token>`**

#### Get Current User
```
GET /api/v1/users/me
Authorization: Bearer <jwt_token>
```

#### Get User by ID
```
GET /api/v1/users/{id}
Authorization: Bearer <jwt_token>
```

#### List All Users
```
GET /api/v1/users
Authorization: Bearer <jwt_token>
```

#### Update User
```
PUT /api/v1/users/{id}
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "name": "Updated Name",
  "email": "updated@example.com"
}
```

#### Delete User
```
DELETE /api/v1/users/{id}
Authorization: Bearer <jwt_token>
```

### Health & Monitoring Endpoints

#### Health Check
```
GET /health
```

#### Readiness Check
```
GET /ready
```

### gRPC Endpoints

The gRPC server runs on port 9000 and provides the following services:

#### gRPC HTTP Gateway (Port 8081)
- `POST /grpc/users` - Create user via gRPC
- `GET /grpc/users` - List users via gRPC  
- `GET /grpc/users/{id}` - Get user by ID via gRPC

#### Native gRPC (Port 9000)
- `UserService.CreateUser` - Create new user
- `UserService.GetUser` - Get user by ID
- `UserService.ListUsers` - List all users
- `UserService.UpdateUser` - Update user
- `UserService.DeleteUser` - Delete user

**gRPC Authentication:**
Include JWT token in metadata:
```
authorization: Bearer <jwt_token>
```

## Architecture

This project follows hexagonal architecture principles:

- **Domain Layer** (`internal/domain/`): Core business entities
- **Ports Layer** (`internal/ports/`): Interfaces for external dependencies
- **Service Layer** (`internal/service/`): Business logic implementation
- **Adapters Layer** (`internal/adapters/`): External integrations
  - HTTP handlers for REST API
  - MongoDB repository implementation

## Running the Application

### Using Docker Compose
```bash
docker-compose up -d
```

### Local Development
1. Start MongoDB
2. Configure environment variables in `.env`:
   ```
   PORT=8000
   MONGO_URI=mongodb://localhost:27017
   DB_NAME=appdb
   JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
   LOG_LEVEL=INFO
   DETAILED_LOGGING=false
   JSON_LOGGING=false
   GRPC_PORT=9000
   ```
3. Run the servers:
   ```bash
   make deps        # Download dependencies
   make dev         # Run HTTP server with auto-reload
   # OR
   make run         # Build and run HTTP server
   
   # For gRPC server
   make run-grpc    # Run gRPC server
   
   # Test gRPC HTTP gateway
   make test-grpc   # Test gRPC endpoints
   ```

**Servers:**
- HTTP Server: Port specified in `.env` (default: 3000)
- gRPC Server: Port 9000
- gRPC HTTP Gateway: Port 8081

## Testing

Run the tests:
```bash
go test ./tests/...
```
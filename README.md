# gRPC User Service

A comprehensive gRPC service implementation demonstrating all four RPC patterns with clean architecture and Go best practices.

## 🏗️ Architecture

This project follows clean architecture principles with clear separation of concerns:

```
.
├── cmd/                    # Application entry points
│   ├── server/            # gRPC server main
│   └── client/            # gRPC client main
├── internal/              # Private application code
│   ├── models/           # Domain models
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic layer
│   ├── server/           # Server configuration
│   └── client/           # Client implementation
├── proto/                # Protocol buffer definitions
├── bin/                  # Compiled binaries (generated)
└── Makefile             # Build automation
```

## 🚀 Features

### gRPC Patterns Implemented

- **Unary RPC**: Simple request-response (GetUser, CreateUser, UpdateUser, DeleteUser)
- **Server Streaming**: Stream multiple responses (StreamUsers with filtering)
- **Client Streaming**: Accept multiple requests (CreateUsers bulk operation)
- **Bidirectional Streaming**: Real-time chat with echo and heartbeat

### Clean Code Practices

- **Dependency Injection**: Loose coupling between layers
- **Interface Segregation**: Repository pattern with clear interfaces
- **Single Responsibility**: Each package has a focused purpose
- **Error Handling**: Proper gRPC status codes and error wrapping
- **Concurrency Safety**: Thread-safe operations with proper mutex usage
- **Context Handling**: Timeout and cancellation support

## 🛠️ Prerequisites

- Go 1.24.5 or later
- Protocol Buffers compiler (`protoc`)
- Make (optional, for using Makefile)

## 📦 Installation

1. **Clone and setup:**
   ```bash
   git clone <repository-url>
   cd grpc-user-service
   ```

2. **Install dependencies:**
   ```bash
   make deps
   # or manually:
   go mod tidy
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

3. **Generate protobuf code (if needed):**
   ```bash
   make proto
   ```

## 🏃‍♂️ Running

### Using Makefile (Recommended)

```bash
# Start the server (in one terminal)
make run-server

# Run client examples (in another terminal)
make run-client

# Build binaries
make build

# Run tests
make test
```

### Using Docker

```bash
# Build and run server
make docker-run

# Run client examples (in another terminal)
make docker-client

# Clean up Docker resources
make docker-clean
```

### Using Go directly

```bash
# Start server
go run cmd/server/main.go

# Run client
go run cmd/client/main.go
```

### Environment Configuration

Copy `.env.example` to `.env` and modify as needed:

```bash
cp .env.example .env
```

## 🧪 Testing

The service includes comprehensive examples demonstrating:

1. **Unary Operations**: User CRUD operations with validation
2. **Server Streaming**: Filtered user listing with real-time streaming
3. **Client Streaming**: Bulk user creation with error aggregation
4. **Bidirectional Streaming**: Real-time chat with echo responses and heartbeats

## 📋 API Reference

### User Management

- `GetUser(UserRequest) → UserResponse`
- `CreateUser(CreateUserRequest) → UserResponse`
- `UpdateUser(UpdateUserRequest) → UserResponse`
- `DeleteUser(UserRequest) → Empty`

### Streaming Operations

- `StreamUsers(UserFilter) → stream UserResponse`
- `CreateUsers(stream CreateUserRequest) → BulkCreateResponse`
- `Chat(stream ChatMessage) → stream ChatMessage`

## 🔧 Development Tools

### gRPC Debugging

```bash
# List available services
grpcurl -plaintext localhost:50051 list

# Describe a service
grpcurl -plaintext localhost:50051 describe user.UserService

# Call a method
grpcurl -plaintext -d '{"id": 1}' localhost:50051 user.UserService/GetUser
```

### Health Check

```bash
# Using grpc_health_probe (if installed)
grpc_health_probe -addr=localhost:50051
```

## 🏛️ Design Patterns

### Repository Pattern
- Abstract data access behind interfaces
- Easy to swap implementations (in-memory → database)
- Testable with mock implementations

### Service Layer
- Business logic separation
- gRPC-specific error handling
- Context management for timeouts/cancellation

### Clean Architecture
- Dependencies point inward
- Framework-independent business logic
- Testable and maintainable code structure

## 🔒 Production Considerations

For production deployment, consider:

- **TLS/SSL**: Replace `insecure.NewCredentials()` with proper TLS
- **Authentication**: Implement proper auth middleware
- **Database**: Replace in-memory repository with persistent storage
- **Logging**: Structured logging with correlation IDs
- **Metrics**: Add Prometheus metrics and health checks
- **Rate Limiting**: Implement request rate limiting
- **Load Balancing**: Use gRPC load balancing strategies

## 📚 Learning Resources

This implementation demonstrates:

- gRPC service patterns and best practices
- Go clean architecture principles
- Concurrent programming with goroutines
- Protocol buffer usage and code generation
- Error handling and context management
- Testing strategies for gRPC services

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
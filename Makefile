# Makefile for gRPC User Service

.PHONY: proto build run-server run-client test clean help

# Variables
PROTO_DIR = proto
SERVER_CMD = cmd/server
CLIENT_CMD = cmd/client
BINARY_DIR = bin

# Default target
help:
	@echo "Available targets:"
	@echo ""
	@echo "Development:"
	@echo "  proto         - Generate protobuf and gRPC code"
	@echo "  deps          - Install dependencies"
	@echo "  build         - Build server and client binaries"
	@echo "  run-server    - Run the gRPC server"
	@echo "  run-client    - Run the gRPC client"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean generated files and binaries"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run server with Docker Compose"
	@echo "  docker-client - Run client with Docker Compose"
	@echo "  docker-clean  - Clean Docker resources"
	@echo ""
	@echo "  help          - Show this help message"

# Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	@export PATH=$$PATH:$$(go env GOPATH)/bin && \
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       $(PROTO_DIR)/user.proto

# Build binaries
build:
	@echo "Building binaries..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_DIR)/server $(SERVER_CMD)/main.go
	@go build -o $(BINARY_DIR)/client $(CLIENT_CMD)/main.go
	@echo "Binaries built in $(BINARY_DIR)/"

# Run server
run-server:
	@echo "Starting gRPC server..."
	@go run $(SERVER_CMD)/main.go

# Run client
run-client:
	@echo "Running gRPC client examples..."
	@go run $(CLIENT_CMD)/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean generated files
clean:
	@echo "Cleaning up..."
	@rm -rf $(BINARY_DIR)
	@go clean

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t grpc-user-service .

docker-run:
	@echo "Running with Docker Compose..."
	@docker-compose up --build

docker-client:
	@echo "Running client with Docker Compose..."
	@docker-compose --profile client up --build

docker-clean:
	@echo "Cleaning Docker resources..."
	@docker-compose down
	@docker rmi grpc-user-service 2>/dev/null || true
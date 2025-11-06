# Build stage
FROM golang:1.24.5-alpine AS builder

# Install protoc and git
RUN apk add --no-cache protobuf git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Expose port
EXPOSE 50051

# Set environment variables
ENV GRPC_PORT=:50051
ENV MAX_CONCURRENT_STREAMS=1000
ENV MAX_MESSAGE_SIZE=4194304

# Run the server
CMD ["./server"]
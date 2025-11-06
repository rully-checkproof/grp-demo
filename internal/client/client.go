package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"example.com/user/internal/config"
	pb "example.com/user/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Client wraps the gRPC client connection and operations
type Client struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
	config *config.Config
}

// New creates a new gRPC client instance
func New() *Client {
	cfg := config.Load()
	
	conn, err := grpc.Dial(cfg.Client.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(cfg.Client.ConnectionTimeout),
	)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	
	return &Client{
		conn:   conn,
		client: pb.NewUserServiceClient(conn),
		config: cfg,
	}
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// RunExamples demonstrates all gRPC patterns
func (c *Client) RunExamples() error {
	defer c.Close()
	
	log.Println("🎯 Starting gRPC Client Examples")
	
	if err := c.UnaryExample(); err != nil {
		return fmt.Errorf("unary example failed: %w", err)
	}
	
	if err := c.ServerStreamingExample(); err != nil {
		return fmt.Errorf("server streaming example failed: %w", err)
	}
	
	if err := c.ClientStreamingExample(); err != nil {
		return fmt.Errorf("client streaming example failed: %w", err)
	}
	
	if err := c.BidirectionalStreamingExample(); err != nil {
		return fmt.Errorf("bidirectional streaming example failed: %w", err)
	}
	
	log.Println("✅ All examples completed successfully!")
	return nil
}

// UnaryExample demonstrates unary RPC calls
func (c *Client) UnaryExample() error {
	log.Println("=== Unary RPC Example ===")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Add metadata for authentication demo
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer token123")
	
	// Test GetUser
	res, err := c.client.GetUser(ctx, &pb.UserRequest{Id: 1})
	if err != nil {
		return fmt.Errorf("GetUser failed: %w", err)
	}
	
	log.Printf("✅ User: %s (%s) - %s", res.Name, res.Email, res.Role)
	
	// Test CreateUser
	createRes, err := c.client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
		Role:  "user",
	})
	if err != nil {
		return fmt.Errorf("CreateUser failed: %w", err)
	}
	
	log.Printf("✅ Created user: %s (ID: %d)", createRes.Name, createRes.Id)
	return nil
}

// ServerStreamingExample demonstrates server streaming RPC
func (c *Client) ServerStreamingExample() error {
	log.Println("=== Server Streaming RPC Example ===")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	stream, err := c.client.StreamUsers(ctx, &pb.UserFilter{
		Keyword: "John",
		Limit:   10,
	})
	if err != nil {
		return fmt.Errorf("StreamUsers failed: %w", err)
	}
	
	count := 0
	for {
		user, err := stream.Recv()
		if err == io.EOF {
			log.Printf("✅ Stream completed - received %d users", count)
			break
		}
		if err != nil {
			return fmt.Errorf("stream receive failed: %w", err)
		}
		
		log.Printf("📨 Streamed user: %s - %s", user.Name, user.Email)
		count++
	}
	
	return nil
}

// ClientStreamingExample demonstrates client streaming RPC
func (c *Client) ClientStreamingExample() error {
	log.Println("=== Client Streaming RPC Example ===")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	stream, err := c.client.CreateUsers(ctx)
	if err != nil {
		return fmt.Errorf("CreateUsers failed: %w", err)
	}
	
	// Send bulk users
	users := []*pb.CreateUserRequest{
		{Name: "Alice Johnson", Email: "alice@example.com", Role: "user"},
		{Name: "Charlie Brown", Email: "charlie@example.com", Role: "user"},
		{Name: "David Wilson", Email: "david@example.com", Role: "admin"},
	}
	
	for _, user := range users {
		if err := stream.Send(user); err != nil {
			return fmt.Errorf("send failed: %w", err)
		}
		log.Printf("📤 Sent user: %s", user.Email)
	}
	
	result, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("close and receive failed: %w", err)
	}
	
	log.Printf("✅ Bulk create result: %d created, %d errors", 
		result.CreatedCount, len(result.Errors))
	
	for _, errMsg := range result.Errors {
		log.Printf("❌ Error: %s", errMsg)
	}
	
	return nil
}

// BidirectionalStreamingExample demonstrates bidirectional streaming RPC
func (c *Client) BidirectionalStreamingExample() error {
	log.Println("=== Bidirectional Streaming RPC Example ===")
	
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	stream, err := c.client.Chat(ctx)
	if err != nil {
		return fmt.Errorf("Chat failed: %w", err)
	}
	
	var wg sync.WaitGroup
	wg.Add(2)
	
	// Message sending goroutine
	go func() {
		defer wg.Done()
		defer stream.CloseSend()
		
		for i := 0; i < 5; i++ {
			msg := &pb.ChatMessage{
				From:      "Client",
				To:        "Server",
				Message:   fmt.Sprintf("Message %d", i+1),
				Timestamp: timestamppb.New(time.Now()),
				Type:      pb.MessageType_MESSAGE_TYPE_TEXT,
			}
			
			if err := stream.Send(msg); err != nil {
				log.Printf("Send error: %v", err)
				return
			}
			
			log.Printf("📤 Sent: %s", msg.Message)
			time.Sleep(1 * time.Second)
		}
	}()
	
	// Message receiving goroutine
	go func() {
		defer wg.Done()
		
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("Receive error: %v", err)
				return
			}
			
			log.Printf("📥 Received: %s -> %s: %s", msg.From, msg.To, msg.Message)
		}
	}()
	
	wg.Wait()
	log.Println("✅ Chat completed")
	return nil
}
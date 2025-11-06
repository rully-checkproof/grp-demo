package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"example.com/user/internal/models"
	"example.com/user/internal/repository"
	pb "example.com/user/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserService implements the gRPC UserService interface
type UserService struct {
	pb.UnimplementedUserServiceServer
	repo repository.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// GetUser implements unary RPC for user retrieval
func (s *UserService) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	log.Printf("GetUser called: ID=%d", req.Id)
	
	// Check context for timeout/cancellation
	if err := s.checkContext(ctx); err != nil {
		return nil, err
	}
	
	user, err := s.repo.GetByID(req.Id)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, status.Errorf(codes.NotFound, "User ID=%d not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "Failed to get user: %v", err)
	}
	
	return user.ToProto(), nil
}

// CreateUser implements unary RPC for user creation
func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	log.Printf("CreateUser called: email=%s", req.Email)
	
	if err := s.checkContext(ctx); err != nil {
		return nil, err
	}
	
	user := models.FromCreateRequest(req, 0) // ID will be set by repository
	
	if err := s.repo.Create(user); err != nil {
		switch err {
		case repository.ErrInvalidInput:
			return nil, status.Error(codes.InvalidArgument, "Name and email are required")
		case repository.ErrEmailExists:
			return nil, status.Errorf(codes.AlreadyExists, "Email %s already in use", req.Email)
		default:
			return nil, status.Errorf(codes.Internal, "Failed to create user: %v", err)
		}
	}
	
	return user.ToProto(), nil
}

// UpdateUser implements unary RPC for user updates
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	log.Printf("UpdateUser called: ID=%d", req.Id)
	
	if err := s.checkContext(ctx); err != nil {
		return nil, err
	}
	
	user, err := s.repo.GetByID(req.Id)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, status.Errorf(codes.NotFound, "User ID=%d not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "Failed to get user: %v", err)
	}
	
	user.Update(req)
	
	if err := s.repo.Update(user); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update user: %v", err)
	}
	
	return user.ToProto(), nil
}

// DeleteUser implements unary RPC for user deletion
func (s *UserService) DeleteUser(ctx context.Context, req *pb.UserRequest) (*emptypb.Empty, error) {
	log.Printf("DeleteUser called: ID=%d", req.Id)
	
	if err := s.checkContext(ctx); err != nil {
		return nil, err
	}
	
	if err := s.repo.Delete(req.Id); err != nil {
		if err == repository.ErrUserNotFound {
			return nil, status.Errorf(codes.NotFound, "User ID=%d not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "Failed to delete user: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

// StreamUsers implements server streaming RPC
func (s *UserService) StreamUsers(filter *pb.UserFilter, stream pb.UserService_StreamUsersServer) error {
	log.Printf("StreamUsers called: filter=%v", filter)
	
	users, err := s.repo.List(filter)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to list users: %v", err)
	}
	
	for _, user := range users {
		// Check if context is cancelled
		if stream.Context().Err() != nil {
			return stream.Context().Err()
		}
		
		if err := stream.Send(user.ToProto()); err != nil {
			return err
		}
		
		// Simulate processing delay
		time.Sleep(100 * time.Millisecond)
	}
	
	return nil
}

// CreateUsers implements client streaming RPC for bulk user creation
func (s *UserService) CreateUsers(stream pb.UserService_CreateUsersServer) error {
	log.Println("CreateUsers called - client streaming")
	
	var createdCount int32
	var userIDs []int32
	var errors []string
	
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		
		user := models.FromCreateRequest(req, 0)
		if err := s.repo.Create(user); err != nil {
			errors = append(errors, fmt.Sprintf("Email %s: %v", req.Email, err))
			continue
		}
		
		createdCount++
		userIDs = append(userIDs, user.ID)
	}
	
	return stream.SendAndClose(&pb.BulkCreateResponse{
		CreatedCount: createdCount,
		UserIds:      userIDs,
		Errors:       errors,
	})
}

// Chat implements bidirectional streaming RPC
func (s *UserService) Chat(stream pb.UserService_ChatServer) error {
	log.Println("Chat called - bidirectional streaming")
	
	var wg sync.WaitGroup
	wg.Add(2)
	
	// Message receiving goroutine
	go func() {
		defer wg.Done()
		for {
			msg, err := stream.Recv()
			if err != nil {
				return
			}
			
			log.Printf("Message received: %s -> %s: %s", msg.From, msg.To, msg.Message)
			
			// Send echo response
			response := &pb.ChatMessage{
				From:      "Server",
				To:        msg.From,
				Message:   fmt.Sprintf("Echo: %s", msg.Message),
				Timestamp: timestamppb.New(time.Now()),
				Type:      pb.MessageType_MESSAGE_TYPE_TEXT,
			}
			
			if err := stream.Send(response); err != nil {
				return
			}
		}
	}()
	
	// Heartbeat sending goroutine
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				heartbeat := &pb.ChatMessage{
					From:      "Server",
					Message:   "Heartbeat",
					Timestamp: timestamppb.New(time.Now()),
					Type:      pb.MessageType_MESSAGE_TYPE_TEXT,
				}
				if err := stream.Send(heartbeat); err != nil {
					return
				}
			case <-stream.Context().Done():
				return
			}
		}
	}()
	
	wg.Wait()
	return nil
}

// checkContext validates the request context for timeout/cancellation
func (s *UserService) checkContext(ctx context.Context) error {
	if ctx.Err() == context.DeadlineExceeded {
		return status.Error(codes.DeadlineExceeded, "Request timeout")
	}
	if ctx.Err() == context.Canceled {
		return status.Error(codes.Canceled, "Request canceled")
	}
	return nil
}
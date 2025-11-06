package models

import (
	"time"

	pb "example.com/user/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// User represents the internal user model
type User struct {
	ID        int32
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ToProto converts internal User model to protobuf UserResponse
func (u *User) ToProto() *pb.UserResponse {
	return &pb.UserResponse{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

// FromCreateRequest creates a User from CreateUserRequest
func FromCreateRequest(req *pb.CreateUserRequest, id int32) *User {
	now := time.Now()
	role := req.Role
	if role == "" {
		role = "user"
	}
	
	return &User{
		ID:        id,
		Name:      req.Name,
		Email:     req.Email,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update modifies user fields from UpdateUserRequest
func (u *User) Update(req *pb.UpdateUserRequest) {
	if req.Name != "" {
		u.Name = req.Name
	}
	if req.Email != "" {
		u.Email = req.Email
	}
	if req.Role != "" {
		u.Role = req.Role
	}
	u.UpdatedAt = time.Now()
}
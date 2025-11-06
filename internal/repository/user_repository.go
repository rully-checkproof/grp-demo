package repository

import (
	"errors"
	"sync"
	"time"

	"example.com/user/internal/models"
	pb "example.com/user/proto"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailExists     = errors.New("email already exists")
	ErrInvalidInput    = errors.New("invalid input")
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	GetByID(id int32) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int32) error
	List(filter *pb.UserFilter) ([]*models.User, error)
	EmailExists(email string) bool
}

// InMemoryUserRepository implements UserRepository using in-memory storage
type InMemoryUserRepository struct {
	users  map[int32]*models.User
	nextID int32
	mutex  sync.RWMutex
}

// NewInMemoryUserRepository creates a new in-memory user repository with sample data
func NewInMemoryUserRepository() *InMemoryUserRepository {
	now := time.Now()
	users := map[int32]*models.User{
		1: {ID: 1, Name: "John Doe", Email: "john@example.com", Role: "admin", CreatedAt: now, UpdatedAt: now},
		2: {ID: 2, Name: "Jane Smith", Email: "jane@example.com", Role: "user", CreatedAt: now, UpdatedAt: now},
		3: {ID: 3, Name: "Bob Johnson", Email: "bob@example.com", Role: "user", CreatedAt: now, UpdatedAt: now},
	}
	
	return &InMemoryUserRepository{
		users:  users,
		nextID: 4,
	}
}

func (r *InMemoryUserRepository) GetByID(id int32) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	user, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	
	// Return a copy to prevent external modifications
	userCopy := *user
	return &userCopy, nil
}

func (r *InMemoryUserRepository) Create(user *models.User) error {
	if user.Name == "" || user.Email == "" {
		return ErrInvalidInput
	}
	
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check for duplicate email
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email {
			return ErrEmailExists
		}
	}
	
	user.ID = r.nextID
	r.nextID++
	r.users[user.ID] = user
	
	return nil
}

func (r *InMemoryUserRepository) Update(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.users[user.ID]; !exists {
		return ErrUserNotFound
	}
	
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) Delete(id int32) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.users[id]; !exists {
		return ErrUserNotFound
	}
	
	delete(r.users, id)
	return nil
}

func (r *InMemoryUserRepository) List(filter *pb.UserFilter) ([]*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var result []*models.User
	count := 0
	
	for _, user := range r.users {
		// Apply keyword filter
		if filter.Keyword != "" && !contains(user.Name, filter.Keyword) {
			continue
		}
		
		// Apply role filter
		if len(filter.Roles) > 0 {
			roleMatch := false
			for _, role := range filter.Roles {
				if user.Role == role {
					roleMatch = true
					break
				}
			}
			if !roleMatch {
				continue
			}
		}
		
		// Apply limit
		if filter.Limit > 0 && count >= int(filter.Limit) {
			break
		}
		
		// Create a copy to prevent external modifications
		userCopy := *user
		result = append(result, &userCopy)
		count++
	}
	
	return result, nil
}

func (r *InMemoryUserRepository) EmailExists(email string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for _, user := range r.users {
		if user.Email == email {
			return true
		}
	}
	return false
}

// contains is a simple string search helper
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    (len(s) > len(substr) && 
		     (s[:len(substr)] == substr || 
		      s[len(s)-len(substr):] == substr ||
		      findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
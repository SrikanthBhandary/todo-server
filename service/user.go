package service

import (
	"context"
	"fmt"

	"github.com/srikanthbhandary/todo-server/entity"
	"github.com/srikanthbhandary/todo-server/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)
	GetUserByUserName(ctx context.Context, UserName string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, userID int) error
	CheckPasswordHash(password, hash string) bool
}

// UserServiceImpl is the implementation of UserService interface
type UserServiceImpl struct {
	repo repository.UserRepository // Assuming you have a repository interface for users
}

// NewUserService creates a new instance of UserServiceImpl
func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

// CreateUser creates a new user with hashed password
func (s *UserServiceImpl) CreateUser(ctx context.Context, user *entity.User) error {
	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	user.Password = string(hashedPassword) // Store the hashed password
	fmt.Println("Stored password", user.Password)
	return s.repo.CreateUser(ctx, user) // Call the repository to add the user
}

// GetUserByID retrieves a user by their ID
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	return s.repo.GetUserByID(ctx, userID) // Call the repository to get the user
}

// GetUserByUserName retrieves a user by their username
func (s *UserServiceImpl) GetUserByUserName(ctx context.Context, username string) (*entity.User, error) {
	return s.repo.GetUserByUserName(ctx, username) // Call the repository to get the user
}

// UpdateUser updates user information
func (s *UserServiceImpl) UpdateUser(ctx context.Context, user *entity.User) error {
	// Optionally, hash the password if it's provided
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword) // Store the hashed password
	}

	return s.repo.UpdateUser(ctx, user) // Call the repository to update the user
}

// DeleteUser deletes a user by their ID
func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID int) error {
	return s.repo.DeleteUser(ctx, userID) // Call the repository to delete the user
}

// checkPasswordHash compares a plain password with its hashed version
func (s *UserServiceImpl) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

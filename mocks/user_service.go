package mocks

import (
	"context"

	"github.com/srikanthbhandary/todo-server/entity"
	"github.com/stretchr/testify/mock"
)

// Mock UserService for testing
type MockUserService struct {
	mock.Mock
}

// Mock method for CreateUser
func (m *MockUserService) CreateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Mock method for GetUserByID
func (m *MockUserService) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*entity.User), args.Error(1)
}

// Mock method for GetUserByUserName
func (m *MockUserService) GetUserByUserName(ctx context.Context, userName string) (*entity.User, error) {
	args := m.Called(ctx, userName)
	return args.Get(0).(*entity.User), args.Error(1)
}

// Mock method for UpdateUser
func (m *MockUserService) UpdateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Mock method for DeleteUser
func (m *MockUserService) DeleteUser(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Mock method for CheckPasswordHash
func (m *MockUserService) CheckPasswordHash(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

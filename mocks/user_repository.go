package mocks

import (
	"context"

	"github.com/srikanthdoc/todo-server/entity"
	"github.com/stretchr/testify/mock"
)

// Mock implementation of the UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUserName(ctx context.Context, userName string) (*entity.User, error) {
	args := m.Called(ctx, userName)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

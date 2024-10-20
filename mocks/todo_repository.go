package mocks

import (
	"context"

	"github.com/srikanthbhandary/todo-server/entity"
	"github.com/stretchr/testify/mock"
)

// Mock implementation of the ToDoRepository
type MockToDoRepository struct {
	mock.Mock
}

func (m *MockToDoRepository) AddToDo(ctx context.Context, todo *entity.ToDo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func (m *MockToDoRepository) GetAllTodos(ctx context.Context, userID int) ([]entity.ToDo, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]entity.ToDo), args.Error(1)
}

func (m *MockToDoRepository) GetTodo(ctx context.Context, userID, todoID int) (entity.ToDo, error) {
	args := m.Called(ctx, userID, todoID)
	return args.Get(0).(entity.ToDo), args.Error(1)
}

func (m *MockToDoRepository) DeleteToDo(ctx context.Context, userID, todoID int) error {
	args := m.Called(ctx, userID, todoID)
	return args.Error(0)
}

func (m *MockToDoRepository) DeleteAllTodos(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

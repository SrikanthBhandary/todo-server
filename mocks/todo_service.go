package mocks

import (
	"context"

	"github.com/srikanthdoc/todo-server/entity"
	"github.com/stretchr/testify/mock"
)

// Mock ToDoService for testing
type MockToDoService struct {
	mock.Mock
}

func (m *MockToDoService) AddToDo(ctx context.Context, todo *entity.ToDo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func (m *MockToDoService) GetAllTodos(ctx context.Context, userID int) ([]entity.ToDo, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]entity.ToDo), args.Error(1)
}

func (m *MockToDoService) GetTodo(ctx context.Context, userID, todoID int) (entity.ToDo, error) {
	args := m.Called(ctx, userID, todoID)
	return args.Get(0).(entity.ToDo), args.Error(1)
}

func (m *MockToDoService) DeleteToDo(ctx context.Context, userID, todoID int) error {
	args := m.Called(ctx, userID, todoID)
	return args.Error(0)
}

func (m *MockToDoService) DeleteAllTodos(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

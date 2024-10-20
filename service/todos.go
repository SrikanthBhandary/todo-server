package service

import (
	"context"

	"time"

	"github.com/srikanthdoc/todo-server/entity"
	"github.com/srikanthdoc/todo-server/repository"
)

type ToDoService interface {
	AddToDo(ctx context.Context, todo *entity.ToDo) error
	GetAllTodos(ctx context.Context, userID int) ([]entity.ToDo, error)
	GetTodo(ctx context.Context, userID, todoID int) (entity.ToDo, error)
	DeleteToDo(ctx context.Context, userID, todoID int) error
	DeleteAllTodos(ctx context.Context, userID int) error
}

// TodoServiceImpl is the implementation of ToDoService interface
type TodoServiceImpl struct {
	repo repository.ToDoRepository
}

// NewTodoService creates a new instance of TodoServiceImpl
func NewTodoService(repo repository.ToDoRepository) *TodoServiceImpl {
	return &TodoServiceImpl{repo: repo}
}

// AddToDo adds a new todo for the specified user
func (s *TodoServiceImpl) AddToDo(ctx context.Context, todo *entity.ToDo) error {
	// You may want to set the datetime here if not set
	if todo.DateTime.IsZero() {
		todo.DateTime = time.Now()
	}
	return s.repo.AddToDo(ctx, todo) // Call the repository to add the todo
}

// GetAllTodos retrieves all todos for a specific user
func (s *TodoServiceImpl) GetAllTodos(ctx context.Context, userID int) ([]entity.ToDo, error) {
	return s.repo.GetAllTodos(ctx, userID) // Call the repository to get all todos for the user
}

// GetTodo retrieves a specific todo for a user
func (s *TodoServiceImpl) GetTodo(ctx context.Context, userID, todoID int) (entity.ToDo, error) {
	return s.repo.GetTodo(ctx, userID, todoID) // Call the repository to get the specific todo
}

// DeleteToDo deletes a specific todo for a user
func (s *TodoServiceImpl) DeleteToDo(ctx context.Context, userID, todoID int) error {
	return s.repo.DeleteToDo(ctx, userID, todoID) // Call the repository to delete the todo
}

// DeleteAllTodos deletes all todos for a specific user.
func (s *TodoServiceImpl) DeleteAllTodos(ctx context.Context, userID int) error {
	return s.repo.DeleteAllTodos(ctx, userID) // Call the repository to delete all todos for the user
}

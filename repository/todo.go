package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/srikanthdoc/todo-server/entity"
)

// ToDoRepository defines the interface for ToDo operations
type ToDoRepository interface {
	AddToDo(ctx context.Context, todo *entity.ToDo) error
	GetAllTodos(ctx context.Context, userID int) ([]entity.ToDo, error)
	GetTodo(ctx context.Context, userID, todoID int) (entity.ToDo, error)
	DeleteToDo(ctx context.Context, userID, todoID int) error
	DeleteAllTodos(ctx context.Context, userID int) error
}

// PostgresToDoRepository implements the ToDoRepository interface using PostgreSQL
type PostgresToDoRepository struct {
	DB *sql.DB
}

// NewPostgresToDoRepository creates a new PostgresToDoRepository
func NewPostgresToDoRepository(db *sql.DB) *PostgresToDoRepository {
	return &PostgresToDoRepository{DB: db}
}

// AddToDo inserts a new todo into the database
func (r *PostgresToDoRepository) AddToDo(ctx context.Context, todo *entity.ToDo) error {
	_, err := r.DB.ExecContext(ctx,
		"INSERT INTO todos (title, datetime, description, user_id) VALUES ($1, $2, $3, $4)",
		todo.Title, todo.DateTime, todo.Description, todo.UserID,
	)
	return err
}

// GetAllTodos retrieves all todos for a specific user from the database
func (r *PostgresToDoRepository) GetAllTodos(ctx context.Context, userID int) ([]entity.ToDo, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT todo_id, title, datetime, description, user_id FROM todos WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []entity.ToDo
	for rows.Next() {
		var todo entity.ToDo
		if err := rows.Scan(&todo.ToDoID, &todo.Title, &todo.DateTime, &todo.Description, &todo.UserID); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// GetTodo retrieves a specific todo for a user from the database
func (r *PostgresToDoRepository) GetTodo(ctx context.Context, userID, todoID int) (entity.ToDo, error) {
	var todo entity.ToDo
	err := r.DB.QueryRowContext(ctx, "SELECT todo_id, title, datetime, description, user_id FROM todos WHERE todo_id = $1 AND user_id = $2", todoID, userID).
		Scan(&todo.ToDoID, &todo.Title, &todo.DateTime, &todo.Description, &todo.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.ToDo{}, fmt.Errorf("todo not found")
		}
		return entity.ToDo{}, err
	}
	return todo, nil
}

// DeleteToDo deletes a specific todo for a user from the database
func (r *PostgresToDoRepository) DeleteToDo(ctx context.Context, userID, todoID int) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM todos WHERE todo_id = $1 AND user_id = $2", todoID, userID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAllTodos deletes all todos for a specific user from the database
func (r *PostgresToDoRepository) DeleteAllTodos(ctx context.Context, userID int) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM todos WHERE user_id = $1", userID)
	return err
}

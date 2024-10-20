package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/srikanthdoc/todo-server/entity"
)

// UserRepository defines the interface for user operations

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)
	GetUserByUserName(ctx context.Context, UserName string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, userID int) error
}

// PostgresUserRepository implements the UserRepository interface using PostgreSQL
type PostgresUserRepository struct {
	DB *sql.DB
}

// NewPostgresUserRepository creates a new PostgresUserRepository
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

// CreateUser inserts a new user into the database
func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	_, err := r.DB.ExecContext(ctx,
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3)",
		user.UserName, user.Email, user.Password,
	)
	return err
}

// GetUserByID retrieves a user by their ID
func (r *PostgresUserRepository) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	var user entity.User
	err := r.DB.QueryRowContext(ctx, "SELECT user_id, username, email, password FROM users WHERE user_id = $1", userID).
		Scan(&user.UserID, &user.UserName, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return &entity.User{}, fmt.Errorf("user not found")
		}
		return &entity.User{}, err
	}
	return &user, nil
}

// GetUserByUserName retrieves a user by their UserName
func (r *PostgresUserRepository) GetUserByUserName(ctx context.Context, UserName string) (*entity.User, error) {
	var user entity.User
	err := r.DB.QueryRowContext(ctx, "SELECT user_id, username, email, password FROM users WHERE username = $1", UserName).
		Scan(&user.UserID, &user.UserName, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return &entity.User{}, fmt.Errorf("user not found")
		}
		return &entity.User{}, err
	}

	return &user, nil
}

// UpdateUser updates the user's information in the database
func (r *PostgresUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	_, err := r.DB.ExecContext(ctx,
		"UPDATE users SET username = $1, email = $2, password = $3 WHERE id = $4",
		user.UserName, user.Email, user.Password, user.UserID,
	)
	return err
}

// DeleteUser deletes a user from the database
func (r *PostgresUserRepository) DeleteUser(ctx context.Context, userID int) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
	return err
}

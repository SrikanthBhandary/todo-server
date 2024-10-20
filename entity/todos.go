package entity

import "time"

// ToDo represents a todo task linked to a specific user
type ToDo struct {
	ToDoID      int       `json:"id"`
	Title       string    `json:"title"`
	DateTime    time.Time `json:"datetime"`
	Description string    `json:"description"`
	UserID      int       `json:"user_id"`
}

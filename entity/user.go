package entity

// User represents a user in the system
type User struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Email    string `json:"email"` //ignoring storing password
}

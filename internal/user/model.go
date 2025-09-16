package user

import "time"

type User struct {
	ID        string
	Email     string
	Password  string
	CreatedAt time.Time
}

type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserResponse struct {
	ID        string
	Email     string
	CreatedAt time.Time
}

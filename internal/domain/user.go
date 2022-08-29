package domain

import (
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user with such credentials not found")

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}

type SignUpInput struct {
	Name     string `json:"name" binding:"required,alpha,gte=2"`
	Surname  string `json:"surname" binding:"required,alpha,gte=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6"`
}

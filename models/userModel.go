package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}

type CreateUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type PublicUser struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

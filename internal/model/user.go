package model

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Username     string         `gorm:"type:varchar(50);unique;not null" json:"username"`
	Email        string         `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	IsAdmin      bool           `gorm:"not null;default:false" json:"is_admin"`
	CreatedAt    time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (User) TableName() string {
	return "go_user"
}

type UserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type UserResponse struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (ur *UserRequest) ValidateUserRequest() error {
	// Add custom validation logic here if needed
	if len(ur.Username) < 3 || len(ur.Username) > 50 {
		return fmt.Errorf("username must be between 3 and 50 characters")
	}
	if len(ur.Password) < 6 || len(ur.Password) > 100 {
		return fmt.Errorf("password must be between 6 and 100 characters")
	}
	// check email format
	if !strings.Contains(ur.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

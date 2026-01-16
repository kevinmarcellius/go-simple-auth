package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/kevinmarcellius/go-simple-auth/internal/model"
	"github.com/kevinmarcellius/go-simple-auth/internal/repository"
	"github.com/kevinmarcellius/go-simple-auth/internal/utils"
)


type UserService interface {
	CreateUser(ctx context.Context, req model.UserRequest) (model.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(ctx context.Context, req model.UserRequest) (model.UserResponse, error) {
	log.Println("Create new user")

	// create password hash here in real application
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return model.UserResponse{
			Message: "Password error",
		}, err
	}
	req.Password = hashedPassword
	

	newUser := model.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		}
	
	err = s.userRepo.CreateUser(newUser)
	if err != nil {
		return model.UserResponse{
			Message: "Failed to create user",
		}, err
	}

	return model.UserResponse{
		Message: "User created successfully",
	}, nil
}


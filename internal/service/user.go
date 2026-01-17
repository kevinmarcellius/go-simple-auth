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
	Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error)
	Refresh(ctx context.Context, req model.RefreshTokenRequest) (model.RefreshTokenResponse, error)
	UpdatePassword(ctx context.Context, userID string, req model.UpdatePasswordRequest) error
}

type userService struct {
	userRepo repository.UserRepository
	jwtKey   string
}

func NewUserService(userRepo repository.UserRepository, jwtKey string) UserService {
	return &userService{userRepo: userRepo, jwtKey: jwtKey}
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

func (s *userService) Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return model.LoginResponse{}, err
	}
	log.Println("User found:", user.Email)

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		log.Println("Invalid password for user:", user.Email)
		return model.LoginResponse{}, utils.ErrInvalidPassword
	}
	
	log.Println("Password verified for user:", user.Email)

	accessToken, refreshToken, err := utils.GenerateJWT(user, s.jwtKey)
	if err != nil {
		return model.LoginResponse{}, err
	}

	return model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) Refresh(ctx context.Context, req model.RefreshTokenRequest) (model.RefreshTokenResponse, error) {
	claims, err := utils.ValidateRefreshToken(req.RefreshToken, s.jwtKey)
	if err != nil {
		return model.RefreshTokenResponse{}, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return model.RefreshTokenResponse{}, err
	}

	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return model.RefreshTokenResponse{}, err
	}

	accessToken, err := utils.GenerateNewAccessToken(user, s.jwtKey)
	if err != nil {
		return model.RefreshTokenResponse{}, err
	}

	return model.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *userService) UpdatePassword(ctx context.Context, userID string, req model.UpdatePasswordRequest) error {
	// Implement password update logic here
	// get user by email
	uuid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	user, err := s.userRepo.FindUserByID(uuid)
	if err != nil {
		return err
	}

	// check old password
	if !utils.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		return utils.ErrInvalidPassword
	}
	// hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	// update user password
	user.PasswordHash = hashedPassword
	err = s.userRepo.UpdateUserById(user.ID, user)
	if err != nil {
		return err
	}
	return nil
}

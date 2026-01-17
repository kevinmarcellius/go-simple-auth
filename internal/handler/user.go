package handler

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/kevinmarcellius/go-simple-auth/internal/model"
	"github.com/kevinmarcellius/go-simple-auth/internal/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()
	// process payload and call userService.CreateUser
	var req model.UserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	err := req.ValidateUserRequest()
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	res, err := h.userService.CreateUser(ctx, req)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create user"})
	}

	return c.JSON(201, res)

}

func (h *UserHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	res, err := h.userService.Login(ctx, req)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "Invalid credentials"})
	}

	return c.JSON(200, res)
}

func (h *UserHandler) Refresh(c echo.Context) error {
	ctx := c.Request().Context()
	var req model.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	res, err := h.userService.Refresh(ctx, req)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "Invalid refresh token"})
	}

	return c.JSON(200, res)
}

func (h *UserHandler) UpdatePassword(c echo.Context) error {
	ctx := c.Request().Context()
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var req model.UpdatePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	// Implement password update logic here
	err := h.userService.UpdatePassword(ctx, userID, req)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	return c.JSON(200, map[string]string{"message": "Password updated successfully"})
}

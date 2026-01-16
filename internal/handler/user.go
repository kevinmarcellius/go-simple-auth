package handler

import (
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

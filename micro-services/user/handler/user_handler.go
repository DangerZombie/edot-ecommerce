package handler

import (
	"errors"
	"monorepo-ecommerce/micro-services/user/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) RegisterUser(c echo.Context) error {
	var req UserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	user, err := h.UserService.RegisterUser(req.Email, req.Phone, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) LoginUser(c echo.Context) error {
	var req UserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	user, err := h.UserService.LoginUser(req.Email, req.Phone, req.Password)
	if err != nil {
		if !errors.Is(err, errors.New("user not found")) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	token, err := service.GenerateToken(user.Id, user.Email, user.Phone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func RegisterUserRoutes(e *echo.Echo, userService service.UserService) {
	handler := NewUserHandler(userService)
	e.POST("/user/register", handler.RegisterUser)
	e.POST("/user/login", handler.LoginUser)
}

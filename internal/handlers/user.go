package handlers

import (
	"blogo/internal/db/sqlc"
	"blogo/internal/services/response"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	db *sqlc.Queries
}

type UserHandler interface {
	CreateUser(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	GetUser(c echo.Context) error
	GetAllUsers(c echo.Context) error
	Login(c echo.Context) error
}

func NewUserHandler(db *sqlc.Queries) UserHandler {
	return &userHandler{
		db: db,
	}
}

func (h *userHandler) CreateUser(c echo.Context) error {
	return response.Success(c, "User created successfully", nil)
}

func (h *userHandler) UpdateUser(c echo.Context) error {
	return response.Success(c, "User updated successfully", nil)
}

func (h *userHandler) DeleteUser(c echo.Context) error {
	return response.Success(c, "User deleted successfully", nil)
}

func (h *userHandler) GetUser(c echo.Context) error {
	return response.Success(c, "User fetched successfully", nil)
}

func (h *userHandler) GetAllUsers(c echo.Context) error {
	return response.Success(c, "All users fetched successfully", nil)
}

func (h *userHandler) Login(c echo.Context) error {
	return response.Success(c, "User logged in successfully", nil)
}
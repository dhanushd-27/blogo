package handlers

import "github.com/labstack/echo/v4"

type userHandler struct{}

type UserHandler interface {
	CreateUser(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	GetUser(c echo.Context) error
	GetAllUsers(c echo.Context) error
}

func NewUserHandler() UserHandler {
	return &userHandler{}
}

func (h *userHandler) CreateUser(c echo.Context) error {
	return nil
}

func (h *userHandler) UpdateUser(c echo.Context) error {
	return nil
}

func (h *userHandler) DeleteUser(c echo.Context) error {
	return nil
}

func (h *userHandler) GetUser(c echo.Context) error {
	return nil
}

func (h *userHandler) GetAllUsers(c echo.Context) error {
	return nil
}

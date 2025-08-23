package routes

import (
	"blogo/internal/handlers"

	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Echo, h handlers.UserHandler) {
	e.POST("/users", h.CreateUser)
	e.PUT("/users/:id", h.UpdateUser)
	e.DELETE("/users/:id", h.DeleteUser)
	e.GET("/users/:id", h.GetUser)
	e.GET("/users", h.GetAllUsers)
}
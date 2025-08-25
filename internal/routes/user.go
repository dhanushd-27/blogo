package routes

import (
	"blogo/internal/handlers"

	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Echo, h handlers.UserHandler) {
	e.POST("/signup", h.CreateUser)
	e.POST("/login", h.Login)
	e.GET("/users/:id", h.GetUser)
	e.GET("/users", h.GetAllUsers)
	e.PUT("/users/:id", h.UpdateUser)
	e.DELETE("/users/:id", h.DeleteUser)
}

package routes

import (
	"blogo/internal/config"
	"blogo/internal/handlers"
	"blogo/internal/middleware"

	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Echo, h handlers.UserHandler, cfg *config.Config) {
	e.POST("/signup", h.CreateUser)
	e.POST("/login", h.Login)

	userGroup := e.Group("/users")
	userGroup.Use(middleware.JWTCookieMiddleware(cfg.JWTSecret))

	userGroup.GET("/me", h.GetUser)
	userGroup.GET("", h.GetAllUsers)
	userGroup.PUT("/:id", h.UpdateUser)
	userGroup.DELETE("/:id", h.DeleteUser)
}

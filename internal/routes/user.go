package routes

import (
	"blogo/internal/config"
	handlers "blogo/internal/handlers/user"
	"blogo/internal/middleware"

	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Echo, h handlers.UserHandler, cfg *config.Config) {
	e.POST("/signup", h.CreateUser)
	e.POST("/login", h.Login)

	e.GET("/me", h.Me, middleware.JWTCookieMiddleware(cfg.JWTSecret))

	userGroup := e.Group("/user")
	userGroup.Use(middleware.JWTCookieMiddleware(cfg.JWTSecret))

	// userGroup.GET("/:id", h.GetUserById)
	userGroup.GET("/all", h.GetAllUsers)
	userGroup.PUT("/:id", h.UpdateUser)
	userGroup.DELETE("/:id", h.DeleteUser)
}

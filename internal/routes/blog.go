package routes

import (
	"blogo/internal/config"
	"blogo/internal/handlers"
	"blogo/internal/middleware"

	"github.com/labstack/echo/v4"
)

func BlogRoutes(e *echo.Echo, h handlers.BlogHandler, cfg *config.Config) {
	// Public route for getting all blogs (no auth required)
	e.GET("/blogs", h.GetAllBlogs)
	e.GET("/blogs/:id", h.GetBlog)

	// Protected routes that require authentication
	blogGroup := e.Group("/blogs")
	blogGroup.Use(middleware.JWTCookieMiddleware(cfg.JWTSecret))

	blogGroup.POST("", h.CreateBlog)
	blogGroup.PUT("/:id", h.UpdateBlog)
	blogGroup.DELETE("/:id", h.DeleteBlog)
}

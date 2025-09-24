package routes

import (
	"blogo/internal/config"
	blog "blogo/internal/handlers/blog"
	"blogo/internal/middleware"

	"github.com/labstack/echo/v4"
)

func BlogRoutes(e *echo.Echo, h blog.BlogHandler, cfg *config.Config) {
	// Public route for getting all blogs (no auth required)
	e.GET("/blogs", h.GetAllBlogs)
	e.GET("/blogs/:id", h.GetBlogByID)

	// Protected routes that require authentication
	blogGroup := e.Group("/blogs")
	blogGroup.Use(middleware.JWTCookieMiddleware(cfg.JWTSecret))

	blogGroup.POST("", h.CreateBlog)
	blogGroup.PATCH("/:id", h.UpdateBlog)
	blogGroup.DELETE("/:id", h.DeleteBlog)
}

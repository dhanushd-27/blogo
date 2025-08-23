package routes

import (
	"blogo/internal/handlers"

	"github.com/labstack/echo/v4"
)

func BlogRoutes(e *echo.Echo, h handlers.BlogHandler) {
	e.POST("/blogs", h.CreateBlog)
	e.PUT("/blogs/:id", h.UpdateBlog)
	e.DELETE("/blogs/:id", h.DeleteBlog)
	e.GET("/blogs/:id", h.GetBlog)
	e.GET("/blogs", h.GetAllBlogs)
}

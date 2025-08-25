package handlers

import (
	"blogo/internal/db/sqlc"
	"blogo/internal/services/response"

	"github.com/labstack/echo/v4"
)

type blogHandler struct {
	db *sqlc.Queries
}

type BlogHandler interface {
	CreateBlog(c echo.Context) error
	UpdateBlog(c echo.Context) error
	DeleteBlog(c echo.Context) error
	GetBlog(c echo.Context) error
	GetAllBlogs(c echo.Context) error
}

func NewBlogHandler(db *sqlc.Queries) BlogHandler {
	return &blogHandler{
		db: db,
	}
}

func (h *blogHandler) CreateBlog(c echo.Context) error {
	return response.Success(c, "Blog created successfully", nil)
}

func (h *blogHandler) UpdateBlog(c echo.Context) error {
	return response.Success(c, "Blog updated successfully", nil)
}

func (h *blogHandler) DeleteBlog(c echo.Context) error {
	return response.Success(c, "Blog deleted successfully", nil)
}

func (h *blogHandler) GetBlog(c echo.Context) error {
	return response.Success(c, "Blog fetched successfully", nil)
}

func (h *blogHandler) GetAllBlogs(c echo.Context) error {
	return response.Success(c, "All blogs fetched successfully", nil)
}

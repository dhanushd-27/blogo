package handlers

import "github.com/labstack/echo/v4"

type blogHandler struct{}

type BlogHandler interface {
	CreateBlog(c echo.Context) error
	UpdateBlog(c echo.Context) error
	DeleteBlog(c echo.Context) error
	GetBlog(c echo.Context) error
	GetAllBlogs(c echo.Context) error
}

func NewBlogHandler() BlogHandler {
	return &blogHandler{}
}

func (h *blogHandler) CreateBlog(c echo.Context) error {
	return nil
}

func (h *blogHandler) UpdateBlog(c echo.Context) error {
	return nil
}

func (h *blogHandler) DeleteBlog(c echo.Context) error {
	return nil
}

func (h *blogHandler) GetBlog(c echo.Context) error {
	return nil
}

func (h *blogHandler) GetAllBlogs(c echo.Context) error {
	return nil
}
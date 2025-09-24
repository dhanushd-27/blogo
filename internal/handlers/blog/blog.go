package blog

import (
	"blogo/internal/db/sqlc"
	"blogo/internal/services/model"
	"blogo/internal/services/response"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type blogHandler struct {
	queries sqlc.Querier
}

type BlogHandler interface {
	CreateBlog(c echo.Context) error
	UpdateBlog(c echo.Context) error
	DeleteBlog(c echo.Context) error
	GetBlogByID(c echo.Context) error
	GetAllBlogs(c echo.Context) error
}

func NewBlogHandler(queries sqlc.Querier) BlogHandler {
	return &blogHandler{
		queries: queries,
	}
}

func (h *blogHandler) CreateBlog(c echo.Context) error {
	blog := model.CreateBlog{}
	if err := c.Bind(&blog); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if err := c.Validate(&blog); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation faild: " + err.Error()})
	}

	userID := c.Get("user_id").(float64)

	blogCreated, err := h.queries.CreateBlog(context.Background(), sqlc.CreateBlogParams{
		Title:   blog.Title,
		Content: blog.Content,
		UserID:  int32(userID),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create blog"})
	}

	return response.Success(c, "Blog created successfully", blogCreated)
}

func (h *blogHandler) UpdateBlog(c echo.Context) error {
	return response.Success(c, "Blog updated successfully", nil)
}

func (h *blogHandler) DeleteBlog(c echo.Context) error {
	return response.Success(c, "Blog deleted successfully", nil)
}

func (h *blogHandler) GetBlogByID(c echo.Context) error {
	return response.Success(c, "Blog fetched successfully", nil)
}

func (h *blogHandler) GetAllBlogs(c echo.Context) error {
	return response.Success(c, "All blogs fetched successfully", nil)
}

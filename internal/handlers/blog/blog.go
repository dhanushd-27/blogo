package blog

import (
	"blogo/internal/db/sqlc"
	"blogo/internal/services/model"
	"blogo/internal/services/response"
	"context"
	"fmt"
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
	// Parse blog ID from path
	idParam := c.Param("id")
	if idParam == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid blog id"})
	}

	// Bind and validate the request body
	req := model.UpdateBlog{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed: " + err.Error()})
	}

	// Require at least one field for PATCH
	if req.Title == nil && req.Content == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "At least one field (title or content) must be provided"})
	}

	// Convert id to int32
	var blogID int32
	if _, err := fmt.Sscanf(idParam, "%d", &blogID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid blog id"})
	}

	// Ensure blog exists; fetch current values
	existing, err := h.queries.GetBlogByID(context.Background(), blogID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Blog not found"})
	}

	// Prepare update params by merging provided fields; don't touch id or user_id
	updateParams := sqlc.UpdateBlogParams{
		ID:      existing.ID,
		Title:   existing.Title,
		Content: existing.Content,
	}
	if req.Title != nil {
		updateParams.Title = *req.Title
	}
	if req.Content != nil {
		updateParams.Content = *req.Content
	}

	updatedBlog, err := h.queries.UpdateBlog(context.Background(), updateParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update blog"})
	}

	return response.Success(c, "Blog updated successfully", updatedBlog)
}

func (h *blogHandler) DeleteBlog(c echo.Context) error {
	return response.Success(c, "Blog deleted successfully", nil)
}

func (h *blogHandler) GetBlogByID(c echo.Context) error {
	return response.Success(c, "Blog fetched successfully", nil)
}

func (h *blogHandler) GetAllBlogs(c echo.Context) error {
	blogs, err := h.queries.ListBlogs(context.Background(), sqlc.ListBlogsParams{
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch blogs"})
	}

	return response.Success(c, "Blogs fetched successfully", blogs)
}

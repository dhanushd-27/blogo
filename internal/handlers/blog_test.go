package handlers

import (
	"blogo/internal/config"
	"blogo/internal/db/mocks"
	"blogo/internal/db/sqlc"
	"blogo/internal/handlers/blog"
	appmw "blogo/internal/middleware"
	"blogo/internal/services/helper"
	"blogo/internal/services/model"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllBlogs(t *testing.T) {
	e := echo.New()
	e.Validator = model.NewValidator()

	cfg := &config.Config{
		JWTSecret: "test",
	}

	t.Run("success - get all blogs", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := blog.NewBlogHandler(mockDB)

		req := httptest.NewRequest(http.MethodGet, "/blogs", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/blogs/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.SetPath("/blogs/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@gmail.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.GetAllBlogs)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "All blogs fetched successfully")
	})

	t.Run("unauthorized - missing jwt cookie", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := blog.NewBlogHandler(mockDB)

		req := httptest.NewRequest(http.MethodGet, "/blogs", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.GetAllBlogs)
		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("unauthorized - expired jwt", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := blog.NewBlogHandler(mockDB)

		req := httptest.NewRequest(http.MethodGet, "/blogs", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@gmail.com", -1*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.GetAllBlogs)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("unauthorized - invalid jwt", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := blog.NewBlogHandler(mockDB)

		req := httptest.NewRequest(http.MethodGet, "/blogs", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		jwtStr, err := helper.MakeTestJWT("wrong-secret", 1, "test", "test@gmail.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.GetAllBlogs)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestGetBlogByID(t *testing.T) {
	e := echo.New()
	e.Validator = model.NewValidator()

	cfg := &config.Config{JWTSecret: "test"}

	fmt.Println(cfg)

	t.Run("success - get blog by id", func(t *testing.T) {
		assert.Equal(t, "test pending", "test pending")
	})

	t.Run("unauthorized - missing jwt cookie", func(t *testing.T) {
		assert.Equal(t, "test pending", "test pending")
	})

	t.Run("unauthorized - expired jwt", func(t *testing.T) {
		assert.Equal(t, "test pending", "test pending")
	})

	t.Run("unauthorized - invalid jwt", func(t *testing.T) {
		assert.Equal(t, "test pending", "test pending")
	})

	t.Run("not found - blog not found", func(t *testing.T) {
		assert.Equal(t, "test pending", "test pending")
	})
}

func TestCreateBlog(t *testing.T) {
	assert.Equal(t, "test pending", "test pending")
}

func TestUpdateBlog(t *testing.T) {
	e := echo.New()
	e.Validator = model.NewValidator()

	cfg := &config.Config{
		JWTSecret: "test",
	}

	t.Run("success - update blog", func(t *testing.T) {
		ctx := context.Background()
		mockDB := mocks.NewMockQuerier(t)

		h := blog.NewBlogHandler(mockDB)
		requestBody := `{"title":"Updated Title","content":"Updated Content"}`
		req := httptest.NewRequest(http.MethodPatch, "/blogs/1", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/blogs/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@gmail.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mockDB.On("GetBlogByID", ctx, int32(1)).Return(sqlc.Blog{
			ID:        1,
			Title:     "Test Blog",
			Content:   "Test Content",
			UserID:    1,
			CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
			UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		}, nil)

		mockDB.On("UpdateBlog", ctx, sqlc.UpdateBlogParams{
			ID:      1,
			Title:   "Updated Title",
			Content: "Updated Content",
		}).Return(sqlc.Blog{
			ID:        1,
			Title:     "Updated Title",
			Content:   "Updated Content",
			UserID:    1,
			CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
			UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		}, nil)

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.UpdateBlog)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Blog updated successfully")
	})

	t.Run("unauthorized - missing jwt cookie", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)

		h := blog.NewBlogHandler(mockDB)

		req := httptest.NewRequest(http.MethodPatch, "/blogs/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/blogs/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.UpdateBlog)
		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("unauthorized - expired jwt", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)

		h := blog.NewBlogHandler(mockDB)

		req := httptest.NewRequest(http.MethodPatch, "/blogs/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/blogs/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@gmail.com", -1*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.UpdateBlog)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("unauthorized - invalid jwt", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)

		h := blog.NewBlogHandler(mockDB)

		req := httptest.NewRequest(http.MethodPatch, "/blogs/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/blogs/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		jwtStr, err := helper.MakeTestJWT("wrong-secret", 1, "test", "test@gmail.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.UpdateBlog)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("not found - blog not found", func(t *testing.T) {
		ctx := context.Background()
		mockDB := mocks.NewMockQuerier(t)

		h := blog.NewBlogHandler(mockDB)

		reqBody := strings.NewReader(`{"title":"Updated Title","content":"Updated content"}`)
		req := httptest.NewRequest(http.MethodPatch, "/blogs/1", reqBody)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/blogs/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@gmail.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mockDB.On("GetBlogByID", ctx, int32(1)).Return(sqlc.Blog{}, echo.ErrNotFound)

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.UpdateBlog)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "Blog not found")
	})

	t.Run("bad request - invalid body", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)

		h := blog.NewBlogHandler(mockDB)

		// malformed JSON to trigger bind error
		req := httptest.NewRequest(http.MethodPatch, "/blogs/1", strings.NewReader("{"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/blogs/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@gmail.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.UpdateBlog)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid request body")
	})
}

func TestDeleteBlog(t *testing.T) {
	assert.Equal(t, "test pending", "test pending")
}

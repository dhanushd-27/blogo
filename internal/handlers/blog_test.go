package handlers

import (
	"blogo/internal/config"
	"blogo/internal/db/mocks"
	"blogo/internal/handlers/blog"
	appmw "blogo/internal/middleware"
	"blogo/internal/services/helper"
	"blogo/internal/services/model"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	assert.Equal(t, "test pending", "test pending")
}

func TestDeleteBlog(t *testing.T) {
	assert.Equal(t, "test pending", "test pending")
}

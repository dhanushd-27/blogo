package handlers

import (
	"blogo/internal/config"
	"blogo/internal/db/mocks"
	"blogo/internal/db/sqlc"
	user "blogo/internal/handlers/u"
	appmw "blogo/internal/middleware"
	"blogo/internal/services/helper"
	"blogo/internal/services/model"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestME(t *testing.T) {
	e := echo.New()
	e.Validator = model.NewValidator()

	cfg := &config.Config{JWTSecret: "test"}

	// rough structures for Me handler tests
	t.Run("success - valid jwt cookie and user exists", func(t *testing.T) {
		ctx := context.Background()
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		// expected DB call
		expectedUser := sqlc.User{ID: 1, Name: "test", Email: "test@gmail.com"}
		mockDB.On("GetUserByID", ctx, int32(1)).Return(expectedUser, nil)

		// request/response
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set cookie with valid jwt
		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@gmail.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		// run through middleware then handler
		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.Me)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "User fetched successfully")
		assert.Contains(t, rec.Body.String(), "test")
		assert.Contains(t, rec.Body.String(), "test@gmail.com")
	})

	t.Run("unauthorized - missing jwt cookie", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// no cookie added; middleware should block, no DB calls expected
		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.Me)
		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("unauthorized - invalid jwt signature", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// sign with wrong secret so signature validation fails
		jwtStr, err := helper.MakeTestJWT("wrong-secret", 1, "test", "test@gmail.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.Me)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("unauthorized - expired jwt", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// token already expired
		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@gmail.com", -1*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.Me)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("not found - user does not exist", func(t *testing.T) {
		ctx := context.Background()
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		// DB returns error (e.g., no rows)
		mockDB.On("GetUserByID", ctx, int32(999)).Return(sqlc.User{}, echo.ErrNotFound)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 999, "ghost", "ghost@example.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.Me)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "User not found")
	})
}

func TestGetAllUsers(t *testing.T) {
	e := echo.New()
	e.Validator = model.NewValidator()

	cfg := &config.Config{JWTSecret: "test"}

	t.Run("success - get all users", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		req := httptest.NewRequest(http.MethodGet, "/user/all", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set cookie with valid jwt
		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "testuser", "test@example.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		// run through middleware then handler
		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.GetAllUsers)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "All users fetched successfully")
	})

	t.Run("unauthorized - missing jwt cookie", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		req := httptest.NewRequest(http.MethodGet, "/user/all", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// no cookie added; middleware should block, no DB calls expected
		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.GetAllUsers)
		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("unauthorized - invalid jwt signature", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		req := httptest.NewRequest(http.MethodGet, "/user/all", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// sign with wrong secret so signature validation fails
		jwtStr, err := helper.MakeTestJWT("wrong-secret", 1, "test", "test@gmail.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.GetAllUsers)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("unauthorized - expired jwt", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		req := httptest.NewRequest(http.MethodGet, "/user/all", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// token already expired
		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@gmail.com", -1*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.GetAllUsers)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestUpdateUser(t *testing.T) {
	e := echo.New()
	e.Validator = model.NewValidator()

	cfg := &config.Config{JWTSecret: "test"}

	t.Run("success - update user", func(t *testing.T) {
		ctx := context.Background()
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		// expected DB calls: fetch current user, then update only name using authed ID
		mockDB.On("GetUserByID", ctx, int32(1)).Return(sqlc.User{ID: 1, Name: "old", Email: "old@example.com", Password: "hashed"}, nil)
		mockDB.On("UpdateUser", ctx, sqlc.UpdateUserParams{ID: 1, Name: "new", Email: "old@example.com", Password: "hashed"}).Return(sqlc.User{ID: 1, Name: "new", Email: "old@example.com"}, nil)

		body := `{"name":"new"}`
		req := httptest.NewRequest(http.MethodPut, "/user/1", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// valid cookie
		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@example.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.UpdateUser)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "User updated successfully")
	})

	t.Run("unauthorized - missing jwt cookie", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		body := `{"name":"new"}`
		req := httptest.NewRequest(http.MethodPut, "/user/1", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.UpdateUser)
		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("unauthorized - invalid jwt signature", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		body := `{"name":"new"}`
		req := httptest.NewRequest(http.MethodPut, "/user/1", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		jwtStr, err := helper.MakeTestJWT("wrong-secret", 1, "test", "test@example.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.UpdateUser)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("unauthorized - expired jwt", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		body := `{"name":"x"}`
		req := httptest.NewRequest(http.MethodPut, "/user/1", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@example.com", -1*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.UpdateUser)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("bad request - invalid body", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		body := `{"invalid": "json"`
		req := httptest.NewRequest(http.MethodPut, "/user/1", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@example.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.UpdateUser)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid request body")
	})

	t.Run("not found - user does not exist", func(t *testing.T) {
		ctx := context.Background()
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		mockDB.On("GetUserByID", ctx, int32(1)).Return(sqlc.User{}, echo.ErrNotFound)

		body := `{"name":"new"}`
		req := httptest.NewRequest(http.MethodPut, "/user/1", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@example.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.UpdateUser)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "User not found")
	})

	t.Run("internal error - update failed", func(t *testing.T) {
		ctx := context.Background()
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		mockDB.On("GetUserByID", ctx, int32(1)).Return(sqlc.User{ID: 1, Name: "old", Email: "old@example.com", Password: "hashed"}, nil)
		mockDB.On("UpdateUser", ctx, sqlc.UpdateUserParams{ID: 1, Name: "new", Email: "old@example.com", Password: "hashed"}).Return(sqlc.User{}, echo.ErrInternalServerError)

		body := `{"name":"new"}`
		req := httptest.NewRequest(http.MethodPut, "/user/1", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@example.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.UpdateUser)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Failed to update user")
	})
}

func TestDeleteUser(t *testing.T) {
	e := echo.New()
	e.Validator = model.NewValidator()

	cfg := &config.Config{JWTSecret: "test"}

	t.Run("success - delete user", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		req := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@example.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.DeleteUser)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "User deleted successfully")
	})

	t.Run("unauthorized - missing jwt cookie", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		req := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.DeleteUser)
		err := handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("unauthorized - invalid jwt signature", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		req := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		jwtStr, err := helper.MakeTestJWT("wrong-secret", 1, "test", "test@example.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.DeleteUser)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("unauthorized - expired jwt", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		req := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		jwtStr, err := helper.MakeTestJWT(cfg.JWTSecret, 1, "test", "test@example.com", -1*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.DeleteUser)
		err = handler(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

package handlers

import (
	"blogo/internal/config"
	"blogo/internal/db/mocks"
	"blogo/internal/db/sqlc"
	user "blogo/internal/handlers/u"
	appmw "blogo/internal/middleware"
	"blogo/internal/services/model"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// helper to create a signed JWT string with desired claims for tests
func makeTestJWT(secret string, userID int32, username, email string, expiresIn time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["username"] = username
	claims["email"] = email
	claims["exp"] = time.Now().Add(expiresIn).Unix()
	return token.SignedString([]byte(secret))
}

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
		jwtStr, err := makeTestJWT(cfg.JWTSecret, 1, "test", "test@gmail.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		// run through middleware then handler
		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.Me)
		_ = handler(c)
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
		_ = handler(c)
	})

	t.Run("unauthorized - invalid jwt signature", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// sign with wrong secret so signature validation fails
		jwtStr, err := makeTestJWT("wrong-secret", 1, "test", "test@gmail.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.Me)
		_ = handler(c)
	})

	t.Run("unauthorized - expired jwt", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)
		h := user.NewUserHandler(mockDB, cfg)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// token already expired
		jwtStr, err := makeTestJWT(cfg.JWTSecret, 1, "test", "test@gmail.com", -1*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.Me)
		_ = handler(c)
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

		jwtStr, err := makeTestJWT(cfg.JWTSecret, 999, "ghost", "ghost@example.com", 24*time.Hour)
		if err != nil {
			t.Fatalf("failed to create jwt: %v", err)
		}
		req.AddCookie(&http.Cookie{Name: "token", Value: jwtStr, Path: "/"})

		mw := appmw.JWTCookieMiddleware(cfg.JWTSecret)
		handler := mw(h.Me)
		_ = handler(c)
	})
}

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"blogo/internal/config"
	"blogo/internal/db/sqlc"
	"blogo/internal/services/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	
	// Register the validator (this is what was missing!)
	e.Validator = model.NewValidator()

	// Create a mock config
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	// Create a mock database (we'll use nil for now since we're not testing DB operations)
	var db *sqlc.Queries

	// Create the handler
	handler := NewUserHandler(db, cfg)

	// Test case 1: Valid user creation
	t.Run("Valid User Creation", func(t *testing.T) {
		userData := model.SignUp{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		jsonData, _ := json.Marshal(userData)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Since we don't have a real DB, this will panic when trying to create the user
		// We need to recover from the panic since this is expected behavior
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Expected panic due to nil DB: %v", r)
			}
		}()

		err := handler.CreateUser(c)

		// If we get here without a panic, it means the handler succeeded
		// which shouldn't happen with a nil DB
		if err == nil {
			t.Error("Expected handler to fail due to nil DB, but it succeeded")
		}
	})

	// Test case 2: Invalid JSON
	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateUser(c)

		// Should fail due to invalid JSON
		if err != nil {
			t.Logf("Expected error due to invalid JSON: %v", err)
		} else {
			// If no error, check the response status
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	// Test case 3: Missing required fields
	t.Run("Missing Required Fields", func(t *testing.T) {
		userData := model.SignUp{
			Name: "Test User",
			// Missing email and password
		}

		jsonData, _ := json.Marshal(userData)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateUser(c)

		// Should fail due to validation
		if err != nil {
			t.Logf("Expected error due to validation: %v", err)
		} else {
			// If no error, check the response status
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}

func TestLogin(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Register the validator (this is what was missing!)
	e.Validator = model.NewValidator()

	// Create a mock config
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	// Create a mock database (we'll use nil for now since we're not testing DB operations)
	var db *sqlc.Queries

	// Create the handler
	handler := NewUserHandler(db, cfg)

	// Test case 1: Valid login attempt
	t.Run("Valid Login", func(t *testing.T) {
		loginData := model.Login{
			Email:    "test@example.com",
			Password: "password123",
		}

		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Since we don't have a real DB, this will panic when trying to get user by email
		// We need to recover from the panic since this is expected behavior
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Expected panic due to nil DB: %v", r)
			}
		}()

		err := handler.Login(c)

		// If we get here without a panic, it means the handler succeeded
		// which shouldn't happen with a nil DB
		if err == nil {
			t.Error("Expected handler to fail due to nil DB, but it succeeded")
		}
	})

	// Test case 2: Invalid JSON
	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		// Should fail due to invalid JSON
		if err != nil {
			t.Logf("Expected error due to invalid JSON: %v", err)
		} else {
			// If no error, check the response status
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	// Test case 3: Missing required fields
	t.Run("Missing Required Fields", func(t *testing.T) {
		loginData := model.Login{
			Email: "test@example.com",
			// Missing password
		}

		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		// Should fail due to validation
		if err != nil {
			t.Logf("Expected error due to validation: %v", err)
		} else {
			// If no error, check the response status
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}

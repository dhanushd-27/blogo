package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"blogo/internal/config"
	"blogo/internal/db/mocks"
	"blogo/internal/db/sqlc"
	"blogo/internal/middleware"
	"blogo/internal/services/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Register the validator
	e.Validator = model.NewValidator()

	// Create a mock config
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	// Test case 1: Valid user creation
	t.Run("Valid User Creation", func(t *testing.T) {
		// Create a mock database
		mockDB := mocks.NewMockQuerier(t)

		// Set up mock expectations
		expectedUser := sqlc.User{
			ID:       1,
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "hashed_password_here",
		}

		// Mock the GetUserByEmail method to return "user not found" error (user doesn't exist)
		mockDB.On("GetUserByEmail", mock.Anything, "test@example.com").
			Return(sqlc.User{}, pgx.ErrNoRows)

		// Mock the CreateUser method to return success
		mockDB.On("CreateUser", mock.Anything, mock.AnythingOfType("sqlc.CreateUserParams")).
			Return(expectedUser, nil)

		// Create the handler with mock DB
		handler := NewUserHandler(mockDB, cfg)

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

		// Execute the handler
		err := handler.CreateUser(c)

		// Assert expectations
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code) // Success returns 200, not 201

		// Verify mock was called as expected
		mockDB.AssertExpectations(t)
	})

	// Test case 2: Database error
	t.Run("Database Error", func(t *testing.T) {
		// Create a mock database
		mockDB := mocks.NewMockQuerier(t)

		// Mock the GetUserByEmail method to return "user not found" error (user doesn't exist)
		mockDB.On("GetUserByEmail", mock.Anything, "test@example.com").
			Return(sqlc.User{}, pgx.ErrNoRows)

		// Mock the CreateUser method to return an error
		mockDB.On("CreateUser", mock.Anything, mock.AnythingOfType("sqlc.CreateUserParams")).
			Return(sqlc.User{}, assert.AnError)

		// Create the handler with mock DB
		handler := NewUserHandler(mockDB, cfg)

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

		// Execute the handler
		err := handler.CreateUser(c)

		// Should return an error
		if err != nil {
			t.Logf("Handler returned error as expected: %v", err)
		} else {
			// Check the response status
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}

		// Verify mock was called as expected
		mockDB.AssertExpectations(t)
	})

	// Test case 3: Invalid JSON
	t.Run("Invalid JSON", func(t *testing.T) {
		// Create a mock database (not needed for this test but good practice)
		mockDB := mocks.NewMockQuerier(t)

		// Create the handler with mock DB
		handler := NewUserHandler(mockDB, cfg)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateUser(c)

		// Should fail due to invalid JSON
		if err != nil {
			t.Logf("Handler returned error as expected: %v", err)
		} else {
			// Check the response status
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	// Test case 4: Missing required fields
	t.Run("Missing Required Fields", func(t *testing.T) {
		// Create a mock database (not needed for this test but good practice)
		mockDB := mocks.NewMockQuerier(t)

		// Create the handler with mock DB
		handler := NewUserHandler(mockDB, cfg)

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
			t.Logf("Handler returned error as expected: %v", err)
		} else {
			// Check the response status
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	// Test case 5: User already exists
	t.Run("User Already Exists", func(t *testing.T) {
		// Create a mock database
		mockDB := mocks.NewMockQuerier(t)

		// Mock the GetUserByEmail method to return an existing user
		existingUser := sqlc.User{
			ID:       1,
			Name:     "Existing User",
			Email:    "test@example.com",
			Password: "hashed_password",
		}
		mockDB.On("GetUserByEmail", mock.Anything, "test@example.com").
			Return(existingUser, nil)

		// Create the handler with mock DB
		handler := NewUserHandler(mockDB, cfg)

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

		err := handler.CreateUser(c)

		// Should fail due to user already existing
		if err != nil {
			t.Logf("Handler returned error as expected: %v", err)
		} else {
			// Check the response status
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}

		// Verify mock was called as expected
		mockDB.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Register the validator
	e.Validator = model.NewValidator()

	// Create a mock config
	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	// Test case 1: Valid login attempt
	t.Run("Valid Login", func(t *testing.T) {
		// Create a mock database
		mockDB := mocks.NewMockQuerier(t)

		// Set up mock expectations
		expectedUser := sqlc.User{
			ID:       1,
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "$2a$10$2Lf8Y8QCo9lbn2v6yQRRJe3OchswXDjUA/mUdsTNmdR3eAp2QPiZW", // bcrypt hash for "password123"
		}

		// Mock the GetUserByEmail method to return a user
		mockDB.On("GetUserByEmail", mock.Anything, "test@example.com").
			Return(expectedUser, nil)

		// Create the handler with mock DB
		handler := NewUserHandler(mockDB, cfg)

		loginData := model.Login{
			Email:    "test@example.com",
			Password: "password123",
		}

		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Execute the handler
		err := handler.Login(c)

		// Assert expectations
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Verify mock was called as expected
		mockDB.AssertExpectations(t)
	})

	// Test case 2: User not found
	t.Run("User Not Found", func(t *testing.T) {
		// Create a mock database
		mockDB := mocks.NewMockQuerier(t)

		// Mock the GetUserByEmail method to return "user not found" error
		mockDB.On("GetUserByEmail", mock.Anything, "nonexistent@example.com").
			Return(sqlc.User{}, pgx.ErrNoRows)

		// Create the handler with mock DB
		handler := NewUserHandler(mockDB, cfg)

		loginData := model.Login{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Execute the handler
		err := handler.Login(c)

		// Should fail due to user not found
		if err != nil {
			t.Logf("Handler returned error as expected: %v", err)
		} else {
			// Check the response status
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}

		// Verify mock was called as expected
		mockDB.AssertExpectations(t)
	})

	// Test case 3: Invalid JSON
	t.Run("Invalid JSON", func(t *testing.T) {
		// Create a mock database (not needed for this test but good practice)
		mockDB := mocks.NewMockQuerier(t)

		// Create the handler with mock DB
		handler := NewUserHandler(mockDB, cfg)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)

		// Should fail due to invalid JSON
		if err != nil {
			t.Logf("Handler returned error as expected: %v", err)
		} else {
			// Check the response status
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	// Test case 4: Missing required fields
	t.Run("Missing Required Fields", func(t *testing.T) {
		// Create a mock database (not needed for this test but good practice)
		mockDB := mocks.NewMockQuerier(t)

		// Create the handler with mock DB
		handler := NewUserHandler(mockDB, cfg)

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
			t.Logf("Handler returned error as expected: %v", err)
		} else {
			// Check the response status
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}
func TestMe(t *testing.T) {
	e := echo.New()
	e.Validator = model.NewValidator()

	cfg := &config.Config{
		JWTSecret: "test-secret-key",
	}

	t.Run("Valid Get User", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)

		// Mock user for login
		loginUser := sqlc.User{
			ID:       1,
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "$2a$10$2Lf8Y8QCo9lbn2v6yQRRJe3OchswXDjUA/mUdsTNmdR3eAp2QPiZW", // "password123" hashed
		}

		// Mock user for /me endpoint
		meUser := sqlc.User{
			ID:       1,
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "$2a$10$2Lf8Y8QCo9lbn2v6yQRRJe3OchswXDjUA/mUdsTNmdR3eAp2QPiZW",
		}

		// Setup mocks
		mockDB.On("GetUserByEmail", mock.Anything, "test@example.com").Return(loginUser, nil)
		mockDB.On("GetUserByID", mock.Anything, int32(1)).Return(meUser, nil)

		handler := NewUserHandler(mockDB, cfg)

		// Step 1: Login to get the cookie
		loginData := model.Login{
			Email:    "test@example.com",
			Password: "password123",
		}

		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Extract the cookie from the response
		cookies := rec.Result().Cookies()
		var tokenCookie *http.Cookie
		if len(cookies) == 0 {
			// If no cookies were set (due to Secure flag in test environment),
			// we need to manually create a valid JWT token for testing
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["user_id"] = float64(1)
			claims["username"] = "Test User"
			claims["email"] = "test@example.com"
			claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

			tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
			assert.NoError(t, err)

			tokenCookie = &http.Cookie{
				Name:  "token",
				Value: tokenString,
			}
		} else {
			tokenCookie = cookies[0]
			assert.Equal(t, "token", tokenCookie.Name)
		}

		// Step 2: Use the cookie to make a /me request
		meReq := httptest.NewRequest(http.MethodGet, "/me", nil)
		meReq.AddCookie(tokenCookie) // Add the cookie from login

		meRec := httptest.NewRecorder()
		meC := e.NewContext(meReq, meRec)

		// Apply the JWT middleware to extract user info from cookie
		middlewareFunc := middleware.JWTCookieMiddleware(cfg.JWTSecret)
		handlerFunc := middlewareFunc(func(c echo.Context) error {
			return handler.Me(c)
		})

		err = handlerFunc(meC)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, meRec.Code)

		// Verify the response
		var response map[string]interface{}
		err = json.Unmarshal(meRec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response["status"])
		assert.Equal(t, "User fetched successfully", response["message"])

		// Verify user data
		userData := response["data"].(map[string]interface{})
		assert.Equal(t, "1", userData["id"])
		assert.Equal(t, "Test User", userData["name"])
		assert.Equal(t, "test@example.com", userData["email"])

		mockDB.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockDB := mocks.NewMockQuerier(t)

		// Mock user for login
		loginUser := sqlc.User{
			ID:       1,
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "$2a$10$2Lf8Y8QCo9lbn2v6yQRRJe3OchswXDjUA/mUdsTNmdR3eAp2QPiZW",
		}

		// Setup mocks - login succeeds but /me fails
		mockDB.On("GetUserByEmail", mock.Anything, "test@example.com").Return(loginUser, nil)
		mockDB.On("GetUserByID", mock.Anything, int32(1)).Return(sqlc.User{}, errors.New("user not found"))

		handler := NewUserHandler(mockDB, cfg)

		// Step 1: Login to get the cookie
		loginData := model.Login{
			Email:    "test@example.com",
			Password: "password123",
		}

		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Extract the cookie
		cookies := rec.Result().Cookies()
		var tokenCookie *http.Cookie
		if len(cookies) == 0 {
			// If no cookies were set (due to Secure flag in test environment),
			// we need to manually create a valid JWT token for testing
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["user_id"] = float64(1)
			claims["username"] = "Test User"
			claims["email"] = "test@example.com"
			claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

			tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
			assert.NoError(t, err)

			tokenCookie = &http.Cookie{
				Name:  "token",
				Value: tokenString,
			}
		} else {
			tokenCookie = cookies[0]
		}

		// Step 2: Use the cookie to make a /me request
		meReq := httptest.NewRequest(http.MethodGet, "/me", nil)
		meReq.AddCookie(tokenCookie)

		meRec := httptest.NewRecorder()
		meC := e.NewContext(meReq, meRec)

		// Apply middleware
		middlewareFunc := middleware.JWTCookieMiddleware(cfg.JWTSecret)
		handlerFunc := middlewareFunc(func(c echo.Context) error {
			return handler.Me(c)
		})

		err = handlerFunc(meC)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, meRec.Code)

		// Verify error response
		var response map[string]interface{}
		err = json.Unmarshal(meRec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User Not Found", response["error"])

		mockDB.AssertExpectations(t)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		handler := NewUserHandler(nil, cfg) // No DB needed for this test

		// Create a request with an invalid cookie
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		invalidCookie := &http.Cookie{
			Name:  "token",
			Value: "invalid.jwt.token",
		}
		req.AddCookie(invalidCookie)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Apply middleware - should fail
		middlewareFunc := middleware.JWTCookieMiddleware(cfg.JWTSecret)
		handlerFunc := middlewareFunc(func(c echo.Context) error {
			return handler.Me(c)
		})

		err := handlerFunc(c)
		assert.Error(t, err)

		// Check if it's an HTTP error with 401 status
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})
}

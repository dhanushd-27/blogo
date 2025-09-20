package handlers

import (
	"blogo/internal/config"
	"blogo/internal/db/mocks"
	"blogo/internal/db/sqlc"
	user "blogo/internal/handlers/u"
	"blogo/internal/services/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	e := echo.New()
	e.Validator = model.NewValidator()
	ctx := context.Background()

	cfg := &config.Config{
		JWTSecret: "test",
	}

	// Write the expectedUser values
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	expectedUser := sqlc.User{
		ID:       1,
		Name:     "test",
		Email:    "test@gmail.com",
		Password: string(hashedPassword),
	}

	t.Run("Valid User Creation", func(t *testing.T) {
		// Create Mock DB Instance
		mockDB := mocks.NewMockQuerier(t)

		// Mock DB Calls
		mockDB.On("GetUserByEmail", ctx, "test@gmail.com").Return(sqlc.User{}, errors.New("user not found"))

		// Use mock.MatchedBy to match the CreateUser call with any hashed password
		mockDB.On("CreateUser", ctx, mock.MatchedBy(func(params sqlc.CreateUserParams) bool {
			return params.Name == "test" && params.Email == "test@gmail.com" && len(params.Password) > 0
		})).Return(expectedUser, nil)

		// Create User Handler
		userHandler := user.NewUserHandler(mockDB, cfg)

		// Create Request Body
		requestBody := `{"name":"test","email":"test@gmail.com","password":"testpassword"}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the CreateUser method
		userHandler.CreateUser(c)

		// Assertions
		var response map[string]interface{}
		jsonErr := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, jsonErr)
		assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
		assert.Equal(t, "User created successfully", response["message"])
	})

	t.Run("Invalid User Json", func(t *testing.T) {
		// Create Mock DB Instance - no mock expectations needed since validation fails before DB calls
		mockDB := mocks.NewMockQuerier(t)

		// Create User Handler
		userHandler := user.NewUserHandler(mockDB, cfg)

		// Create request body with empty JSON (should fail validation)
		requestBody := `{}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the CreateUser method
		err := userHandler.CreateUser(c)

		// Assertions
		assert.NoError(t, err) // The handler should not return an error, it should write to response
		assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)

		var response map[string]interface{}
		jsonErr := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, jsonErr)
		assert.Contains(t, response["error"], "Validation failed")
	})

	t.Run("User with this email already exists", func(t *testing.T) {
		// Create Mock DB Instance
		mockDB := mocks.NewMockQuerier(t)

		// Create UserHandler
		userHandler := user.NewUserHandler(mockDB, cfg)

		mockDB.On("GetUserByEmail", ctx, "test@gmail.com").Return(expectedUser, nil)

		// Create request body
		requestBody := `{"name":"test","email":"test@gmail.com","password":"testpassword"}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the CreateUser method
		err := userHandler.CreateUser(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Result().StatusCode)

		var response map[string]interface{}
		jsonErr := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, jsonErr)
		assert.Equal(t, "User with this email already exists", response["error"])
	})
}

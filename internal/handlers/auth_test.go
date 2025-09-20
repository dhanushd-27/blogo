package handlers

// import (
// 	"blogo/internal/config"
// 	"blogo/internal/db/mocks"
// 	"blogo/internal/db/sqlc"
// 	user "blogo/internal/handlers/u"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/labstack/echo/v4"
// 	"github.com/stretchr/testify/assert"
// 	"golang.org/x/crypto/bcrypt"
// )

// func TestCreateUser(t *testing.T) {
// 	e := echo.New()
// 	ctx := context.Background()

// 	cfg := &config.Config{
// 		JWTSecret: "test",
// 	}

// 	t.Run("Valid User Creation", func(t *testing.T) {
// 		// Create Mock DB Instance
// 		mockDB := mocks.NewMockQuerier(t)

// 		// Write the expectedUser values
// 		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
// 		expectedUser := sqlc.User{
// 			ID:       1,
// 			Name:     "test",
// 			Email:    "test@gmail.com",
// 			Password: string(hashedPassword),
// 		}

// 		// Mock DB Calls
// 		mockDB.On("GetUserByEmail", ctx, expectedUser.Email).Return(sqlc.User{}, errors.New("user not found"))

// 		mockDB.On("CreateUser", ctx, sqlc.CreateUserParams{
// 			Name:     expectedUser.Name,
// 			Email:    expectedUser.Email,
// 			Password: expectedUser.Password,
// 		}).Return(expectedUser, nil)

// 		// Create User Handler
// 		userHandler := user.NewUserHandler(mockDB, cfg)

// 		// Create Request Body
// 		requestBody := `{"name":"test","email":"test@gmail.com","password":"testpassword"}`
// 		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(requestBody))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)

// 		// Call the CreateUser method
// 		err := userHandler.CreateUser(c)
// 		assert.NoError(t, err)

// 		// Assertions
// 		var response map[string]interface{}
// 		jsonErr := json.Unmarshal(rec.Body.Bytes(), &response)
// 		assert.NoError(t, jsonErr)
// 		assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
// 		assert.Equal(t, "User created successfully", response["message"])
// 	})
// }

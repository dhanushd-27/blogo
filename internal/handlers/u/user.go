package u

import (
	"blogo/internal/config"
	"blogo/internal/db/sqlc"
	"blogo/internal/services/model"
	"blogo/internal/services/response"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	queries sqlc.Querier
	cfg     *config.Config
}

type UserHandlerInterface interface {
	CreateUser(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	Me(c echo.Context) error
	GetAllUsers(c echo.Context) error
	Login(c echo.Context) error
}

func NewUserHandler(queries sqlc.Querier, cfg *config.Config) UserHandlerInterface {
	return &UserHandler{
		queries: queries,
		cfg:     cfg,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	u := model.SignUp{}
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed: " + err.Error()})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process password"})
	}

	_, err = h.queries.GetUserByEmail(context.Background(), u.Email)
	if err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "User with this email already exists"})
	}

	user, err := h.queries.CreateUser(context.Background(), sqlc.CreateUserParams{
		Name:     u.Name,
		Email:    u.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return response.Success(c, "User created successfully", user)
}

func (h *UserHandler) Login(c echo.Context) error {
	u := model.Login{}

	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed: " + err.Error()})
	}

	user, err := h.queries.GetUserByEmail(context.Background(), u.Email)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	// Generate a JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["username"] = user.Name
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate authentication token"})
	}

	// setting up a http only cookie
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(time.Hour * 24)
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteNoneMode

	c.SetCookie(cookie)

	return response.Success(c, "User logged in successfully", nil)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	u := model.UpdateUser{}

	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed: " + err.Error()})
	}

	_, err := h.queries.GetUserByID(context.Background(), u.ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	_, err = h.queries.UpdateUser(context.Background(), sqlc.UpdateUserParams{
		Name:     *u.Name,
		Email:    *u.Email,
		Password: *u.Password,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
	}

	return response.Success(c, "User updated successfully", nil)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	return response.Success(c, "User deleted successfully", nil)
}

func (h *UserHandler) Me(c echo.Context) error {
	userID := c.Get("user_id").(float64)

	user, err := h.queries.GetUserByID(context.Background(), int32(userID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return response.Success(c, "User fetched successfully", map[string]string{
		"id":    strconv.Itoa(int(user.ID)),
		"name":  user.Name,
		"email": user.Email,
	})
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
	return response.Success(c, "All users fetched successfully", nil)
}

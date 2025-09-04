package handlers

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

type userHandler struct {
	db  sqlc.Querier
	cfg *config.Config
}

type UserHandler interface {
	CreateUser(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	Me(c echo.Context) error
	GetAllUsers(c echo.Context) error
	Login(c echo.Context) error
}

func NewUserHandler(db sqlc.Querier, cfg *config.Config) UserHandler {
	return &userHandler{
		db:  db,
		cfg: cfg,
	}
}

func (h *userHandler) CreateUser(c echo.Context) error {
	u := model.SignUp{}
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	_, err = h.db.GetUserByEmail(context.Background(), u.Email)
	if err == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User already exists"})
	}

	user, err := h.db.CreateUser(context.Background(), sqlc.CreateUserParams{
		Name:     u.Name,
		Email:    u.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return response.Success(c, "User created successfully", user)
}

func (h *userHandler) Login(c echo.Context) error {
	u := model.Login{}

	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := h.db.GetUserByEmail(context.Background(), u.Email)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid password"})
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
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

func (h *userHandler) UpdateUser(c echo.Context) error {
	u := model.UpdateUser{}

	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid body type"})
	}

	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid body fields"})
	}

	_, err := h.db.GetUserByID(context.Background(), u.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	_, err = h.db.UpdateUser(context.Background(), sqlc.UpdateUserParams{
		Name:     *u.Name,
		Email:    *u.Email,
		Password: *u.Password,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User Update Failed"})
	}

	return response.Success(c, "User updated successfully", nil)
}

func (h *userHandler) DeleteUser(c echo.Context) error {
	return response.Success(c, "User deleted successfully", nil)
}

func (h *userHandler) Me(c echo.Context) error {
	userID := c.Get("user_id").(float64)

	user, err := h.db.GetUserByID(context.Background(), int32(userID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User Not Found"})
	}

	return response.Success(c, "User fetched successfully", map[string]string{
		"id":    strconv.Itoa(int(user.ID)),
		"name":  user.Name,
		"email": user.Email,
	})
}

func (h *userHandler) GetAllUsers(c echo.Context) error {
	return response.Success(c, "All users fetched successfully", nil)
}

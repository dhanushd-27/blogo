package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Created(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{
		Status: "created",
		Message: message,
		Data:    data,
	})
}

func BadRequest(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusBadRequest, Response{
		Status:  "bad request",
		Message: message,
		Data:    data,
	})
}

func NotFound(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusNotFound, Response{
		Status:  "not found",
		Message: message,
		Data:    data,
	})
}
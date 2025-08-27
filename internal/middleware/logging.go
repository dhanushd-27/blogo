package middleware

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Start the timer
		start := time.Now()
		// Call the next middleware or handler
		err := next(c)
		// If there is an error, log it
		// Log the request and response
		log.Printf("--> %s %s %d %s", c.Request().Method, c.Request().URL.Path, c.Response().Status, time.Since(start).String())
		return err
	}
}

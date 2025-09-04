package middleware

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		log.Printf("--> %s %s %s", c.Request().Method, c.Request().URL.Path, time.Since(start).String())
		return err
	}
}

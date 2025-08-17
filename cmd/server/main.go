package main

import (
	"blogo/internal/config"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":" + config.Port))

	log.Printf("Server is running on port %s", config.Port)
}
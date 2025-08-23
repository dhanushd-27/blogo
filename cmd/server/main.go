package main

import (
	"blogo/internal/config"
	"blogo/internal/db"
	"blogo/internal/handlers"
	"blogo/internal/routes"

	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")
	defer db.Close()

	e := echo.New()

	// Register routes here
	routes.HealthCheck(e)
	routes.BlogRoutes(e, handlers.NewBlogHandler())
	routes.UserRoutes(e, handlers.NewUserHandler())

	e.Logger.Fatal(e.Start(":" + cfg.Port))

	log.Printf("Server is running on port %s", cfg.Port)
}

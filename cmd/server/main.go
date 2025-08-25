package main

import (
	"blogo/internal/config"
	"blogo/internal/db"
	"blogo/internal/db/sqlc"
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

	dbPool, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")
	defer dbPool.Close()

	// Create SQLC queries instance
	queries := sqlc.New(dbPool)

	e := echo.New()

	// Register routes here
	routes.HealthCheck(e)
	routes.BlogRoutes(e, handlers.NewBlogHandler(queries))
	routes.UserRoutes(e, handlers.NewUserHandler(queries))

	e.Logger.Fatal(e.Start(":" + cfg.Port))

	log.Printf("Server is running on port %s", cfg.Port)
}

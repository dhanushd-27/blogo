package main

import (
	"blogo/internal/config"
	"blogo/internal/db"
	"blogo/internal/db/sqlc"
	blog "blogo/internal/handlers/blog"
	user "blogo/internal/handlers/user"
	"blogo/internal/middleware"
	"blogo/internal/routes"
	"blogo/internal/services/model"

	"log"

	"github.com/labstack/echo/v4"
	em "github.com/labstack/echo/v4/middleware"
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

	e.Use(middleware.Logger)
	e.Use(em.CORSWithConfig(em.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	e.Validator = model.NewValidator()

	// Register routes here
	routes.HealthCheck(e)
	routes.BlogRoutes(e, blog.NewBlogHandler(queries), cfg)
	routes.UserRoutes(e, user.NewUserHandler(queries, cfg), cfg)

	e.Logger.Fatal(e.Start(":" + cfg.Port))

	log.Printf("Server is running on port %s", cfg.Port)
}

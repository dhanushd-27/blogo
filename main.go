package main

import (
	"log"
	"os"

	"github.com/dhanushd-27/blog_go/db"
	"github.com/dhanushd-27/blog_go/helper"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	database := db.NewDB(db.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	})

	server := helper.NewApiServer(":8080", database)
	err := server.Run()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

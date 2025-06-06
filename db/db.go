package db

import (
	"fmt"
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
	err error
}

type DBConfig struct {
	Host string
	Port string
	User string
	Password string
	DBName string
	SSLMode string
}

func NewDB(config DBConfig) *DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	DB := &DB{}

	DB.db, DB.err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if DB.err != nil {
		log.Fatalf("Failed to connect to database: %v", DB.err)
	}

	fmt.Println("Connected to database")

	return DB
}
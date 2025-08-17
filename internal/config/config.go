package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

const (
	envFilePath = ".env"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
}

var (
	instance *Config
	once     sync.Once
	loadErr  error
)

func Load() (*Config, error) {
	once.Do(func() {
		err := godotenv.Load(envFilePath)
		if err != nil {
			if !os.IsNotExist(err) {
				loadErr = fmt.Errorf("%s file doesn't exist", envFilePath)
			}
			loadErr = fmt.Errorf("error loading .env file %v", err)
			return
		}

		instance = &Config{
			Port:       getEnv("PORT"),
			DBHost:     getEnv("DB_HOST"),
			DBPort:     getEnv("DB_PORT"),
			DBUser:     getEnv("DB_USER"),
			DBPassword: getEnv("DB_PASSWORD"),
		}
	})

	return instance, loadErr
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	panic(fmt.Sprintf("Environment variable %s is not set", key))
}

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
	Port                  string
	DBHost                string
	DBPort                string
	DBUser                string
	DBPassword            string
	DBName                string
	SSLMode               string
	DBMaxConn             string
	DBMinConn             string
	DBConnMaxLifetime     string
	DBConnMaxIdleLifetime string
	DBHealthCheckPeriod   string
	ConnectTimeout        string
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
			Port:                  getEnv("PORT"),
			DBHost:                getEnv("DB_HOST"),
			DBPort:                getEnv("DB_PORT"),
			DBUser:                getEnv("DB_USER"),
			DBPassword:            getEnv("DB_PASSWORD"),
			DBName:                getEnv("DB_NAME"),
			SSLMode:               getEnv("SSL_MODE"),
			DBMaxConn:             getEnv("DB_MAX_CONN"),
			DBMinConn:             getEnv("DB_MIN_CONN"),
			DBConnMaxLifetime:     getEnv("DB_CONN_MAX_LIFETIME"),
			DBConnMaxIdleLifetime: getEnv("DB_CONN_MAX_IDLE_LIFETIME"),
			DBHealthCheckPeriod:   getEnv("DB_HEALTH_CHECK_PERIOD"),
			ConnectTimeout:        getEnv("CONNECT_TIMEOUT"),
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

package infrastructure

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port   int
	DBHost string
	DBPort int
	DBUser string
	DBPass string
	DBName string
}

func LoadConfig() (*Config, error) {
	port, err := strconv.Atoi(getEnv("APP_PORT", "8080"))
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
		return nil, err
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:   port,
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: dbPort,
		DBUser: getEnv("DB_USER", "postgres"),
		DBPass: getEnv("DB_PASS", ""),
		DBName: getEnv("DB_NAME", "aTES"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

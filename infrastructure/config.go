package infrastructure

import (
	"log"
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
	}
}

package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	AppPort     string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
	JWTExpireHr int
}

// LoadConfig loads settings from environment variables (and .env)
func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "sistergo"),
		JWTSecret:  getEnv("JWT_SECRET", "secret"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

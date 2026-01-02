package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret     string
	JWTExpiration string

	// Logger
	LogLevel  string
	LogFormat string

	// App
	Environment string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file (ignore error in production)
	_ = godotenv.Load()

	config := &Config{
		ServerPort: getEnv("PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "digital_wallet"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),

		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "text"),

		Environment: getEnv("ENVIRONMENT", "development"),
	}

	return config, nil
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

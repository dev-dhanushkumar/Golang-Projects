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

	// App
	Environment string

	// DATA SECRET
	DataSecret string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		ServerPort: getEnv("PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "personal-ess"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		JWTSecret:     getEnv("JWT_SECRET", "enter secret key here"),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),

		Environment: getEnv("ENVIRONMENT", "development"),

		DataSecret: getEnv("DATA_SECRET", "Enter data secret key"),
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

package database

import (
	"digital-wallet-api/internal/models"
	"digital-wallet-api/pkg/logger"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Connect establishes database connection
func Connect(config Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("Database connected successfully", map[string]interface{}{
		"host":     config.Host,
		"port":     config.Port,
		"database": config.DBName,
	})

	return nil
}

// AutoMigrate runs database migrations
func AutoMigrate() error {
	logger.Info("Running database migrations...")

	// Enable UUID extension for PostgreSQL
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		logger.Warn("Failed to enable uuid-ossp extension", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Enable pgcrypto for gen_random_uuid()
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error; err != nil {
		logger.Warn("Failed to enable pgcrypto extension", map[string]interface{}{
			"error": err.Error(),
		})
	}

	err := DB.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Category{},
		&models.Transaction{},
		&models.Transfer{},
		&models.Budget{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	logger.Info("Database migrations completed successfully")

	// Seed default categories
	if err := seedDefaultCategories(); err != nil {
		logger.Warn("Failed to seed default categories", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return nil
}

// seedDefaultCategories creates default categories if they don't exist
func seedDefaultCategories() error {
	var count int64
	DB.Model(&models.Category{}).Where("is_default = ?", true).Count(&count)

	if count == 0 {
		categories := models.DefaultCategories()
		result := DB.Create(&categories)
		if result.Error != nil {
			return result.Error
		}
		logger.Info("Default categories seeded", map[string]interface{}{
			"count": len(categories),
		})
	}

	return nil
}

// Close closes database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	logger.Info("Closing database connection...")
	return sqlDB.Close()
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

package database

import (
	"fmt"
	"personal-expense-splitting-settlement/internal/models"
	"time"

	"go.uber.org/zap"
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

// connect establishes database connection
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
		return fmt.Errorf("failed to connec to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

// AutoMigration runs datbase migration
func AutoMigration(sugar *zap.SugaredLogger) error {
	// Enable UUID extension for PostgresSQL
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		sugar.Warnln("Failed to enable uuid-ossp extension", "error", err)
	}

	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error; err != nil {
		sugar.Warnln("Failed to enable pgcrypt extension", "error", err)
	}

	err := DB.AutoMigrate(
		&models.User{},
		&models.UserSession{},
	)

	if err != nil {
		return fmt.Errorf("Failed to migrate database: %w", err)
	}

	sugar.Info("Database migration completed successfully")

	return nil
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// Logger implementation here
	return sqlDB.Close()
}

func GetDB() *gorm.DB {
	return DB
}

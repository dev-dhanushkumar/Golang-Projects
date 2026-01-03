package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	Version  string
	Name     string
	UpSQL    string
	DownSQL  string
	FilePath string
}

// RunMigrations executes all pending migrations
func RunMigrations(sugar *zap.SugaredLogger) error {
	// Create migrations table if not exists
	if err := createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all migration files
	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Apply pending migrations
	for _, migration := range migrations {
		if !contains(appliedMigrations, migration.Version) {
			sugar.Infow("Applying migration", "version", migration.Version, "name", migration.Name)
			if err := applyMigration(migration); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
			}
			sugar.Infow("Migration applied successfully", "version", migration.Version)
		}
	}

	sugar.Info("All migrations completed successfully")
	return nil
}

// createMigrationsTable creates the migrations tracking table
func createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`
	return DB.Exec(query).Error
}

// loadMigrations reads all migration files from the migrations directory
func loadMigrations() ([]Migration, error) {
	migrationsDir := "migrations"

	// Check if migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("migrations directory does not exist: %s", migrationsDir)
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	migrationsMap := make(map[string]*Migration)

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Parse filename: 000001_create_users_table.up.sql
		parts := strings.Split(file.Name(), "_")
		if len(parts) < 2 {
			continue
		}

		version := parts[0]
		isUp := strings.HasSuffix(file.Name(), ".up.sql")
		isDown := strings.HasSuffix(file.Name(), ".down.sql")

		if !isUp && !isDown {
			continue
		}

		// Extract migration name
		nameWithExt := strings.Join(parts[1:], "_")
		name := strings.TrimSuffix(strings.TrimSuffix(nameWithExt, ".up.sql"), ".down.sql")

		// Get or create migration entry
		if _, exists := migrationsMap[version]; !exists {
			migrationsMap[version] = &Migration{
				Version: version,
				Name:    name,
			}
		}

		// Read file content
		content, err := os.ReadFile(filepath.Join(migrationsDir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		if isUp {
			migrationsMap[version].UpSQL = string(content)
			migrationsMap[version].FilePath = file.Name()
		} else {
			migrationsMap[version].DownSQL = string(content)
		}
	}

	// Convert map to sorted slice
	var migrations []Migration
	for _, migration := range migrationsMap {
		if migration.UpSQL != "" { // Only include migrations with up SQL
			migrations = append(migrations, *migration)
		}
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// getAppliedMigrations returns list of applied migration versions
func getAppliedMigrations() ([]string, error) {
	var versions []string
	err := DB.Raw("SELECT version FROM schema_migrations ORDER BY version").Scan(&versions).Error
	return versions, err
}

// applyMigration applies a single migration
func applyMigration(migration Migration) error {
	// Start a transaction
	return DB.Transaction(func(tx *gorm.DB) error {
		// Execute the migration SQL
		if err := tx.Exec(migration.UpSQL).Error; err != nil {
			return fmt.Errorf("failed to execute migration SQL: %w", err)
		}

		// Record the migration
		query := "INSERT INTO schema_migrations (version, name) VALUES (?, ?)"
		if err := tx.Exec(query, migration.Version, migration.Name).Error; err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}

		return nil
	})
}

// RollbackMigration rolls back the last applied migration
func RollbackMigration(sugar *zap.SugaredLogger) error {
	// Get the last applied migration
	var lastVersion string
	err := DB.Raw("SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&lastVersion).Error
	if err != nil {
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	if lastVersion == "" {
		sugar.Info("No migrations to rollback")
		return nil
	}

	// Load all migrations
	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Find the migration to rollback
	var migrationToRollback *Migration
	for _, m := range migrations {
		if m.Version == lastVersion {
			migrationToRollback = &m
			break
		}
	}

	if migrationToRollback == nil {
		return fmt.Errorf("migration %s not found", lastVersion)
	}

	if migrationToRollback.DownSQL == "" {
		return fmt.Errorf("no down migration found for version %s", lastVersion)
	}

	sugar.Infow("Rolling back migration", "version", lastVersion, "name", migrationToRollback.Name)

	// Execute rollback in transaction
	return DB.Transaction(func(tx *gorm.DB) error {
		// Execute the down migration
		if err := tx.Exec(migrationToRollback.DownSQL).Error; err != nil {
			return fmt.Errorf("failed to execute down migration: %w", err)
		}

		// Remove migration record
		if err := tx.Exec("DELETE FROM schema_migrations WHERE version = ?", lastVersion).Error; err != nil {
			return fmt.Errorf("failed to remove migration record: %w", err)
		}

		sugar.Infow("Migration rolled back successfully", "version", lastVersion)
		return nil
	})
}

// GetMigrationStatus returns the current migration status
func GetMigrationStatus() ([]map[string]interface{}, error) {
	appliedMigrations, err := getAppliedMigrations()
	if err != nil {
		return nil, err
	}

	allMigrations, err := loadMigrations()
	if err != nil {
		return nil, err
	}

	var status []map[string]interface{}
	for _, migration := range allMigrations {
		applied := contains(appliedMigrations, migration.Version)
		status = append(status, map[string]interface{}{
			"version": migration.Version,
			"name":    migration.Name,
			"applied": applied,
		})
	}

	return status, nil
}

// Helper function to check if slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

package main

import (
	"fmt"
	"os"
	"personal-expense-splitting-settlement/internal/config"
	"personal-expense-splitting-settlement/internal/database"
	"personal-expense-splitting-settlement/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Initialize logger
	sugar, err := logger.InitLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer sugar.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		sugar.Fatalw("Failed to load config", "error", err)
	}

	// Connect to database
	dbConfig := database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	}

	if err := database.Connect(dbConfig); err != nil {
		sugar.Fatalw("Failed to connect to database", "error", err)
	}
	defer database.Close()

	command := os.Args[1]

	switch command {
	case "migrate":
		handleMigrate(sugar)
	case "rollback":
		handleRollback(sugar)
	case "status":
		handleStatus(sugar)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleMigrate(sugar *zap.SugaredLogger) {
	fmt.Println("Running migrations...")
	if err := database.RunMigrationsFromFiles(sugar); err != nil {
		sugar.Fatalw("Migration failed", "error", err)
	}
	fmt.Println("✅ All migrations completed successfully")
}

func handleRollback(sugar *zap.SugaredLogger) {
	fmt.Println("Rolling back last migration...")
	if err := database.RollbackMigration(sugar); err != nil {
		sugar.Fatalw("Rollback failed", "error", err)
	}
	fmt.Println("✅ Rollback completed successfully")
}

func handleStatus(sugar *zap.SugaredLogger) {
	status, err := database.GetMigrationStatus()
	if err != nil {
		sugar.Fatalw("Failed to get migration status", "error", err)
	}

	fmt.Println("\nMigration Status:")
	fmt.Println("================")
	for _, s := range status {
		applied := "❌ Not Applied"
		if s["applied"].(bool) {
			applied = "✅ Applied"
		}
		fmt.Printf("%s - %s [%s]\n", s["version"], s["name"], applied)
	}
	fmt.Println()
}

func printUsage() {
	fmt.Println("Migration Tool Usage:")
	fmt.Println("  go run cmd/migrate/main.go <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  migrate   - Run all pending migrations")
	fmt.Println("  rollback  - Rollback the last applied migration")
	fmt.Println("  status    - Show migration status")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migrate/main.go migrate")
	fmt.Println("  go run cmd/migrate/main.go status")
	fmt.Println("  go run cmd/migrate/main.go rollback")
}

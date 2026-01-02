package main

import (
	"digital-wallet-api/internal/config"
	"digital-wallet-api/internal/database"
	"digital-wallet-api/internal/handler"
	"digital-wallet-api/internal/repository"
	"digital-wallet-api/internal/router"
	"digital-wallet-api/internal/service"
	"digital-wallet-api/pkg/logger"
	"fmt"
	"os"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	enableJSON := cfg.LogFormat == "json"
	logger.InitLogger(cfg.LogLevel, enableJSON)

	logger.Info("Starting Digital Wallet API", map[string]interface{}{
		"environment": cfg.Environment,
		"port":        cfg.ServerPort,
	})

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
		logger.Fatal("Failed to connect to database", map[string]interface{}{
			"error": err.Error(),
		})
	}
	defer database.Close()

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		logger.Fatal("Failed to run migrations", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Initialize repositories
	db := database.GetDB()
	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	budgetRepo := repository.NewBudgetRepository(db)

	// Parse JWT expiration
	jwtExpiry, err := time.ParseDuration(cfg.JWTExpiration)
	if err != nil {
		logger.Fatal("Invalid JWT expiration format", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Initialize services
	authService := service.NewAuthService(userRepo, walletRepo, cfg.JWTSecret, jwtExpiry)
	walletService := service.NewWalletService(walletRepo)
	transactionService := service.NewTransactionService(walletRepo, transactionRepo, db)
	budgetService := service.NewBudgetService(budgetRepo, categoryRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	walletHandler := handler.NewWalletHandler(walletService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	budgetHandler := handler.NewBudgetHandler(budgetService)

	// Setup router
	r := router.SetupRouter(router.RouterConfig{
		AuthHandler:        authHandler,
		WalletHandler:      walletHandler,
		TransactionHandler: transactionHandler,
		BudgetHandler:      budgetHandler,
		JWTSecret:          cfg.JWTSecret,
	})

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	logger.Info("Server starting", map[string]interface{}{
		"address": addr,
	})

	if err := r.Run(addr); err != nil {
		logger.Fatal("Failed to start server", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

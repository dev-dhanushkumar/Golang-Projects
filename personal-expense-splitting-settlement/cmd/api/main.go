package main

import (
	"fmt"
	"os"
	"personal-expense-splitting-settlement/internal/config"
	"personal-expense-splitting-settlement/internal/database"
	"personal-expense-splitting-settlement/internal/handler"
	"personal-expense-splitting-settlement/internal/repository"
	"personal-expense-splitting-settlement/internal/router"
	"personal-expense-splitting-settlement/internal/services"
	"personal-expense-splitting-settlement/pkg/logger"
	"time"
)

func main() {
	// Initialize the logger early in the application
	sugar, err := logger.InitLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Ensure all buffered logs are written before exiting
	defer sugar.Sync()
	sugar.Info("Application starting up...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		sugar.Fatalw("Failed to load config", "error", err)
	}

	sugar.Debugw("Configuration loaded successfully", "db_host", cfg.DBHost)

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
		sugar.Fatalw("Failed to connect to database", "error", err, "db_name", cfg.DBName)
	}

	sugar.Info("Database connection established successfully")
	defer database.Close()

	// Run SQL migrations
	if err := database.RunMigrationsFromFiles(sugar); err != nil {
		sugar.Fatalw("Failed to run migrations", "error", err, "db_name", cfg.DBName)
	}

	// Initialize repository
	db := database.GetDB()
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	friendshipRepo := repository.NewFriendshipRepository(db)

	// Parse JWT expiration
	jwtExpiry, err := time.ParseDuration(cfg.JWTExpiration)
	if err != nil {
		sugar.Fatalln("Invalid JWT expiration format", "error", err)
	}

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, jwtExpiry, cfg.DataSecret)
	sessionService := services.NewSessionService(sessionRepo, cfg.JWTSecret, jwtExpiry)
	friendshipService := services.NewFriendshipService(friendshipRepo, userRepo)

	// Initialize Handler
	authHandler := handler.NewAuthHandler(authService, sessionService)
	friendshipHandler := handler.NewFriendshipHandler(friendshipService)

	// Setup Router
	r := router.SetupRouter(router.RouterConfig{
		AuthHandler:       authHandler,
		FriendshipHandler: friendshipHandler,
		JWTSecret:         cfg.JWTSecret,
		Logger:            sugar,
	})

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	sugar.Info("Server Starting", "address", addr)

	if err := r.Run(addr); err != nil {
		sugar.Fatalw("Failed to start server", "error", err)
	}
}

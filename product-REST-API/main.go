package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dev-dhanushkumar/Golang-Projects/product-api/db"
	"github.com/dev-dhanushkumar/Golang-Projects/product-api/handlers"
)

func main() {
	// connect to MongoDB
	client, err := db.GetMongoClient()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	productCollection := db.GetProductsCollection(client)
	// fmt.Println("Database connected Successfully Final")

	// Create product handler
	productHandler := handlers.NewProductHandler(productCollection)

	// Set up routes
	mux := http.NewServeMux()
	mux.HandleFunc("/products", productHandler.HandleProducts)
	mux.HandleFunc("/products/", productHandler.HandleProducts)

	// Configure server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Disconnect MongoDB client
	if err := client.Disconnect(ctx); err != nil {
		log.Fatalf("Error disconnecting from MongoDB: %v", err)
	}

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")

}

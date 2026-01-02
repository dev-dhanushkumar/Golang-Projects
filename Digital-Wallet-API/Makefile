.PHONY: help run build test clean migrate-up migrate-down dev

help: ## Display this help screen
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Run the application
	@echo "Starting application..."
	@go run cmd/api/main.go

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/digital-wallet-api cmd/api/main.go
	@echo "Build complete: bin/digital-wallet-api"

dev: ## Run in development mode with hot reload (requires air)
	@echo "Starting development mode..."
	@air

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

install: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

migrate: ## Run database migrations
	@echo "Running database migrations..."
	@go run cmd/api/main.go
	@echo "Migrations complete"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t digital-wallet-api:latest .
	@echo "Docker image built"

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env digital-wallet-api:latest

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

format: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted"

# 🚀 Quick Start Guide

## Prerequisites Checklist

Before you begin, ensure you have:

- ✅ Go 1.21 or higher installed
- ✅ PostgreSQL 12+ installed and running
- ✅ Git (for cloning)
- ✅ A terminal/command line interface

## Step-by-Step Setup

### 1. Verify Go Installation

```bash
go version
# Should output: go version go1.21 or higher
```

### 2. Verify PostgreSQL Installation

```bash
psql --version
# Should output: psql (PostgreSQL) 12.x or higher
```

### 3. Create PostgreSQL Database

**Option A: Using psql command line**

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE digital_wallet;

# Create user (optional)
CREATE USER wallet_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE digital_wallet TO wallet_user;

# Exit psql
\q
```

**Option B: Using createdb command**

```bash
createdb -U postgres digital_wallet
```

### 4. Setup Environment Variables

```bash
# Copy the example env file
cp .env.example .env

# Edit .env with your preferred editor
nano .env  # or vim .env or code .env
```

**Required configurations in .env:**

```env
# Server Configuration
PORT=8080
ENVIRONMENT=development

# Database Configuration (IMPORTANT: Update these!)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres           # Change to your PostgreSQL username
DB_PASSWORD=your_password  # Change to your PostgreSQL password
DB_NAME=digital_wallet
DB_SSL_MODE=disable

# JWT Configuration (IMPORTANT: Change in production!)
JWT_SECRET=your-super-secret-jwt-key-minimum-32-characters-long
JWT_EXPIRATION=24h

# Logger Configuration
LOG_LEVEL=debug
LOG_FORMAT=text
```

### 5. Install Dependencies

```bash
# Download all Go modules
go mod download

# Tidy up dependencies
go mod tidy
```

### 6. Run the Application

```bash
# Option 1: Using go run
go run cmd/api/main.go

# Option 2: Using Makefile
make run

# Option 3: Build then run
make build
./bin/digital-wallet-api
```

The application will:
- ✅ Connect to PostgreSQL
- ✅ Enable UUID extensions (uuid-ossp, pgcrypto)
- ✅ Auto-migrate database tables
- ✅ Seed default categories
- ✅ Start HTTP server on configured port

### 7. Verify Installation

Open another terminal and test the health endpoint:

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "service": "digital-wallet-api",
  "status": "ok"
}
```

## 🧪 Testing the API

### 1. Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User",
    "phone": "1234567890"
  }'
```

Expected response:
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "uuid-here",
    "email": "test@example.com",
    "full_name": "Test User",
    "phone": "1234567890",
    "is_active": true
  }
}
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

Save the `token` from the response - you'll need it for authenticated requests!

### 3. Get Wallet Balance (Authenticated)

```bash
curl -X GET http://localhost:8080/api/v1/wallet/balance \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 4. Add Money to Wallet

```bash
curl -X POST http://localhost:8080/api/v1/transactions/credit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "amount": 1000.00,
    "description": "Initial deposit"
  }'
```

### 5. Check Transaction History

```bash
curl -X GET http://localhost:8080/api/v1/transactions \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## 🐛 Troubleshooting

### Problem: "Failed to connect to database"

**Solution:**
1. Verify PostgreSQL is running: `pg_isready`
2. Check database exists: `psql -U postgres -l | grep digital_wallet`
3. Verify credentials in `.env` file
4. Check PostgreSQL is listening on port 5432

### Problem: "uuid-ossp extension not found"

**Solution:**
```sql
-- Connect to your database
psql -U postgres -d digital_wallet

-- Create extension manually
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
```

### Problem: Port 8080 already in use

**Solution:**
Change the PORT in `.env` file:
```env
PORT=8081  # or any other available port
```

### Problem: go.mod issues

**Solution:**
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download
go mod tidy
```

## 📚 Next Steps

1. ✅ **Explore the API**: Check [README.md](README.md) for complete API documentation
2. ✅ **Test with Postman**: Import the API endpoints and create a collection
3. ✅ **Read the code**: Start with `cmd/api/main.go` to understand the flow
4. ✅ **Customize**: Add your own features and endpoints

## 🛠️ Development Commands

```bash
# Run tests
make test

# Run with coverage
make test-coverage

# Format code
make format

# Build binary
make build

# Clean build artifacts
make clean

# See all available commands
make help
```

## 🎓 Understanding UUID Implementation

This project uses UUIDs instead of auto-incrementing integers:

**Benefits:**
- 🔐 Better security (non-sequential)
- 🌐 Distributed system friendly
- 🔗 No collision across databases
- 📊 Better for microservices architecture

**Database:**
- PostgreSQL automatically generates UUIDs using `gen_random_uuid()`
- All foreign keys use UUID type
- Indexes work efficiently with UUIDs

**API:**
- All IDs in requests/responses are UUID strings
- Example: `"id": "550e8400-e29b-41d4-a716-446655440000"`

## 📧 Support

For issues or questions:
1. Check the troubleshooting section
2. Review the [README.md](README.md)
3. Check Go and PostgreSQL documentation

---

**Happy Coding! 🎉**

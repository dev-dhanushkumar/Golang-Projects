# Digital Wallet & Expense Management API

A comprehensive digital wallet and expense management REST API built with Go, featuring JWT authentication, transaction management, budget tracking, and more.

## 🚀 Features

- **User Authentication**: Register, login with JWT tokens
- **Digital Wallet**: Manage wallet balance with credit/debit operations
- **Transactions**: Track all money movements with detailed history
- **Money Transfer**: Send money to other users securely
- **Budget Management**: Set and track budgets by category
- **Expense Categories**: Organize expenses with custom categories
- **Transaction Summary**: Get detailed financial reports
- **UUID-based IDs**: All entities use UUIDs instead of integers for better security

## 📋 Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT (JSON Web Tokens)
- **Logger**: Custom structured logger
- **Validation**: go-playground/validator
- **Password Hashing**: bcrypt
- **Decimal Handling**: shopspring/decimal

## 🏗️ Project Structure

```
digital-wallet-api/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/                     # Configuration management
│   ├── database/                   # Database connection & migrations
│   ├── models/                     # Data models (with UUID)
│   ├── repository/                 # Data access layer
│   ├── service/                    # Business logic
│   ├── handler/                    # HTTP handlers
│   ├── middleware/                 # HTTP middleware
│   ├── dto/                        # Data transfer objects
│   └── router/                     # Route definitions
├── pkg/
│   ├── logger/                     # Custom logger
│   ├── utils/                      # Utility functions
│   └── validator/                  # Validation helpers
├── .env.example                    # Environment variables template
├── .gitignore
├── Makefile                        # Build and development commands
├── go.mod
└── README.md
```

## 🛠️ Installation & Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile commands)

### 1. Clone the repository

```bash
git clone <repository-url>
cd Digital-Wallet-API
```

### 2. Install dependencies

```bash
go mod download
# or
make install
```

### 3. Configure environment variables

```bash
cp .env.example .env
```

Edit `.env` file with your configuration:

```env
# Server
PORT=8080
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=digital_wallet
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production-min-32-chars
JWT_EXPIRATION=24h

# Logger
LOG_LEVEL=debug
LOG_FORMAT=text
```

### 4. Create PostgreSQL database

```bash
createdb digital_wallet
```

Or using SQL:

```sql
CREATE DATABASE digital_wallet;
```

### 5. Run the application

The application will automatically run migrations on startup.

```bash
# Using Go
go run cmd/api/main.go

# Or using Make
make run

# Or build and run
make build
./bin/digital-wallet-api
```

## 📡 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Register new user | No |
| POST | `/api/v1/auth/login` | Login user | No |
| GET | `/api/v1/auth/profile` | Get user profile | Yes |
| PUT | `/api/v1/auth/profile` | Update profile | Yes |

### Wallet

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/wallet` | Get wallet info | Yes |
| GET | `/api/v1/wallet/balance` | Get balance | Yes |

### Transactions

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/transactions/credit` | Add money | Yes |
| POST | `/api/v1/transactions/debit` | Spend money | Yes |
| POST | `/api/v1/transactions/transfer` | Transfer to user | Yes |
| GET | `/api/v1/transactions` | Get history | Yes |
| GET | `/api/v1/transactions/:id` | Get transaction | Yes |
| GET | `/api/v1/transactions/summary` | Get summary | Yes |

### Budgets

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/budgets` | Create budget | Yes |
| GET | `/api/v1/budgets` | Get all budgets | Yes |
| GET | `/api/v1/budgets/:id` | Get budget | Yes |
| PUT | `/api/v1/budgets/:id` | Update budget | Yes |
| DELETE | `/api/v1/budgets/:id` | Delete budget | Yes |
| GET | `/api/v1/budgets/alerts` | Get budget alerts | Yes |

## 📝 API Usage Examples

### Register User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "full_name": "John Doe",
    "phone": "1234567890"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Credit Money (Add to Wallet)

```bash
curl -X POST http://localhost:8080/api/v1/transactions/credit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "amount": 1000.00,
    "description": "Salary deposit"
  }'
```

### Debit Money (Spend)

```bash
curl -X POST http://localhost:8080/api/v1/transactions/debit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "amount": 50.00,
    "category_id": "category-uuid-here",
    "description": "Lunch"
  }'
```

## 🔧 Development

### Run tests

```bash
make test
```

### Run with coverage

```bash
make test-coverage
```

### Format code

```bash
make format
```

### Build binary

```bash
make build
```

## 🗄️ Database Schema (UUID-based)

All tables use UUID as primary keys instead of auto-incrementing integers.

### Key Features:
- All IDs are UUID v4
- PostgreSQL's `gen_random_uuid()` used for auto-generation
- Foreign key relationships use UUID references
- Better security and distribution compatibility

## 🔐 Security Features

- JWT-based authentication
- Password hashing with bcrypt
- UUID-based IDs (non-sequential, harder to guess)
- CORS middleware
- Request validation
- Error handling middleware
- Structured logging with request IDs

## 📊 Logging

The application uses a custom structured logger that supports:
- Multiple log levels (DEBUG, INFO, WARN, ERROR, FATAL)
- JSON and text formats
- Colored console output
- Context logging with fields
- Request tracking with unique IDs

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

MIT License

## 👨‍💻 Author

Dhanush Kumar

## 🙏 Acknowledgments

- Gin Web Framework
- GORM
- shopspring/decimal
- golang-jwt

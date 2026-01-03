# Personal Expense Splitting & Settlement - Documentation

Welcome to the documentation for the Personal Expense Splitting & Settlement application!

---

## 📚 Documentation Structure

### 🔐 Authentication Module
Located in root `docs/` folder:

- **[AUTH_MODULE_REVIEW.md](AUTH_MODULE_REVIEW.md)** - Complete authentication system review
- **[API_TESTING.md](API_TESTING.md)** - Authentication API testing guide
- **[QUICK_REFERENCE.md](QUICK_REFERENCE.md)** - Quick reference for common operations
- **[LOGGING_MIGRATIONS_GUIDE.md](LOGGING_MIGRATIONS_GUIDE.md)** - Detailed guide for logging and migrations
- **[QUICKSTART_LOGGING_MIGRATIONS.md](QUICKSTART_LOGGING_MIGRATIONS.md)** - Quick start for logging and migrations

---

### 🤝 Friendship Module  
Located in `docs/friendship-module/`:

- **[FRIENDSHIP_API_TESTING.md](friendship-module/FRIENDSHIP_API_TESTING.md)** - Complete friendship API testing guide
  - Send/Accept/Reject friend requests
  - Block users
  - View friends list
  - Manage pending requests

---

### 📊 Logging System
Located in `docs/logger/`:

- **[README.md](logger/README.md)** - Overview and quick start
- **[LOGGING_SYSTEM.md](logger/LOGGING_SYSTEM.md)** - Complete logging documentation
- **[QUICK_REFERENCE.md](logger/QUICK_REFERENCE.md)** - Common logging commands
- **[IMPLEMENTATION_SUMMARY.md](logger/IMPLEMENTATION_SUMMARY.md)** - Implementation details

**Features**:
- ✅ Request/Response tracking with unique IDs
- ✅ Duration metrics for all requests
- ✅ Full request/response body logging
- ✅ JSON format for easy querying
- ✅ Dual output (console + file)

---

## 🚀 Quick Start

### 1. Start the Server
```bash
go run ./cmd/api/main.go
```

### 2. Monitor Logs
```bash
tail -f logs/app.log | jq .
```

### 3. Test Authentication
```bash
# Register a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@test.com",
    "password": "Test123!",
    "first_name": "Test",
    "last_name": "User"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@test.com",
    "password": "Test123!"
  }'

# Save token
export TOKEN="<your_access_token>"

# Test protected endpoint
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Test Friendships
```bash
# Send friend request
curl -X POST http://localhost:8080/api/v1/friends/request \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"friend_email":"friend@test.com"}'

# View pending requests
curl -X GET http://localhost:8080/api/v1/friends/pending \
  -H "Authorization: Bearer $TOKEN"

# View friends list
curl -X GET http://localhost:8080/api/v1/friends \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📋 Implementation Status

### ✅ Completed Modules

#### Phase 0: Authentication (8 endpoints)
- ✅ User registration with encrypted phone numbers
- ✅ Login with session creation
- ✅ Token refresh with rotation
- ✅ Logout with session revocation
- ✅ Get user profile
- ✅ Update user profile
- ✅ Get active sessions
- ✅ Strong password validation

#### Phase 1: Friendships (7 endpoints)
- ✅ Send friend requests
- ✅ Accept friend requests
- ✅ Reject friend requests
- ✅ Block users
- ✅ Remove friends
- ✅ List friends
- ✅ View pending requests

#### Infrastructure
- ✅ Manual SQL migrations with version tracking
- ✅ File + console logging with Zap
- ✅ Request/Response logging with unique IDs
- ✅ Password encryption (Bcrypt)
- ✅ Phone number encryption (AES-256)
- ✅ JWT authentication
- ✅ Session management

---

### 🔄 Upcoming Modules

#### Phase 2: Groups (8 endpoints)
- Create/update/delete groups
- Add/remove members
- Role-based permissions
- Group details

#### Phase 3: Expenses (6 endpoints)
- Add expenses with split methods
- Update/delete expenses
- View group/user expenses
- Split calculations (equal, exact, percentage, shares)

#### Phase 4: Balance & Settlement (7 endpoints)
- View balances
- Settlement suggestions
- Record settlements
- Minimize transactions

#### Phase 5-7: Additional Features
- Notifications
- Payment methods
- Audit logs & analytics

---

## 🗂️ Project Structure

```
project-root/
├── cmd/
│   ├── api/
│   │   └── main.go                  # Application entry point
│   └── migrate/
│       └── main.go                  # Migration CLI tool
├── internal/
│   ├── config/                      # Configuration management
│   ├── database/                    # Database connection & migrations
│   ├── dto/                         # Data Transfer Objects
│   ├── handler/                     # HTTP handlers
│   ├── middleware/                  # Auth & logging middleware
│   ├── models/                      # Database models
│   ├── repository/                  # Data access layer
│   ├── router/                      # Route definitions
│   └── services/                    # Business logic
├── pkg/
│   ├── logger/                      # Logging utilities
│   ├── utils/                       # Helper utilities
│   └── validator/                   # Custom validators
├── migrations/                      # SQL migration files
├── logs/                            # Application logs
└── docs/                            # Documentation (this folder)
```

---

## 🔧 Configuration

Environment variables (`.env` file):
```bash
# Server
SERVER_PORT=8080
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=personal-ess
DB_SSL_MODE=disable

# Security
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION=168h  # 7 days
DATA_SECRET=your-32-byte-encryption-key
```

---

## 📊 Key Features

### Security
- 🔐 JWT-based authentication
- 🔄 Automatic token rotation
- 🔒 Bcrypt password hashing (cost 12)
- 🔑 AES-256-GCM phone encryption
- ✅ Custom password validation
- 📝 Session tracking with IP/User-Agent

### Database
- 🗄️ PostgreSQL with UUID primary keys
- 📝 Manual SQL migrations with version tracking
- ⚡ Optimized indexes
- 🔗 Proper foreign key constraints
- ⏰ Automatic timestamp triggers

### Logging
- 📊 Structured JSON logging (Zap)
- 🆔 Unique request IDs
- ⏱️ Duration tracking
- 📝 Request/Response body logging
- 🎯 Automatic log levels (INFO/WARN/ERROR)
- 🖥️ Dual output (console + file)

### Architecture
- 🏗️ Clean architecture (Handler→Service→Repository→DB)
- 💉 Dependency injection
- 🎯 Single responsibility principle
- ✨ Consistent error handling
- 📦 Modular design

---

## 🧪 Testing

### Run the Application
```bash
go run ./cmd/api/main.go
```

### Run Migrations
```bash
# Check migration status
go run ./cmd/migrate/main.go status

# Run migrations
go run ./cmd/migrate/main.go migrate

# Rollback last migration
go run ./cmd/migrate/main.go rollback
```

### Database Access
```bash
psql -h localhost -U postgres -d personal-ess
```

---

## 📈 Monitoring & Debugging

### View Logs
```bash
# Real-time monitoring
tail -f logs/app.log | jq .

# Find errors
cat logs/app.log | jq 'select(.level == "error")'

# Track request
cat logs/app.log | jq 'select(.request_id == "YOUR_ID")'

# Find slow requests
cat logs/app.log | jq 'select(.duration_ms > 1000)'
```

### Database Queries
```sql
-- View users
SELECT id, email, first_name, last_name FROM users;

-- View friendships
SELECT f.id, f.status, 
       u1.email as user_1, 
       u2.email as user_2
FROM friendships f
JOIN users u1 ON f.user_id_1 = u1.id
JOIN users u2 ON f.user_id_2 = u2.id;

-- View active sessions
SELECT s.id, u.email, s.expire_at, s.revoked_at
FROM user_sessions s
JOIN users u ON s.user_id = u.id
WHERE s.revoked_at IS NULL;
```

---

## 🆘 Troubleshooting

### Common Issues

**Server won't start**:
- Check PostgreSQL is running
- Verify database credentials in `.env`
- Ensure port 8080 is available

**Migration errors**:
- Check database connection
- Verify UUID extension is enabled
- Review migration SQL syntax

**Authentication fails**:
- Verify JWT_SECRET is set
- Check token expiration
- Ensure session is not revoked

**Logs not appearing**:
- Check `logs/` directory exists
- Verify file permissions
- Check ENVIRONMENT variable

---

## 📚 Additional Resources

- **[IMPLEMENTATION_PLAN.md](../IMPLEMENTATION_PLAN.md)** - Complete roadmap for remaining features
- **[ROADMAP.md](../ROADMAP.md)** - Visual implementation timeline
- **[plan.md](../plan.md)** - Original project specification

---

## 🎯 Next Steps

1. **Complete Phase 2**: Groups module (8 endpoints)
2. **Complete Phase 3**: Expenses module (6 endpoints) - CRITICAL
3. **Complete Phase 4**: Balance & Settlement (7 endpoints) - CRITICAL
4. **Add remaining features**: Notifications, Payment Methods, Audit

---

## 🤝 Contributing

When adding new features:
1. Follow existing code patterns
2. Add migrations for schema changes
3. Update relevant documentation
4. Test all endpoints
5. Check logs for errors

---

## ✅ Summary

This application provides a robust foundation for expense splitting with:
- ✅ Secure authentication with session management
- ✅ Social features (friendships)
- ✅ Comprehensive logging and monitoring
- ✅ Clean architecture and best practices
- ✅ Complete documentation

**Current Progress**: 15/48 endpoints (31% complete)

**Ready to proceed with Phase 2: Groups Module!**

---

For detailed information on specific topics, refer to the relevant documentation files listed above.

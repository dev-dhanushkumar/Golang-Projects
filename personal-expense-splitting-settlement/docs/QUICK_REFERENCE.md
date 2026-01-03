# Authentication Module - Quick Reference

## 🚀 Quick Start

```bash
# 1. Start PostgreSQL
docker container start postgres_arch

# 2. Run the app
go run cmd/api/main.go

# 3. Test health
curl http://localhost:8080/health
```

---

## 📡 API Endpoints

### Public Endpoints
```bash
POST   /api/v1/auth/register     # Register new user
POST   /api/v1/auth/login        # Login user
POST   /api/v1/auth/refresh      # Refresh token
```

### Protected Endpoints (Requires Bearer Token)
```bash
POST   /api/v1/auth/logout       # Logout current session
GET    /api/v1/auth/me           # Get current user
GET    /api/v1/auth/sessions     # Get active sessions
GET    /api/v1/users/me          # Get user profile
PATCH  /api/v1/users/me          # Update user profile
```

---

## 🔑 Authentication Flow

### 1. Register
```json
POST /api/v1/auth/register
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "+1234567890"
}
```

### 2. Login
```json
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "user": { ... },
    "tokens": {
      "access_token": "eyJ...",
      "refresh_token": "xyz...",
      "expires_in": 86400
    }
  }
}
```

### 3. Use Protected Endpoint
```bash
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  http://localhost:8080/api/v1/users/me
```

### 4. Refresh Token
```json
POST /api/v1/auth/refresh
{
  "refresh_token": "YOUR_REFRESH_TOKEN"
}
```

### 5. Logout
```bash
POST /api/v1/auth/logout
Authorization: Bearer YOUR_ACCESS_TOKEN
```

---

## 🔒 Password Requirements

✅ **Valid Password Must Have**:
- Minimum 8 characters
- At least 1 uppercase letter (A-Z)
- At least 1 lowercase letter (a-z)
- At least 1 number (0-9)
- At least 1 special character (!@#$%^&*()_+-=[]{};':"\\|,.<>/?)

**Examples**:
- ✅ `SecurePass123!`
- ✅ `MyP@ssw0rd`
- ❌ `weak`
- ❌ `password123` (no uppercase, no special)

---

## 🗂️ Project Structure

```
personal-expense-splitting-settlement/
├── cmd/api/
│   └── main.go                    # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration loading
│   ├── database/
│   │   └── database.go            # Database connection & migrations
│   ├── dto/
│   │   ├── auth_dto.go            # Request/Response DTOs
│   │   └── session_dto.go
│   ├── handler/
│   │   └── auth_handler.go        # HTTP handlers
│   ├── middleware/
│   │   └── auth_middleware.go     # JWT authentication
│   ├── models/
│   │   ├── user.go                # User model
│   │   ├── session.go             # Session model
│   │   └── errors.go              # Custom errors
│   ├── repository/
│   │   ├── user_repository.go     # User data access
│   │   └── session_repository.go  # Session data access
│   ├── router/
│   │   └── router.go              # Route definitions
│   └── services/
│       ├── auth_service.go        # Auth business logic
│       └── session_service.go     # Session business logic
├── pkg/
│   ├── logger/
│   │   └── logger.go              # Logging setup
│   ├── utils/
│   │   ├── jwt.go                 # JWT utilities
│   │   ├── random.go              # Random generation
│   │   └── response.go            # HTTP response helpers
│   └── validator/
│       └── validator.go           # Custom validators
└── plan.md                        # Project plan
```

---

## 🔧 Configuration (.env)

```bash
# Server
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=personal-ess
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION=24h

# Encryption
DATA_SECRET=base64-encoded-32-byte-key

# Environment
ENVIRONMENT=development
```

---

## 📦 Database Tables

### users
```sql
id              UUID PRIMARY KEY
email           VARCHAR UNIQUE NOT NULL
password_hash   VARCHAR NOT NULL
first_name      VARCHAR
last_name       VARCHAR
phone_number    VARCHAR (ENCRYPTED)
profile_image_url VARCHAR
default_currency VARCHAR DEFAULT 'USD'
email_verified  BOOLEAN DEFAULT false
is_active       BOOLEAN DEFAULT true
created_at      TIMESTAMP
updated_at      TIMESTAMP
last_login_at   TIMESTAMP
```

### user_sessions
```sql
id                  UUID PRIMARY KEY
user_id             UUID FK -> users.id
token_hash          VARCHAR INDEXED
refresh_token_hash  VARCHAR
ip_address          INET
user_agent          TEXT
expire_at           TIMESTAMP
created_at          TIMESTAMP
revoked_at          TIMESTAMP
```

---

## 🛠️ Common Commands

### Run Application
```bash
go run cmd/api/main.go
```

### Run with Hot Reload (Air)
```bash
air
```

### Database Commands
```bash
# Connect to PostgreSQL
docker exec -it postgres_arch psql -U postgres -d personal-ess

# Check users
SELECT * FROM users;

# Check sessions
SELECT * FROM user_sessions WHERE revoked_at IS NULL;

# Clean up test data
DELETE FROM user_sessions;
DELETE FROM users WHERE email LIKE 'test%';
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run tests
go test ./...

# Check dependencies
go mod tidy
go mod verify
```

---

## 🐛 Debugging

### Enable Verbose Logging
The app uses `zap.NewDevelopment()` which provides detailed logs.

### Common Issues

**1. Database Connection Error**
```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Restart if needed
docker container start postgres_arch
```

**2. JWT Token Invalid**
- Check JWT_SECRET in .env matches
- Verify token hasn't expired
- Use refresh token to get new access token

**3. Password Validation Fails**
- Ensure password meets all requirements
- Check custom validator is registered

**4. Session Not Found**
- Check session hasn't expired
- Verify session wasn't revoked on logout

---

## 📝 Code Examples

### Create Custom Handler
```go
func (h *AuthHandler) YourHandler(ctx *gin.Context) {
    // Get user from context (set by middleware)
    userID, _ := ctx.Get("user_id")
    
    // Your logic here
    
    // Return success
    utils.OK(ctx, "Success message", data)
}
```

### Add New Repository Method
```go
func (r *userRepository) YourMethod(id uuid.UUID) error {
    return r.db.Where("id = ?", id).First(&model).Error
}
```

### Add New Service Method
```go
func (s *authService) YourService(param string) (*Response, error) {
    // Business logic
    result, err := s.userRepo.YourMethod(param)
    if err != nil {
        return nil, err
    }
    return result, nil
}
```

---

## ✅ Testing Checklist

Before considering auth module complete:

- [ ] Register new user works
- [ ] Strong password validation works
- [ ] Weak password rejected
- [ ] Login with correct credentials works
- [ ] Login with wrong password fails
- [ ] Access protected endpoint with token works
- [ ] Access protected endpoint without token fails
- [ ] Update profile works
- [ ] Get active sessions works
- [ ] Token refresh works
- [ ] Old token invalidated after refresh
- [ ] Logout works
- [ ] Token invalid after logout
- [ ] Phone number encrypted in database
- [ ] LastLoginAt updated on login

---

## 🎯 Next Module: Groups

After auth is complete, implement:

1. **Groups Module**
   - Create group
   - Add/remove members
   - Get group details
   - List user groups

2. **Friendships Module**
   - Send friend request
   - Accept/reject request
   - List friends
   - Get friend balance

---

## 📚 Resources

- **Plan**: See [plan.md](plan.md)
- **Full Review**: See [AUTH_MODULE_REVIEW.md](AUTH_MODULE_REVIEW.md)
- **Testing Guide**: See [API_TESTING.md](API_TESTING.md)
- **Gin Framework**: https://gin-gonic.com/docs/
- **GORM**: https://gorm.io/docs/
- **JWT**: https://github.com/golang-jwt/jwt

---

## 💡 Tips

1. **Always use the middleware** for protected routes
2. **Validate all inputs** using the validator package
3. **Handle errors properly** with custom error types
4. **Use transactions** for multi-step operations
5. **Log important events** for debugging
6. **Never log sensitive data** (passwords, tokens)
7. **Test with different scenarios** (happy path & error cases)

---

**Need Help?**
- Check the review document
- Read the code comments
- Test with provided examples
- Review the plan.md

**Last Updated**: January 3, 2026

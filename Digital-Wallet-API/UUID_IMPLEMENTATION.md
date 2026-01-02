# UUID Implementation Summary

## Overview

This Digital Wallet API project has been implemented with **UUID (Universally Unique Identifier)** as the primary key for all entities instead of traditional auto-incrementing integers.

## Key Changes from Integer IDs to UUIDs

### 1. **Base Model** (internal/models/base.go)

**Before (Integer):**
```go
type BaseModel struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**After (UUID):**
```go
type BaseModel struct {
    ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
    if base.ID == uuid.Nil {
        base.ID = uuid.New()
    }
    return nil
}
```

### 2. **Model Relationships**

All foreign key relationships now use UUID:

**User Model:**
```go
type User struct {
    BaseModel  // UUID ID
    // ... fields
}
```

**Wallet Model:**
```go
type Wallet struct {
    BaseModel
    UserID   uuid.UUID       `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
    // ... other fields
}
```

**Transaction Model:**
```go
type Transaction struct {
    BaseModel
    WalletID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"wallet_id"`
    CategoryID *uuid.UUID `gorm:"type:uuid;index" json:"category_id,omitempty"` // Nullable
    // ... other fields
}
```

### 3. **Database Migrations**

PostgreSQL extensions are automatically enabled:

```go
// Enable UUID extension for PostgreSQL
DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
DB.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"")
```

### 4. **Repository Layer**

All repository methods updated to work with UUID:

**Before:**
```go
func (r *userRepository) FindByID(id uint) (*models.User, error)
```

**After:**
```go
func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error)
```

### 5. **Service Layer**

Service methods accept and return UUIDs:

```go
func (s *authService) GetProfile(userID uuid.UUID) (*dto.UserResponse, error)
func (s *transactionService) Credit(userID uuid.UUID, req dto.CreditRequest) (*dto.TransactionResponse, error)
```

### 6. **Handler Layer**

Handlers parse UUID strings from request parameters:

```go
func (h *AuthHandler) GetProfile(c *gin.Context) {
    userID := c.GetString("user_id")
    uid, err := uuid.Parse(userID)  // Parse string to UUID
    if err != nil {
        utils.BadRequest(c, "Invalid user ID", err)
        return
    }
    // ... use uid
}
```

### 7. **DTOs (Data Transfer Objects)**

All DTOs use uuid.UUID type:

```go
type UserResponse struct {
    ID       uuid.UUID `json:"id"`
    Email    string    `json:"email"`
    FullName string    `json:"full_name"`
    // ...
}
```

### 8. **JWT Token Claims**

JWT claims store UUID for user identification:

```go
type Claims struct {
    UserID uuid.UUID `json:"user_id"`
    Email  string    `json:"email"`
    jwt.RegisteredClaims
}
```

## Benefits of UUID Implementation

### 1. **Security**
- ✅ Non-sequential IDs prevent enumeration attacks
- ✅ Harder to guess valid IDs
- ✅ No information leakage about entity count

**Example:**
- Integer: `/api/v1/transactions/123` → easy to guess 122, 124
- UUID: `/api/v1/transactions/550e8400-e29b-41d4-a716-446655440000` → impossible to guess

### 2. **Distributed Systems**
- ✅ Generate IDs on any server without coordination
- ✅ No central ID authority needed
- ✅ Perfect for microservices architecture
- ✅ Can merge databases without ID conflicts

### 3. **Scalability**
- ✅ No bottleneck from auto-increment sequences
- ✅ Better for horizontal scaling
- ✅ Can generate IDs client-side if needed

### 4. **Database Performance**
- ✅ No lock contention on ID generation
- ✅ Indexes still work efficiently
- ✅ Better for distributed databases

## API Response Examples

### Register User Response (UUID)
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "email": "user@example.com",
    "full_name": "John Doe",
    "phone": "1234567890",
    "is_active": true
  }
}
```

### Transaction Response (UUID)
```json
{
  "success": true,
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "wallet_id": "b2c3d4e5-f678-90ab-cdef-1234567890ab",
    "category_id": "c3d4e5f6-7890-abcd-ef12-34567890abcd",
    "type": "credit",
    "amount": "1000.00",
    "balance_after": "1000.00",
    "description": "Initial deposit",
    "transaction_date": "2026-01-02T10:30:00Z"
  }
}
```

## Database Schema (PostgreSQL)

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### Wallets Table
```sql
CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    balance DECIMAL(15,2) DEFAULT 0.00,
    currency VARCHAR(3) DEFAULT 'USD',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(user_id)
);
```

### Transactions Table
```sql
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    type VARCHAR(20) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    balance_after DECIMAL(15,2) NOT NULL,
    description TEXT,
    reference_id VARCHAR(100) UNIQUE,
    status VARCHAR(20) DEFAULT 'completed',
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## Usage in Code

### Creating a New Entity
```go
// UUID is automatically generated
user := &models.User{
    Email:    "user@example.com",
    FullName: "John Doe",
}
// Before create hook will generate UUID if not set
db.Create(user)
// user.ID is now a valid UUID
```

### Finding by UUID
```go
userID := uuid.MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479")
user, err := userRepo.FindByID(userID)
```

### Parsing UUID from String (API requests)
```go
idStr := c.Param("id")
id, err := uuid.Parse(idStr)
if err != nil {
    // Handle invalid UUID
    return
}
```

## Migration from Integer to UUID

If you were to migrate an existing system:

1. **Add UUID columns** alongside integer IDs
2. **Populate UUIDs** for existing records
3. **Update foreign keys** to reference UUIDs
4. **Update application code** to use UUIDs
5. **Remove integer columns** once verified

## Performance Considerations

### Pros:
- No sequence lock contention
- Better for distributed systems
- Can generate offline

### Cons:
- Slightly larger storage (16 bytes vs 4-8 bytes)
- Indexes may be less compact
- Random UUIDs can cause page splits

### Optimization:
- Use UUID v1 or v6 for time-ordered UUIDs if needed
- Consider indexing strategies for large tables
- PostgreSQL handles UUIDs efficiently with proper indexes

## Testing UUID Implementation

```bash
# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'

# Response will contain UUID:
# "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479"

# Use UUID in subsequent requests
curl -X GET http://localhost:8080/api/v1/transactions/f47ac10b-58cc-4372-a567-0e02b2c3d479 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Summary

✅ **All entities use UUID as primary keys**
✅ **PostgreSQL extensions (uuid-ossp, pgcrypto) enabled**
✅ **Automatic UUID generation on entity creation**
✅ **All foreign key relationships use UUID**
✅ **API endpoints accept/return UUID strings**
✅ **JWT tokens store user UUID**
✅ **Better security and scalability**

---

**This implementation provides a modern, secure, and scalable foundation for the Digital Wallet API!** 🎉

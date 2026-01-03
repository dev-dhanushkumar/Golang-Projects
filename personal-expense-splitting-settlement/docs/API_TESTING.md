# API Testing Guide

## Quick Start

### 1. Start the Server

```bash
# Make sure PostgreSQL is running
docker container start postgres_arch

# Start the application
go run cmd/api/main.go
```

The server should start on `http://localhost:8080`

### 2. Test Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "Personal expense splitting and settlement"
}
```

---

## Authentication Flow Tests

### Test 1: Register New User

**Endpoint**: `POST /api/v1/auth/register`

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "SecurePass123!",
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "+1234567890"
  }'
```

**Expected Response** (201 Created):
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": "uuid-here",
      "email": "john.doe@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "created_at": "2026-01-03T10:00:00Z"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "random-base64-string",
      "expires_in": 86400
    }
  }
}
```

**Save the access_token and refresh_token for next tests!**

---

### Test 2: Invalid Password (Should Fail)

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test2@example.com",
    "password": "weak",
    "first_name": "Test",
    "last_name": "User"
  }'
```

**Expected Response** (400 Bad Request):
```json
{
  "success": false,
  "message": "Validation failed",
  "data": "Password validation failed"
}
```

**Reason**: Password must be min 8 chars with uppercase, lowercase, number, and special character.

---

### Test 3: Login

**Endpoint**: `POST /api/v1/auth/login`

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "SecurePass123!"
  }'
```

**Expected Response** (200 OK):
```json
{
  "success": true,
  "message": "Login Successfully",
  "data": {
    "user": {
      "id": "uuid-here",
      "email": "john.doe@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "created_at": "2026-01-03T10:00:00Z"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "random-base64-string",
      "expires_in": 86400
    }
  }
}
```

---

### Test 4: Get User Profile (Protected)

**Endpoint**: `GET /api/v1/users/me`

```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN_HERE"
```

**Expected Response** (200 OK):
```json
{
  "success": true,
  "message": "User profile retrived successfully",
  "data": {
    "id": "uuid-here",
    "email": "john.doe@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2026-01-03T10:00:00Z"
  }
}
```

**Without Token** (401 Unauthorized):
```bash
curl -X GET http://localhost:8080/api/v1/users/me
```

```json
{
  "success": false,
  "message": "Authorization header required"
}
```

---

### Test 5: Update Profile

**Endpoint**: `PATCH /api/v1/users/me`

```bash
curl -X PATCH http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jonathan",
    "default_currency": "EUR"
  }'
```

**Expected Response** (200 OK):
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "id": "uuid-here",
    "email": "john.doe@example.com",
    "first_name": "Jonathan",
    "last_name": "Doe",
    "created_at": "2026-01-03T10:00:00Z"
  }
}
```

---

### Test 6: Get Active Sessions

**Endpoint**: `GET /api/v1/auth/sessions`

```bash
curl -X GET http://localhost:8080/api/v1/auth/sessions \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN_HERE"
```

**Expected Response** (200 OK):
```json
{
  "success": true,
  "message": "Active sessions retrived seccessfully",
  "data": [
    {
      "id": "session-uuid",
      "user_id": "user-uuid",
      "ip_address": "127.0.0.1",
      "user_agent": "curl/7.81.0",
      "expires_at": "2026-01-10T10:00:00Z",
      "created_at": "2026-01-03T10:00:00Z"
    }
  ]
}
```

---

### Test 7: Refresh Token

**Endpoint**: `POST /api/v1/auth/refresh`

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN_HERE"
  }'
```

**Expected Response** (200 OK):
```json
{
  "success": true,
  "message": "Token Rotation perform successfully",
  "data": {
    "access_token": "new-access-token",
    "refresh_token": "new-refresh-token",
    "expires_in": 3600
  }
}
```

**Note**: Old refresh token is invalidated after successful rotation.

---

### Test 8: Logout

**Endpoint**: `POST /api/v1/auth/logout`

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN_HERE"
```

**Expected Response** (204 No Content):
No body returned, just HTTP status 204.

**After Logout**: Using the same access token should return 401 Unauthorized.

---

## Complete Test Sequence

Here's a complete bash script to test the entire flow:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"
EMAIL="test_$(date +%s)@example.com"  # Unique email
PASSWORD="SecurePass123!"

echo "=== 1. Testing Registration ==="
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\",
    \"first_name\": \"Test\",
    \"last_name\": \"User\",
    \"phone_number\": \"+1234567890\"
  }")

echo $REGISTER_RESPONSE | jq .

# Extract access token
ACCESS_TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.data.tokens.access_token')
REFRESH_TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.data.tokens.refresh_token')

echo -e "\n=== 2. Testing Get Profile ==="
curl -s -X GET $BASE_URL/api/v1/users/me \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

echo -e "\n=== 3. Testing Update Profile ==="
curl -s -X PATCH $BASE_URL/api/v1/users/me \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Updated",
    "default_currency": "EUR"
  }' | jq .

echo -e "\n=== 4. Testing Get Sessions ==="
curl -s -X GET $BASE_URL/api/v1/auth/sessions \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

echo -e "\n=== 5. Testing Token Refresh ==="
REFRESH_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")

echo $REFRESH_RESPONSE | jq .

# Extract new tokens
NEW_ACCESS_TOKEN=$(echo $REFRESH_RESPONSE | jq -r '.data.access_token')

echo -e "\n=== 6. Testing Login ==="
curl -s -X POST $BASE_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }" | jq .

echo -e "\n=== 7. Testing Logout ==="
curl -s -X POST $BASE_URL/api/v1/auth/logout \
  -H "Authorization: Bearer $NEW_ACCESS_TOKEN" -w "\nHTTP Status: %{http_code}\n"

echo -e "\n=== 8. Testing Access After Logout (Should Fail) ==="
curl -s -X GET $BASE_URL/api/v1/users/me \
  -H "Authorization: Bearer $NEW_ACCESS_TOKEN" | jq .

echo -e "\n=== All Tests Completed ==="
```

**To run**:
```bash
chmod +x test_auth.sh
./test_auth.sh
```

---

## Password Validation Tests

### Valid Passwords ✅
- `SecurePass123!`
- `MyP@ssw0rd`
- `Test1234!@#$`
- `Welcome2024#`

### Invalid Passwords ❌
- `weak` - Too short
- `password123` - No uppercase, no special char
- `PASSWORD123!` - No lowercase
- `PasswordOnly!` - No number
- `Pass123` - Too short, no special char

---

## Postman Collection

You can also import this into Postman:

### Environment Variables
```
base_url: http://localhost:8080
access_token: (will be set automatically)
refresh_token: (will be set automatically)
```

### Pre-request Script (for Login/Register)
```javascript
// No pre-request needed
```

### Test Script (for Login/Register)
```javascript
if (pm.response.code === 200 || pm.response.code === 201) {
    const response = pm.response.json();
    pm.environment.set("access_token", response.data.tokens.access_token);
    pm.environment.set("refresh_token", response.data.tokens.refresh_token);
}
```

### Authorization Header (for protected routes)
```
Type: Bearer Token
Token: {{access_token}}
```

---

## Database Verification

### Check Users Table
```sql
SELECT id, email, first_name, last_name, email_verified, is_active, 
       created_at, last_login_at 
FROM users;
```

### Check Sessions Table
```sql
SELECT id, user_id, ip_address, user_agent, expire_at, 
       created_at, revoked_at 
FROM user_sessions 
WHERE revoked_at IS NULL;
```

### Check Encrypted Phone Number
```sql
SELECT email, phone_number, length(phone_number) as encrypted_length 
FROM users;
```

**Note**: Phone number should be base64 encoded and longer than the original.

---

## Common Issues & Solutions

### Issue 1: "Failed to connect to database"
**Solution**: Make sure PostgreSQL is running
```bash
docker container start postgres_arch
```

### Issue 2: "User already exists"
**Solution**: Use a different email or check existing users
```sql
DELETE FROM users WHERE email = 'test@example.com';
```

### Issue 3: "Invalid or expired token"
**Solution**: 
- Check if token has expired
- Verify JWT_SECRET matches in .env
- Use refresh token to get new access token

### Issue 4: "Validation failed"
**Solution**: Check password requirements:
- Min 8 characters
- At least 1 uppercase
- At least 1 lowercase
- At least 1 number
- At least 1 special character

---

## Performance Testing

### Load Test with Apache Bench

```bash
# Test login endpoint (create users first)
ab -n 100 -c 10 -p login.json -T application/json \
  http://localhost:8080/api/v1/auth/login
```

**login.json**:
```json
{
  "email": "test@example.com",
  "password": "SecurePass123!"
}
```

### Expected Performance
- Login: < 200ms
- Register: < 300ms (due to encryption)
- Get Profile: < 50ms
- Update Profile: < 100ms

---

## Security Testing

### Test 1: SQL Injection Prevention
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com OR 1=1--",
    "password": "anything"
  }'
```
**Expected**: Should return "Invalid credentials", not SQL error.

### Test 2: XSS Prevention
```bash
curl -X PATCH http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "<script>alert(\"XSS\")</script>"
  }'
```
**Expected**: Should store as plain text, not execute.

### Test 3: Token Reuse After Logout
```bash
# 1. Logout
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $ACCESS_TOKEN"

# 2. Try to use same token
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```
**Expected**: 401 Unauthorized

---

## Next Steps

After successful testing:

1. ✅ Verify all tests pass
2. ✅ Check database for correct data
3. ✅ Review logs for any errors
4. 🚀 Move to next module (Groups/Expenses)

---

**Happy Testing! 🧪**

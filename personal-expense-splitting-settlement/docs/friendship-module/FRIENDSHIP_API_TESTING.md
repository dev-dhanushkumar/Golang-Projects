# Phase 1: Friendship API Testing Guide

## ✅ Phase 1 Complete!

All 7 friendship endpoints are now implemented and running.

---

## 🎯 Endpoints Implemented

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/friends/request` | Send a friend request |
| POST | `/api/v1/friends/:id/accept` | Accept a friend request |
| POST | `/api/v1/friends/:id/reject` | Reject a friend request |
| POST | `/api/v1/friends/:id/block` | Block a user |
| DELETE | `/api/v1/friends/:id` | Remove a friend |
| GET | `/api/v1/friends` | Get list of accepted friends |
| GET | `/api/v1/friends/pending` | Get pending requests (sent & received) |

---

## 📝 Testing Workflow

### Prerequisites
Server is running on `http://localhost:8080`

### Step 1: Create Two Test Users

**User 1 - Alice:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@test.com",
    "password": "Alice123!",
    "first_name": "Alice",
    "last_name": "Smith"
  }'
```

**User 2 - Bob:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "bob@test.com",
    "password": "Bob123!",
    "first_name": "Bob",
    "last_name": "Johnson"
  }'
```

---

### Step 2: Login and Get Tokens

**Alice Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@test.com",
    "password": "Alice123!"
  }'
```

**Save Alice's token:**
```bash
export ALICE_TOKEN="<alice_access_token>"
```

**Bob Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "bob@test.com",
    "password": "Bob123!"
  }'
```

**Save Bob's token:**
```bash
export BOB_TOKEN="<bob_access_token>"
```

---

### Step 3: Send Friend Request

**Alice sends friend request to Bob:**
```bash
curl -X POST http://localhost:8080/api/v1/friends/request \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -d '{
    "friend_email": "bob@test.com"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Friend request sent successfully",
  "data": null
}
```

---

### Step 4: Check Pending Requests

**Bob checks pending requests (received):**
```bash
curl -X GET http://localhost:8080/api/v1/friends/pending \
  -H "Authorization: Bearer $BOB_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Pending requests retrieved successfully",
  "data": {
    "sent": [],
    "received": [
      {
        "id": "friendship-uuid",
        "requester_id": "alice-uuid",
        "email": "alice@test.com",
        "first_name": "Alice",
        "last_name": "Smith",
        "profile_image": null,
        "created_at": "2026-01-03T12:15:00Z"
      }
    ]
  }
}
```

**Alice checks pending requests (sent):**
```bash
curl -X GET http://localhost:8080/api/v1/friends/pending \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

---

### Step 5: Accept Friend Request

**Bob accepts Alice's friend request:**

First, get the friendship ID from the pending requests response above.

```bash
export FRIENDSHIP_ID="<friendship-uuid-from-above>"

curl -X POST http://localhost:8080/api/v1/friends/$FRIENDSHIP_ID/accept \
  -H "Authorization: Bearer $BOB_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Friend request accepted successfully",
  "data": null
}
```

---

### Step 6: View Friends List

**Alice views her friends:**
```bash
curl -X GET http://localhost:8080/api/v1/friends \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Friends retrieved successfully",
  "data": {
    "friends": [
      {
        "id": "friendship-uuid",
        "friend_id": "bob-uuid",
        "friend_email": "bob@test.com",
        "friend_name": "Bob Johnson",
        "status": "accepted",
        "requested_by": "alice-uuid",
        "is_requester": true,
        "created_at": "2026-01-03T12:15:00Z",
        "updated_at": "2026-01-03T12:16:00Z"
      }
    ],
    "total": 1
  }
}
```

**Bob views his friends:**
```bash
curl -X GET http://localhost:8080/api/v1/friends \
  -H "Authorization: Bearer $BOB_TOKEN"
```

---

### Step 7: Additional Tests

#### Test Reject Friend Request

**Create new users (Charlie & Dave), then:**

```bash
# Charlie sends request to Dave
curl -X POST http://localhost:8080/api/v1/friends/request \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CHARLIE_TOKEN" \
  -d '{"friend_email": "dave@test.com"}'

# Dave rejects the request
curl -X POST http://localhost:8080/api/v1/friends/$FRIENDSHIP_ID/reject \
  -H "Authorization: Bearer $DAVE_TOKEN"
```

#### Test Block User

```bash
curl -X POST http://localhost:8080/api/v1/friends/$FRIENDSHIP_ID/block \
  -H "Authorization: Bearer $BOB_TOKEN"
```

#### Test Remove Friend

```bash
curl -X DELETE http://localhost:8080/api/v1/friends/$FRIENDSHIP_ID \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

---

## 🔍 Testing Edge Cases

### 1. Cannot send friend request to yourself
```bash
curl -X POST http://localhost:8080/api/v1/friends/request \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -d '{
    "friend_email": "alice@test.com"
  }'
```

**Expected:** Error message "cannot send friend request to yourself"

---

### 2. Cannot send duplicate friend request
```bash
# Send same request twice
curl -X POST http://localhost:8080/api/v1/friends/request \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -d '{"friend_email": "bob@test.com"}'
```

**Expected:** Error message "friendship or friend request already exists"

---

### 3. Cannot accept own friend request
```bash
# Alice tries to accept her own request
curl -X POST http://localhost:8080/api/v1/friends/$FRIENDSHIP_ID/accept \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

**Expected:** Error message "cannot accept your own friend request"

---

### 4. User not found
```bash
curl -X POST http://localhost:8080/api/v1/friends/request \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -d '{
    "friend_email": "nonexistent@test.com"
  }'
```

**Expected:** Error message "user not found"

---

## 🗄️ Database Verification

**Check friendships table:**
```sql
psql -h localhost -U postgres -d personal-ess

-- View all friendships
SELECT * FROM friendships;

-- View with user details
SELECT 
    f.id,
    f.status,
    u1.email as user_1_email,
    u2.email as user_2_email,
    req.email as requester_email
FROM friendships f
JOIN users u1 ON f.user_id_1 = u1.id
JOIN users u2 ON f.user_id_2 = u2.id
JOIN users req ON f.requested_by = req.id;
```

---

## ✅ Validation Checklist

- [x] Migration 000003 applied successfully
- [x] Friendships table created with proper constraints
- [x] All 7 endpoints registered and accessible
- [x] Send friend request works
- [x] Accept friend request works
- [x] Reject friend request works
- [x] Block user works
- [x] Remove friend works
- [x] Get friends list works
- [x] Get pending requests works
- [x] Edge cases handled properly
- [x] User authentication required for all endpoints
- [x] Proper error messages returned
- [x] Database constraints enforced (user_id_1 < user_id_2)

---

## 📊 Phase 1 Status

**Completed:** 7/7 endpoints ✅

**Files Created:**
- ✅ migrations/000003_create_friendships_table.up.sql
- ✅ migrations/000003_create_friendships_table.down.sql
- ✅ internal/models/friendship.go
- ✅ internal/dto/friendship_dto.go
- ✅ internal/repository/friendship_repository.go
- ✅ internal/services/friendship_service.go
- ✅ internal/handler/friendship_handler.go

**Files Updated:**
- ✅ internal/router/router.go (added friendship routes)
- ✅ cmd/api/main.go (initialized friendship components)
- ✅ internal/models/errors.go (added ErrUserNotFound)

---

## 🎉 Next Steps

**Ready for Phase 2: Groups Module**

When ready, we'll implement:
- Groups table (8 endpoints)
- Group members management
- Role-based permissions
- Group CRUD operations

Just say "Let's start Phase 2" when you're ready! 🚀

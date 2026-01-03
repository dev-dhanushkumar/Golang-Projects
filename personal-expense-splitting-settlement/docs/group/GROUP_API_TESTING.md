# Group API Testing Guide

Complete guide for testing all Group module endpoints.

## Prerequisites

1. **Server Running**: Ensure the server is running on `http://localhost:8080`
2. **Authentication**: You need a valid JWT access token
3. **Friends**: Some endpoints require users to be friends first

### Get Your Access Token

```bash
# Register a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "Pass@1234",
    "first_name": "John",
    "last_name": "Doe"
  }'

# Login to get token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "Pass@1234"
  }'

# Save the access_token from response
export TOKEN="your_access_token_here"
```

---

## 1. Create Group

**Endpoint**: `POST /api/v1/groups`

**Description**: Create a new group with the authenticated user as admin.

### Request

```bash
curl -X POST http://localhost:8080/api/v1/groups \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Weekend Trip 2026",
    "description": "Goa beach vacation with friends",
    "type": "trip",
    "image_url": "https://example.com/group-image.jpg",
    "member_ids": []
  }'
```

### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Group name (1-255 chars) |
| description | string | No | Group description (max 1000 chars) |
| type | string | Yes | Group type: `general`, `trip`, `home`, `couple`, `event`, `project`, `other` |
| image_url | string | No | Valid URL for group image |
| member_ids | array | No | Array of user UUIDs to add as members (must be friends) |

### Success Response (201)

```json
{
  "success": true,
  "message": "Group created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Weekend Trip 2026",
    "description": "Goa beach vacation with friends",
    "type": "trip",
    "image_url": "https://example.com/group-image.jpg",
    "created_by": "123e4567-e89b-12d3-a456-426614174000",
    "members": [
      {
        "id": "member-uuid",
        "user_id": "123e4567-e89b-12d3-a456-426614174000",
        "name": "John Doe",
        "email": "user@example.com",
        "role": "admin",
        "joined_at": "2026-01-03T10:30:00Z"
      }
    ],
    "created_at": "2026-01-03T10:30:00Z",
    "updated_at": "2026-01-03T10:30:00Z"
  }
}
```

### Error Responses

**400 - Invalid Request**
```json
{
  "success": false,
  "message": "Invalid request body"
}
```

**400 - Validation Error**
```json
{
  "success": false,
  "message": "Key: 'CreateGroupRequest.Type' Error:Field validation for 'Type' failed on the 'oneof' tag"
}
```

### Testing Scenarios

```bash
# 1. Create a simple group
curl -X POST http://localhost:8080/api/v1/groups \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Home Expenses",
    "type": "home"
  }'

# 2. Create group with description
curl -X POST http://localhost:8080/api/v1/groups \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Office Party",
    "description": "Annual office celebration expenses",
    "type": "event"
  }'

# 3. Invalid type (should fail)
curl -X POST http://localhost:8080/api/v1/groups \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Group",
    "type": "invalid_type"
  }'

# Save group ID from successful response
export GROUP_ID="your_group_id_here"
```

---

## 2. List User's Groups

**Endpoint**: `GET /api/v1/groups`

**Description**: Retrieve all groups where the authenticated user is a member.

### Request

```bash
curl -X GET http://localhost:8080/api/v1/groups \
  -H "Authorization: Bearer $TOKEN"
```

### Success Response (200)

```json
{
  "success": true,
  "message": "Groups retrieved successfully",
  "data": {
    "groups": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Weekend Trip 2026",
        "description": "Goa beach vacation with friends",
        "type": "trip",
        "image_url": "https://example.com/group-image.jpg",
        "created_by": "123e4567-e89b-12d3-a456-426614174000",
        "member_count": 3,
        "created_at": "2026-01-03T10:30:00Z",
        "updated_at": "2026-01-03T10:30:00Z"
      },
      {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "name": "Home Expenses",
        "description": "",
        "type": "home",
        "image_url": "",
        "created_by": "123e4567-e89b-12d3-a456-426614174000",
        "member_count": 2,
        "created_at": "2026-01-03T11:00:00Z",
        "updated_at": "2026-01-03T11:00:00Z"
      }
    ],
    "total": 2
  }
}
```

### Error Responses

**401 - Unauthorized**
```json
{
  "success": false,
  "message": "Unauthorized"
}
```

---

## 3. Get Group Details

**Endpoint**: `GET /api/v1/groups/:id`

**Description**: Retrieve detailed information about a specific group including all members.

### Request

```bash
curl -X GET http://localhost:8080/api/v1/groups/$GROUP_ID \
  -H "Authorization: Bearer $TOKEN"
```

### Success Response (200)

```json
{
  "success": true,
  "message": "Group retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Weekend Trip 2026",
    "description": "Goa beach vacation with friends",
    "type": "trip",
    "image_url": "https://example.com/group-image.jpg",
    "created_by": "123e4567-e89b-12d3-a456-426614174000",
    "members": [
      {
        "id": "member-1-uuid",
        "user_id": "123e4567-e89b-12d3-a456-426614174000",
        "name": "John Doe",
        "email": "john@example.com",
        "role": "admin",
        "joined_at": "2026-01-03T10:30:00Z"
      },
      {
        "id": "member-2-uuid",
        "user_id": "223e4567-e89b-12d3-a456-426614174001",
        "name": "Jane Smith",
        "email": "jane@example.com",
        "role": "member",
        "joined_at": "2026-01-03T10:35:00Z"
      }
    ],
    "created_at": "2026-01-03T10:30:00Z",
    "updated_at": "2026-01-03T10:30:00Z"
  }
}
```

### Error Responses

**400 - Invalid Group ID**
```json
{
  "success": false,
  "message": "Invalid group ID"
}
```

**403 - Not a Member**
```json
{
  "success": false,
  "message": "You are not a member of this group"
}
```

**404 - Group Not Found**
```json
{
  "success": false,
  "message": "group not found"
}
```

### Testing Scenarios

```bash
# 1. Get existing group details
curl -X GET http://localhost:8080/api/v1/groups/$GROUP_ID \
  -H "Authorization: Bearer $TOKEN"

# 2. Invalid group ID (should fail)
curl -X GET http://localhost:8080/api/v1/groups/invalid-uuid \
  -H "Authorization: Bearer $TOKEN"

# 3. Group you're not a member of (should fail with 403)
export OTHER_GROUP_ID="some-other-group-uuid"
curl -X GET http://localhost:8080/api/v1/groups/$OTHER_GROUP_ID \
  -H "Authorization: Bearer $TOKEN"
```

---

## 4. Update Group

**Endpoint**: `PATCH /api/v1/groups/:id`

**Description**: Update group information. Only admins can update.

### Request

```bash
curl -X PATCH http://localhost:8080/api/v1/groups/$GROUP_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Trip Name",
    "description": "Updated description for the trip"
  }'
```

### Request Body (All fields optional)

| Field | Type | Description |
|-------|------|-------------|
| name | string | Updated group name (1-255 chars) |
| description | string | Updated description (max 1000 chars) |
| type | string | Updated group type |
| image_url | string | Updated image URL |

### Success Response (200)

```json
{
  "success": true,
  "message": "Group updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Updated Trip Name",
    "description": "Updated description for the trip",
    "type": "trip",
    "image_url": "https://example.com/group-image.jpg",
    "created_by": "123e4567-e89b-12d3-a456-426614174000",
    "members": [...],
    "created_at": "2026-01-03T10:30:00Z",
    "updated_at": "2026-01-03T12:00:00Z"
  }
}
```

### Error Responses

**403 - Not Admin**
```json
{
  "success": false,
  "message": "only admins can update group information"
}
```

### Testing Scenarios

```bash
# 1. Update only name
curl -X PATCH http://localhost:8080/api/v1/groups/$GROUP_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Group Name"
  }'

# 2. Update multiple fields
curl -X PATCH http://localhost:8080/api/v1/groups/$GROUP_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bali Trip 2026",
    "description": "Beach vacation in Bali",
    "type": "trip",
    "image_url": "https://example.com/bali.jpg"
  }'

# 3. Non-admin trying to update (should fail with 403)
# First, get a token for a non-admin member
export NON_ADMIN_TOKEN="token_of_non_admin_member"
curl -X PATCH http://localhost:8080/api/v1/groups/$GROUP_ID \
  -H "Authorization: Bearer $NON_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Should Fail"
  }'
```

---

## 5. Delete Group

**Endpoint**: `DELETE /api/v1/groups/:id`

**Description**: Soft delete a group. Only admins can delete.

### Request

```bash
curl -X DELETE http://localhost:8080/api/v1/groups/$GROUP_ID \
  -H "Authorization: Bearer $TOKEN"
```

### Success Response (200)

```json
{
  "success": true,
  "message": "Group deleted successfully"
}
```

### Error Responses

**403 - Not Admin**
```json
{
  "success": false,
  "message": "only admins can delete groups"
}
```

### Testing Scenarios

```bash
# 1. Admin deleting group
curl -X DELETE http://localhost:8080/api/v1/groups/$GROUP_ID \
  -H "Authorization: Bearer $TOKEN"

# 2. Non-admin trying to delete (should fail)
curl -X DELETE http://localhost:8080/api/v1/groups/$GROUP_ID \
  -H "Authorization: Bearer $NON_ADMIN_TOKEN"

# 3. Verify group is deleted (should return 404)
curl -X GET http://localhost:8080/api/v1/groups/$GROUP_ID \
  -H "Authorization: Bearer $TOKEN"
```

---

## 6. Add Member to Group

**Endpoint**: `POST /api/v1/groups/:id/members`

**Description**: Add a new member to the group. Only admins can add members. Users must be friends to be added.

### Request

```bash
curl -X POST http://localhost:8080/api/v1/groups/$GROUP_ID/members \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "223e4567-e89b-12d3-a456-426614174001",
    "role": "member"
  }'
```

### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| user_id | UUID | Yes | User ID to add as member |
| role | string | No | Role: `admin` or `member` (default: `member`) |

### Success Response (201)

```json
{
  "success": true,
  "message": "Member added successfully",
  "data": {
    "id": "member-uuid",
    "user_id": "223e4567-e89b-12d3-a456-426614174001",
    "name": "Jane Smith",
    "email": "jane@example.com",
    "role": "member",
    "joined_at": "2026-01-03T12:00:00Z"
  }
}
```

### Error Responses

**400 - Not Friends**
```json
{
  "success": false,
  "message": "can only add friends to the group"
}
```

**403 - Not Admin**
```json
{
  "success": false,
  "message": "only admins can add members"
}
```

**404 - User Not Found**
```json
{
  "success": false,
  "message": "user not found"
}
```

### Testing Scenarios

```bash
# First, create friendship
export FRIEND_ID="223e4567-e89b-12d3-a456-426614174001"

# Send friend request (from another account)
curl -X POST http://localhost:8080/api/v1/friends/request \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "friend@example.com"
  }'

# Accept friend request (from friend's account)
# ... (see friendship API testing)

# 1. Add friend as member
curl -X POST http://localhost:8080/api/v1/groups/$GROUP_ID/members \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "'$FRIEND_ID'"
  }'

# 2. Add friend as admin
curl -X POST http://localhost:8080/api/v1/groups/$GROUP_ID/members \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "'$FRIEND_ID'",
    "role": "admin"
  }'

# 3. Try to add non-friend (should fail)
export NON_FRIEND_ID="some-random-user-uuid"
curl -X POST http://localhost:8080/api/v1/groups/$GROUP_ID/members \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "'$NON_FRIEND_ID'"
  }'

# 4. Non-admin trying to add member (should fail)
curl -X POST http://localhost:8080/api/v1/groups/$GROUP_ID/members \
  -H "Authorization: Bearer $NON_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "'$FRIEND_ID'"
  }'
```

---

## 7. Remove Member from Group

**Endpoint**: `DELETE /api/v1/groups/:id/members/:user_id`

**Description**: Remove a member from the group. Admins can remove others, or users can remove themselves.

### Request

```bash
# Admin removing a member
curl -X DELETE http://localhost:8080/api/v1/groups/$GROUP_ID/members/$MEMBER_ID \
  -H "Authorization: Bearer $TOKEN"
```

### Success Response (200)

```json
{
  "success": true,
  "message": "Member removed successfully"
}
```

### Error Responses

**400 - Last Admin**
```json
{
  "success": false,
  "message": "cannot remove the last admin from the group"
}
```

**403 - Not Authorized**
```json
{
  "success": false,
  "message": "only admins can remove other members"
}
```

### Testing Scenarios

```bash
export MEMBER_TO_REMOVE="member-user-uuid"

# 1. Admin removing regular member
curl -X DELETE http://localhost:8080/api/v1/groups/$GROUP_ID/members/$MEMBER_TO_REMOVE \
  -H "Authorization: Bearer $TOKEN"

# 2. User removing themselves
curl -X DELETE http://localhost:8080/api/v1/groups/$GROUP_ID/members/$MY_USER_ID \
  -H "Authorization: Bearer $TOKEN"

# 3. Non-admin trying to remove someone else (should fail)
curl -X DELETE http://localhost:8080/api/v1/groups/$GROUP_ID/members/$MEMBER_TO_REMOVE \
  -H "Authorization: Bearer $NON_ADMIN_TOKEN"

# 4. Try to remove the last admin (should fail)
# First ensure there's only one admin
curl -X DELETE http://localhost:8080/api/v1/groups/$GROUP_ID/members/$ONLY_ADMIN_ID \
  -H "Authorization: Bearer $TOKEN"
```

---

## 8. Update Member Role

**Endpoint**: `PATCH /api/v1/groups/:id/members/:user_id`

**Description**: Update a member's role (promote to admin or demote to member). Only admins can change roles.

### Request

```bash
curl -X PATCH http://localhost:8080/api/v1/groups/$GROUP_ID/members/$MEMBER_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "admin"
  }'
```

### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| role | string | Yes | New role: `admin` or `member` |

### Success Response (200)

```json
{
  "success": true,
  "message": "Member role updated successfully"
}
```

### Error Responses

**400 - Last Admin**
```json
{
  "success": false,
  "message": "cannot demote the last admin"
}
```

**403 - Not Admin**
```json
{
  "success": false,
  "message": "only admins can change member roles"
}
```

**403 - Cannot Change Own Role**
```json
{
  "success": false,
  "message": "cannot change your own role"
}
```

### Testing Scenarios

```bash
export MEMBER_TO_PROMOTE="member-user-uuid"

# 1. Promote member to admin
curl -X PATCH http://localhost:8080/api/v1/groups/$GROUP_ID/members/$MEMBER_TO_PROMOTE \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "admin"
  }'

# 2. Demote admin to member
curl -X PATCH http://localhost:8080/api/v1/groups/$GROUP_ID/members/$ADMIN_TO_DEMOTE \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "member"
  }'

# 3. Try to change your own role (should fail)
curl -X PATCH http://localhost:8080/api/v1/groups/$GROUP_ID/members/$MY_USER_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "member"
  }'

# 4. Non-admin trying to change roles (should fail)
curl -X PATCH http://localhost:8080/api/v1/groups/$GROUP_ID/members/$MEMBER_ID \
  -H "Authorization: Bearer $NON_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "admin"
  }'

# 5. Try to demote the last admin (should fail)
curl -X PATCH http://localhost:8080/api/v1/groups/$GROUP_ID/members/$ONLY_ADMIN_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "member"
  }'
```

---

## Complete Testing Flow

Here's a complete workflow to test all group functionality:

```bash
# 1. Setup - Register and login users
export TOKEN1="token_for_user1"
export TOKEN2="token_for_user2"
export USER2_ID="user2-uuid"

# 2. Create friendship between users
curl -X POST http://localhost:8080/api/v1/friends/request \
  -H "Authorization: Bearer $TOKEN1" \
  -H "Content-Type: application/json" \
  -d '{"email": "user2@example.com"}'

# Accept friend request
curl -X POST http://localhost:8080/api/v1/friends/$FRIENDSHIP_ID/accept \
  -H "Authorization: Bearer $TOKEN2"

# 3. Create group (User1 is admin)
GROUP_RESPONSE=$(curl -X POST http://localhost:8080/api/v1/groups \
  -H "Authorization: Bearer $TOKEN1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Group",
    "description": "Testing group functionality",
    "type": "general"
  }')

export GROUP_ID=$(echo $GROUP_RESPONSE | jq -r '.data.id')

# 4. List groups
curl -X GET http://localhost:8080/api/v1/groups \
  -H "Authorization: Bearer $TOKEN1"

# 5. Get group details
curl -X GET http://localhost:8080/api/v1/groups/$GROUP_ID \
  -H "Authorization: Bearer $TOKEN1"

# 6. Update group
curl -X PATCH http://localhost:8080/api/v1/groups/$GROUP_ID \
  -H "Authorization: Bearer $TOKEN1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Test Group",
    "description": "Updated description"
  }'

# 7. Add member to group
curl -X POST http://localhost:8080/api/v1/groups/$GROUP_ID/members \
  -H "Authorization: Bearer $TOKEN1" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "'$USER2_ID'",
    "role": "member"
  }'

# 8. User2 views group (should succeed as they're now a member)
curl -X GET http://localhost:8080/api/v1/groups/$GROUP_ID \
  -H "Authorization: Bearer $TOKEN2"

# 9. Promote User2 to admin
curl -X PATCH http://localhost:8080/api/v1/groups/$GROUP_ID/members/$USER2_ID \
  -H "Authorization: Bearer $TOKEN1" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "admin"
  }'

# 10. User2 can now update group (as admin)
curl -X PATCH http://localhost:8080/api/v1/groups/$GROUP_ID \
  -H "Authorization: Bearer $TOKEN2" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated by User2"
  }'

# 11. User2 removes themselves from group
curl -X DELETE http://localhost:8080/api/v1/groups/$GROUP_ID/members/$USER2_ID \
  -H "Authorization: Bearer $TOKEN2"

# 12. Delete group
curl -X DELETE http://localhost:8080/api/v1/groups/$GROUP_ID \
  -H "Authorization: Bearer $TOKEN1"
```

---

## Monitoring Logs

To monitor request logs while testing:

```bash
# Watch logs in real-time
tail -f logs/app.log | jq .

# Filter only group-related requests
tail -f logs/app.log | jq 'select(.path | contains("/groups"))'

# Watch for errors
tail -f logs/app.log | jq 'select(.level == "error" or .level == "warn")'
```

---

## Common Validation Errors

### Invalid Group Type
```bash
curl -X POST http://localhost:8080/api/v1/groups \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "type": "invalid"
  }'
# Error: Field validation for 'Type' failed on the 'oneof' tag
```

### Invalid UUID Format
```bash
curl -X GET http://localhost:8080/api/v1/groups/not-a-uuid \
  -H "Authorization: Bearer $TOKEN"
# Error: Invalid group ID
```

### Missing Required Fields
```bash
curl -X POST http://localhost:8080/api/v1/groups \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Missing name and type"
  }'
# Error: Field validation errors
```

---

## Next Steps

- **Phase 3**: Expenses Module - Add expense tracking to groups
- **Phase 4**: Balance & Settlement - Calculate balances and settle debts
- **Phase 5**: Notifications - Notify members of group activities

---

## Notes

1. **Admin Protection**: Database triggers ensure at least one admin always exists in a group
2. **Friendship Requirement**: Can only add users who are accepted friends
3. **Member Count**: Cached in the group list response for performance
4. **Soft Delete**: Groups are soft-deleted (deleted_at timestamp) for data integrity
5. **Audit Trail**: All membership changes tracked with joined_at and left_at timestamps

---

## Troubleshooting

### "User not authenticated" error
- Ensure your token is valid and not expired
- Re-login to get a fresh token

### "You are not a member of this group" error
- Verify you're accessing a group you belong to
- Check group ID is correct

### "only admins can..." errors
- Verify your role in the group
- Use an admin token for admin-only operations

### "can only add friends to the group" error
- Ensure friendship exists and is accepted
- Check friendship status with GET /api/v1/friends

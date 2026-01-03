# Expense Module API Testing Documentation

This guide provides comprehensive testing instructions for all expense-related endpoints in the Personal Expense Splitting and Settlement application.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Authentication Setup](#authentication-setup)
- [Expense Endpoints](#expense-endpoints)
  - [1. Create Expense](#1-create-expense)
  - [2. Get Expense Details](#2-get-expense-details)
  - [3. Get User Expenses](#3-get-user-expenses)
  - [4. Get Group Expenses](#4-get-group-expenses)
  - [5. Get Expenses with Filters](#5-get-expenses-with-filters)
  - [6. Update Expense](#6-update-expense)
  - [7. Delete Expense](#7-delete-expense)
- [Split Methods Explained](#split-methods-explained)
- [Complete Testing Flow](#complete-testing-flow)
- [Error Scenarios](#error-scenarios)
- [Monitoring Logs](#monitoring-logs)

## Prerequisites

1. Server running on `http://localhost:8080`
2. At least 2 registered users (we'll call them User1 and User2)
3. Users should be friends (use friendship endpoints from Phase 1)
4. At least one group created with both users as members (use group endpoints from Phase 2)

## Authentication Setup

First, register and login users to get JWT tokens:

```bash
# Register User 1
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user1@example.com",
    "password": "Password123!",
    "first_name": "John",
    "last_name": "Doe"
  }'

# Register User 2
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user2@example.com",
    "password": "Password123!",
    "first_name": "Jane",
    "last_name": "Smith"
  }'

# Login as User 1
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user1@example.com",
    "password": "Password123!"
  }'

# Save the access_token from the response
export TOKEN_USER1="<access_token_from_response>"

# Login as User 2
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user2@example.com",
    "password": "Password123!"
  }'

export TOKEN_USER2="<access_token_from_response>"
```

## Expense Endpoints

### 1. Create Expense

Creates a new expense with participants. Supports 4 split methods: equal, exact, percentage, and shares.

#### Example 1: Equal Split (Most Common)

```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER1" \
  -d '{
    "description": "Dinner at restaurant",
    "amount": 150.00,
    "category": "food",
    "date": "2026-01-03",
    "split_method": "equal",
    "participants": [
      {
        "user_id": "<USER1_ID>",
        "paid_amount": 150.00
      },
      {
        "user_id": "<USER2_ID>",
        "paid_amount": 0.00
      }
    ]
  }'
```

**Response:**
```json
{
  "status": "success",
  "message": "Expense created successfully",
  "data": {
    "id": "uuid",
    "description": "Dinner at restaurant",
    "amount": 150.00,
    "category": "food",
    "date": "2026-01-03",
    "receipt_url": "",
    "created_by": "uuid",
    "group_id": null,
    "creator_name": "John Doe",
    "group_name": "",
    "participants": [
      {
        "id": "uuid",
        "user_id": "uuid",
        "user_name": "John Doe",
        "paid_amount": 150.00,
        "owed_amount": 75.00,
        "net_amount": 75.00,
        "is_settled": false
      },
      {
        "id": "uuid",
        "user_id": "uuid",
        "user_name": "Jane Smith",
        "paid_amount": 0.00,
        "owed_amount": 75.00,
        "net_amount": -75.00,
        "is_settled": false
      }
    ],
    "created_at": "2026-01-03T14:52:37Z",
    "updated_at": "2026-01-03T14:52:37Z"
  }
}
```

#### Example 2: Exact Amounts Split

```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER1" \
  -d '{
    "description": "Grocery shopping",
    "amount": 100.00,
    "category": "shopping",
    "split_method": "exact",
    "participants": [
      {
        "user_id": "<USER1_ID>",
        "paid_amount": 100.00,
        "owed_amount": 60.00
      },
      {
        "user_id": "<USER2_ID>",
        "paid_amount": 0.00,
        "owed_amount": 40.00
      }
    ]
  }'
```

#### Example 3: Percentage Split

```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER1" \
  -d '{
    "description": "Movie tickets",
    "amount": 200.00,
    "category": "entertainment",
    "split_method": "percentage",
    "participants": [
      {
        "user_id": "<USER1_ID>",
        "paid_amount": 200.00,
        "percentage": 60.0
      },
      {
        "user_id": "<USER2_ID>",
        "paid_amount": 0.00,
        "percentage": 40.0
      }
    ]
  }'
```

#### Example 4: Shares Split

```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER1" \
  -d '{
    "description": "Taxi ride",
    "amount": 300.00,
    "category": "transport",
    "split_method": "shares",
    "participants": [
      {
        "user_id": "<USER1_ID>",
        "paid_amount": 300.00,
        "shares": 2
      },
      {
        "user_id": "<USER2_ID>",
        "paid_amount": 0.00,
        "shares": 1
      }
    ]
  }'
```

#### Example 5: Group Expense

```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER1" \
  -d '{
    "description": "Team lunch",
    "amount": 500.00,
    "category": "food",
    "group_id": "<GROUP_ID>",
    "split_method": "equal",
    "participants": [
      {
        "user_id": "<USER1_ID>",
        "paid_amount": 500.00
      },
      {
        "user_id": "<USER2_ID>",
        "paid_amount": 0.00
      }
    ]
  }'
```

### 2. Get Expense Details

Retrieves detailed information about a specific expense including all participants.

```bash
curl -X GET http://localhost:8080/api/v1/expenses/<EXPENSE_ID> \
  -H "Authorization: Bearer $TOKEN_USER1"
```

**Response:**
```json
{
  "status": "success",
  "message": "Expense retrieved successfully",
  "data": {
    "id": "uuid",
    "description": "Dinner at restaurant",
    "amount": 150.00,
    "category": "food",
    "participants": [
      {
        "id": "uuid",
        "user_name": "John Doe",
        "paid_amount": 150.00,
        "owed_amount": 75.00,
        "net_amount": 75.00,
        "is_settled": false
      }
    ]
  }
}
```

### 3. Get User Expenses

Retrieves all expenses where the authenticated user is a participant.

```bash
# Get all expenses for the user
curl -X GET http://localhost:8080/api/v1/expenses \
  -H "Authorization: Bearer $TOKEN_USER1"

# With pagination
curl -X GET "http://localhost:8080/api/v1/expenses?limit=10&offset=0" \
  -H "Authorization: Bearer $TOKEN_USER1"
```

**Response:**
```json
{
  "status": "success",
  "message": "Expenses retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "description": "Dinner at restaurant",
      "amount": 150.00,
      "category": "food",
      "creator_name": "John Doe",
      "group_name": "",
      "date": "2026-01-03",
      "created_at": "2026-01-03T14:52:37Z"
    }
  ]
}
```

### 4. Get Group Expenses

Retrieves all expenses for a specific group. User must be a member of the group.

```bash
# Get all expenses for a group
curl -X GET http://localhost:8080/api/v1/groups/<GROUP_ID>/expenses \
  -H "Authorization: Bearer $TOKEN_USER1"

# With pagination
curl -X GET "http://localhost:8080/api/v1/groups/<GROUP_ID>/expenses?limit=10&offset=0" \
  -H "Authorization: Bearer $TOKEN_USER1"
```

**Response:**
```json
{
  "status": "success",
  "message": "Expenses retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "description": "Team lunch",
      "amount": 500.00,
      "category": "food",
      "group_name": "Trip to Goa",
      "date": "2026-01-03"
    }
  ]
}
```

### 5. Get Expenses with Filters

Retrieves expenses with various filters: group, category, date range.

```bash
# Filter by category
curl -X GET "http://localhost:8080/api/v1/expenses/filter?category=food&limit=20" \
  -H "Authorization: Bearer $TOKEN_USER1"

# Filter by group
curl -X GET "http://localhost:8080/api/v1/expenses/filter?group_id=<GROUP_ID>" \
  -H "Authorization: Bearer $TOKEN_USER1"

# Filter by date range
curl -X GET "http://localhost:8080/api/v1/expenses/filter?start_date=2026-01-01&end_date=2026-01-31" \
  -H "Authorization: Bearer $TOKEN_USER1"

# Multiple filters
curl -X GET "http://localhost:8080/api/v1/expenses/filter?category=food&start_date=2026-01-01&end_date=2026-01-31&limit=10&offset=0" \
  -H "Authorization: Bearer $TOKEN_USER1"
```

**Response:**
```json
{
  "status": "success",
  "message": "Expenses retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "description": "Dinner at restaurant",
      "amount": 150.00,
      "category": "food",
      "date": "2026-01-03"
    }
  ]
}
```

### 6. Update Expense

Updates expense details. Only the creator can update, and amount cannot be changed after creation.

```bash
curl -X PATCH http://localhost:8080/api/v1/expenses/<EXPENSE_ID> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER1" \
  -d '{
    "description": "Dinner at Italian restaurant",
    "category": "food",
    "receipt_url": "https://example.com/receipt.pdf"
  }'
```

**Response:**
```json
{
  "status": "success",
  "message": "Expense updated successfully",
  "data": {
    "id": "uuid",
    "description": "Dinner at Italian restaurant",
    "amount": 150.00,
    "category": "food",
    "receipt_url": "https://example.com/receipt.pdf"
  }
}
```

### 7. Delete Expense

Deletes an expense. Only the creator can delete, and cannot delete if any settlements have been made.

```bash
curl -X DELETE http://localhost:8080/api/v1/expenses/<EXPENSE_ID> \
  -H "Authorization: Bearer $TOKEN_USER1"
```

**Response:**
```json
{
  "status": "success",
  "message": "Expense deleted successfully",
  "data": null
}
```

## Split Methods Explained

### 1. Equal Split (`split_method: "equal"`)
- Amount is divided equally among all participants
- No need to specify `owed_amount`, `percentage`, or `shares`
- Example: ₹150 split between 2 people = ₹75 each

### 2. Exact Amounts (`split_method: "exact"`)
- Specify exact `owed_amount` for each participant
- Total owed amounts must equal total expense amount
- Example: ₹100 expense where User1 owes ₹60, User2 owes ₹40

### 3. Percentage Split (`split_method: "percentage"`)
- Specify `percentage` for each participant
- Total percentages must equal 100%
- Example: ₹200 expense with 60%-40% split = ₹120 and ₹80

### 4. Shares Split (`split_method: "shares"`)
- Specify `shares` for each participant
- Amount divided proportionally based on shares
- Example: ₹300 with 2:1 share ratio = ₹200 and ₹100

## Complete Testing Flow

### Step 1: Setup (Users and Group)
```bash
# 1. Register two users
# 2. Login both users and save tokens
# 3. Send friend request from User1 to User2
# 4. Accept friend request from User2
# 5. Create a group with both users as members
# Save USER1_ID, USER2_ID, and GROUP_ID
```

### Step 2: Create Various Expenses
```bash
# Create equal split expense
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER1" \
  -d '{
    "description": "Restaurant bill",
    "amount": 200.00,
    "category": "food",
    "split_method": "equal",
    "participants": [
      {"user_id": "<USER1_ID>", "paid_amount": 200.00},
      {"user_id": "<USER2_ID>", "paid_amount": 0.00}
    ]
  }'

# Save EXPENSE_ID from response
```

### Step 3: View Expenses
```bash
# Get expense details
curl -X GET http://localhost:8080/api/v1/expenses/<EXPENSE_ID> \
  -H "Authorization: Bearer $TOKEN_USER1"

# Get all user expenses
curl -X GET http://localhost:8080/api/v1/expenses \
  -H "Authorization: Bearer $TOKEN_USER1"

# Get group expenses
curl -X GET http://localhost:8080/api/v1/groups/<GROUP_ID>/expenses \
  -H "Authorization: Bearer $TOKEN_USER1"
```

### Step 4: Filter Expenses
```bash
# Filter by category
curl -X GET "http://localhost:8080/api/v1/expenses/filter?category=food" \
  -H "Authorization: Bearer $TOKEN_USER1"

# Filter by date range
curl -X GET "http://localhost:8080/api/v1/expenses/filter?start_date=2026-01-01&end_date=2026-01-31" \
  -H "Authorization: Bearer $TOKEN_USER1"
```

### Step 5: Update Expense
```bash
curl -X PATCH http://localhost:8080/api/v1/expenses/<EXPENSE_ID> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER1" \
  -d '{
    "description": "Updated description"
  }'
```

### Step 6: Delete Expense
```bash
curl -X DELETE http://localhost:8080/api/v1/expenses/<EXPENSE_ID> \
  -H "Authorization: Bearer $TOKEN_USER1"
```

## Error Scenarios

### 1. Unauthorized Access
```bash
# No token
curl -X GET http://localhost:8080/api/v1/expenses
# Expected: 401 Unauthorized
```

### 2. Invalid Split Method
```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER1" \
  -d '{
    "description": "Test",
    "amount": 100.00,
    "split_method": "equal",
    "participants": [
      {"user_id": "<USER1_ID>", "paid_amount": 50.00},
      {"user_id": "<USER2_ID>", "paid_amount": 30.00}
    ]
  }'
# Expected: 400 Bad Request - total paid amount must equal expense amount
```

### 3. Non-participant Viewing Expense
```bash
# User2 tries to view User1's personal expense
curl -X GET http://localhost:8080/api/v1/expenses/<EXPENSE_ID> \
  -H "Authorization: Bearer $TOKEN_USER2"
# Expected: 404 Not Found - you are not a participant
```

### 4. Non-creator Updating Expense
```bash
curl -X PATCH http://localhost:8080/api/v1/expenses/<EXPENSE_ID> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER2" \
  -d '{"description": "Updated"}'
# Expected: 403 Forbidden - only creator can update
```

### 5. Invalid Category
```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER1" \
  -d '{
    "description": "Test",
    "amount": 100.00,
    "category": "invalid_category",
    "split_method": "equal",
    "participants": [...]
  }'
# Expected: 400 Bad Request - invalid expense category
```

## Monitoring Logs

Watch server logs in real-time to track expense operations:

```bash
# In a separate terminal
tail -f logs/app.log | grep -i "expense"

# Or view JSON logs
tail -f logs/app.json | jq 'select(.msg | contains("expense"))'
```

## Available Categories

- `general` - General expenses
- `food` - Food and dining
- `transport` - Transportation
- `entertainment` - Entertainment and recreation
- `utilities` - Utilities (electricity, water, etc.)
- `shopping` - Shopping
- `healthcare` - Healthcare and medical
- `education` - Education
- `travel` - Travel and tourism
- `other` - Other expenses

## Validation Rules

1. **Description**: Required, max 500 characters, cannot be empty
2. **Amount**: Required, must be positive
3. **Category**: Must be one of the valid categories
4. **Date**: Cannot be in the future
5. **Split Method**: Must be one of: equal, exact, percentage, shares
6. **Participants**: At least one participant required
7. **Paid Amounts**: Total must equal expense amount
8. **Owed Amounts** (exact method): Total must equal expense amount
9. **Percentages** (percentage method): Total must equal 100%
10. **Shares** (shares method): Must be positive integers
11. **Group Expenses**: All participants must be group members
12. **Amount Update**: Cannot update amount after creation
13. **Delete**: Cannot delete if any participant has settled

## Notes

- All amounts are stored as DECIMAL(12,2) for precision
- Timestamps are in UTC
- Net amount = Paid amount - Owed amount
- Positive net amount = User is owed money (creditor)
- Negative net amount = User owes money (debtor)
- Settlement status is tracked per participant
- Expenses support soft delete (deleted_at field)
- Migrations 000006 and 000007 handle expense tables

## Troubleshooting

1. **Migration not applied**: Check logs for migration errors
2. **Foreign key errors**: Ensure users and groups exist
3. **Split calculation errors**: Verify total paid = total amount
4. **Permission errors**: Verify user is participant/creator
5. **Group expense errors**: Verify all participants are group members

# Phase 4: Balance & Settlement API Testing Guide

## ✅ Phase 4 Complete!

All 13 balance and settlement endpoints are now implemented and running.

---

## 🎯 Endpoints Implemented

### Settlement Endpoints (8)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/settlements` | Create a new settlement |
| GET | `/api/v1/settlements` | Get user's settlement history |
| GET | `/api/v1/settlements/between?user_id={uuid}` | Get settlements between two users |
| GET | `/api/v1/settlements/:id` | Get specific settlement details |
| PATCH | `/api/v1/settlements/:id` | Update settlement (payer only, before confirmation) |
| PATCH | `/api/v1/settlements/:id/confirm` | Confirm settlement (payee only) |
| DELETE | `/api/v1/settlements/:id` | Delete settlement (payer only, before confirmation) |
| GET | `/api/v1/groups/:id/settlements` | Get all settlements for a group |

### Balance & Suggestion Endpoints (5)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users/me/balance-summary` | Get total balance summary |
| GET | `/api/v1/users/me/balances` | Get detailed balances with each person |
| GET | `/api/v1/groups/:id/balances` | Get balance breakdown for group members |
| GET | `/api/v1/settlements/suggestions` | Get smart settlement suggestions |
| GET | `/api/v1/groups/:id/settlement-suggestions` | Get group settlement suggestions |

---

## 📝 Testing Workflow

### Prerequisites
- Server is running on `http://localhost:8080`
- You have completed Phase 3 (Expenses) testing
- You have at least 2 users with expenses between them

### Quick Setup
If you need test users and expenses:

**Create and login users (Alice & Bob):**
```bash
# Register Alice
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@test.com",
    "password": "Alice123!",
    "first_name": "Alice",
    "last_name": "Smith"
  }'

# Login Alice
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@test.com",
    "password": "Alice123!"
  }'

# Save token
export ALICE_TOKEN="<alice_access_token>"
export ALICE_ID="<alice_user_id>"

# Register and login Bob
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "bob@test.com",
    "password": "Bob123!",
    "first_name": "Bob",
    "last_name": "Johnson"
  }'

curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "bob@test.com",
    "password": "Bob123!"
  }'

export BOB_TOKEN="<bob_access_token>"
export BOB_ID="<bob_user_id>"
```

**Create some expenses (for balance calculation):**
```bash
# Alice creates an expense where she paid for Bob
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -d '{
    "description": "Lunch together",
    "amount": 500.00,
    "category": "food",
    "date": "2026-01-03",
    "participants": [
      {
        "user_id": "'$ALICE_ID'",
        "paid_amount": 500.00,
        "owed_amount": 250.00
      },
      {
        "user_id": "'$BOB_ID'",
        "paid_amount": 0,
        "owed_amount": 250.00
      }
    ]
  }'
```

---

## 📋 Part 1: Balance Checking

### Step 1: Check Balance Summary

**Alice checks her balance summary:**
```bash
curl -X GET http://localhost:8080/api/v1/users/me/balance-summary \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Balance summary retrieved successfully",
  "data": {
    "total_owed": 250.00,
    "total_owing": 0,
    "net_balance": 250.00,
    "balances": [
      {
        "user_id": "bob-uuid",
        "user_name": "Bob Johnson",
        "user_email": "bob@test.com",
        "amount": 250.00
      }
    ],
    "last_updated": "2026-01-03T12:00:00Z"
  }
}
```

**Note:** 
- `amount > 0`: They owe you
- `amount < 0`: You owe them

---

### Step 2: Check Detailed Balances

**Alice checks detailed balances:**
```bash
curl -X GET http://localhost:8080/api/v1/users/me/balances \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

**Bob checks his balances:**
```bash
curl -X GET http://localhost:8080/api/v1/users/me/balances \
  -H "Authorization: Bearer $BOB_TOKEN"
```

**Expected for Bob:**
```json
{
  "success": true,
  "message": "Balances retrieved successfully",
  "data": {
    "total_owed": 0,
    "total_owing": 250.00,
    "net_balance": -250.00,
    "balances": [
      {
        "user_id": "alice-uuid",
        "user_name": "Alice Smith",
        "user_email": "alice@test.com",
        "amount": -250.00
      }
    ],
    "last_updated": "2026-01-03T12:00:00Z"
  }
}
```

---

### Step 3: Get Settlement Suggestions

**Alice gets settlement suggestions:**
```bash
curl -X GET http://localhost:8080/api/v1/settlements/suggestions \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Settlement suggestions generated successfully",
  "data": {
    "suggestions": [
      {
        "from": "bob-uuid",
        "from_name": "Bob Johnson",
        "to": "alice-uuid",
        "to_name": "Alice Smith",
        "amount": 250.00
      }
    ],
    "total_amount": 250.00,
    "message": "Here are suggested settlements to simplify your balances"
  }
}
```

---

## 💰 Part 2: Settlement Management

### Step 4: Create Settlement

**Bob creates a settlement to pay Alice:**
```bash
curl -X POST http://localhost:8080/api/v1/settlements \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $BOB_TOKEN" \
  -d '{
    "payee_id": "'$ALICE_ID'",
    "amount": 250.00,
    "payment_method": "upi",
    "notes": "Payment for lunch"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Settlement created successfully",
  "data": {
    "id": "settlement-uuid",
    "payer_id": "bob-uuid",
    "payer_name": "Bob Johnson",
    "payee_id": "alice-uuid",
    "payee_name": "Alice Smith",
    "amount": 250.00,
    "payment_method": "upi",
    "notes": "Payment for lunch",
    "is_confirmed": false,
    "confirmed_at": null,
    "group_id": null,
    "group_name": "",
    "created_at": "2026-01-03T12:05:00Z"
  }
}
```

**Save the settlement ID:**
```bash
export SETTLEMENT_ID="<settlement-uuid-from-response>"
```

---

### Step 5: View Settlement History

**Bob checks his settlements:**
```bash
curl -X GET http://localhost:8080/api/v1/settlements \
  -H "Authorization: Bearer $BOB_TOKEN"
```

**Alice checks her settlements:**
```bash
curl -X GET http://localhost:8080/api/v1/settlements \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

**With pagination:**
```bash
curl -X GET "http://localhost:8080/api/v1/settlements?limit=10&offset=0" \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

---

### Step 6: View Settlements Between Users

**Alice views settlements with Bob:**
```bash
curl -X GET "http://localhost:8080/api/v1/settlements/between?user_id=$BOB_ID" \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Settlements retrieved successfully",
  "data": [
    {
      "id": "settlement-uuid",
      "payer_id": "bob-uuid",
      "payer_name": "Bob Johnson",
      "payee_id": "alice-uuid",
      "payee_name": "Alice Smith",
      "amount": 250.00,
      "payment_method": "upi",
      "notes": "Payment for lunch",
      "is_confirmed": false,
      "confirmed_at": null,
      "group_id": null,
      "group_name": "",
      "created_at": "2026-01-03T12:05:00Z"
    }
  ]
}
```

---

### Step 7: Get Specific Settlement

**Alice views the settlement details:**
```bash
curl -X GET http://localhost:8080/api/v1/settlements/$SETTLEMENT_ID \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

---

### Step 8: Update Settlement (Before Confirmation)

**Bob updates the settlement (only payer can update):**
```bash
curl -X PATCH http://localhost:8080/api/v1/settlements/$SETTLEMENT_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $BOB_TOKEN" \
  -d '{
    "payment_method": "bank_transfer",
    "notes": "Payment via bank transfer"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Settlement updated successfully",
  "data": {
    "id": "settlement-uuid",
    "payer_id": "bob-uuid",
    "payer_name": "Bob Johnson",
    "payee_id": "alice-uuid",
    "payee_name": "Alice Smith",
    "amount": 250.00,
    "payment_method": "bank_transfer",
    "notes": "Payment via bank transfer",
    "is_confirmed": false,
    "confirmed_at": null,
    "group_id": null,
    "group_name": "",
    "created_at": "2026-01-03T12:05:00Z"
  }
}
```

---

### Step 9: Confirm Settlement

**Alice confirms receipt of payment (only payee can confirm):**
```bash
curl -X PATCH http://localhost:8080/api/v1/settlements/$SETTLEMENT_ID/confirm \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Settlement confirmed successfully",
  "data": {
    "id": "settlement-uuid",
    "payer_id": "bob-uuid",
    "payer_name": "Bob Johnson",
    "payee_id": "alice-uuid",
    "payee_name": "Alice Smith",
    "amount": 250.00,
    "payment_method": "bank_transfer",
    "notes": "Payment via bank transfer",
    "is_confirmed": true,
    "confirmed_at": "2026-01-03T12:10:00Z",
    "group_id": null,
    "group_name": "",
    "created_at": "2026-01-03T12:05:00Z"
  }
}
```

---

### Step 10: Verify Balance Updated

**Alice checks balance after confirmation:**
```bash
curl -X GET http://localhost:8080/api/v1/users/me/balance-summary \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Balance summary retrieved successfully",
  "data": {
    "total_owed": 0,
    "total_owing": 0,
    "net_balance": 0,
    "balances": [],
    "last_updated": "2026-01-03T12:11:00Z"
  }
}
```

**Bob's balance should also be zero now.**

---

## 👥 Part 3: Group Balances & Settlements

### Step 11: Create a Group-

**Alice creates a trip group:**
```bash
curl -X POST http://localhost:8080/api/v1/groups \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -d '{
    "name": "Weekend Trip",
    "description": "Trip to Goa",
    "type": "Trip"
    "member_emails": ["bob@test.com"]
  }'

export GROUP_ID="<group-uuid-from-response>"
```

---

### Step 12: Create Group Expenses

**Alice creates a group expense:**
```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -d '{
    "description": "Hotel booking",
    "amount": 6000.00,
    "category": "accommodation",
    "date": "2026-01-03",
    "split_method": "exact",
    "group_id": "'$GROUP_ID'",
    "participants": [
      {
        "user_id": "'$ALICE_ID'",
        "paid_amount": 6000.00,
        "owed_amount": 3000.00
      },
      {
        "user_id": "'$BOB_ID'",
        "paid_amount": 0,
        "owed_amount": 3000.00
      }
    ]
  }'
```

**Bob creates another group expense:**
```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $BOB_TOKEN" \
  -d '{
    "description": "Car rental",
    "amount": 4000.00,
    "category": "transport",
    "date": "2026-01-03",
    "split_method": "exact",
    "group_id": "'$GROUP_ID'",
    "participants": [
      {
        "user_id": "'$ALICE_ID'",
        "paid_amount": 0,
        "owed_amount": 2000.00
      },
      {
        "user_id": "'$BOB_ID'",
        "paid_amount": 4000.00,
        "owed_amount": 2000.00
      }
    ]
  }'
```

---

### Step 13: View Group Balances

**View group balance breakdown:**
```bash
curl -X GET http://localhost:8080/api/v1/groups/$GROUP_ID/balances \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Group balances retrieved successfully",
  "data": {
    "group_id": "group-uuid",
    "group_name": "Weekend Trip",
    "total_expense": 10000.00,
    "balances": [
      {
        "user_id": "alice-uuid",
        "user_name": "Alice Smith",
        "total_paid": 6000.00,
        "total_owed": 5000.00,
        "net_balance": 1000.00
      },
      {
        "user_id": "bob-uuid",
        "user_name": "Bob Johnson",
        "total_paid": 4000.00,
        "total_owed": 5000.00,
        "net_balance": -1000.00
      }
    ],
    "last_updated": "2026-01-03T12:20:00Z"
  }
}
```

**Note:**
- `net_balance > 0`: User paid more than they owe (should receive money)
- `net_balance < 0`: User owes money to the group

---

### Step 14: Get Group Settlement Suggestions

**Get smart suggestions for the group:**
```bash
curl -X GET http://localhost:8080/api/v1/groups/$GROUP_ID/settlement-suggestions \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Group settlement suggestions generated successfully",
  "data": {
    "suggestions": [
      {
        "from": "bob-uuid",
        "from_name": "Bob Johnson",
        "to": "alice-uuid",
        "to_name": "Alice Smith",
        "amount": 1000.00
      }
    ],
    "total_amount": 1000.00,
    "message": "Here are suggested settlements for the group"
  }
}
```

---

### Step 15: Create Group Settlement

**Bob creates a group settlement:**
```bash
curl -X POST http://localhost:8080/api/v1/settlements \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $BOB_TOKEN" \
  -d '{
    "payee_id": "'$ALICE_ID'",
    "amount": 1000.00,
    "payment_method": "upi",
    "notes": "Trip settlement",
    "group_id": "'$GROUP_ID'"
  }'

export GROUP_SETTLEMENT_ID="<settlement-uuid-from-response>"
```

---

### Step 16: View Group Settlements

**View all settlements for the group:**
```bash
curl -X GET http://localhost:8080/api/v1/groups/$GROUP_ID/settlements \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Settlements retrieved successfully",
  "data": [
    {
      "id": "settlement-uuid",
      "payer_id": "bob-uuid",
      "payer_name": "Bob Johnson",
      "payee_id": "alice-uuid",
      "payee_name": "Alice Smith",
      "amount": 1000.00,
      "payment_method": "upi",
      "notes": "Trip settlement",
      "is_confirmed": false,
      "confirmed_at": null,
      "group_id": "group-uuid",
      "group_name": "Weekend Trip",
      "created_at": "2026-01-03T12:25:00Z"
    }
  ]
}
```

---

### Step 17: Confirm Group Settlement

**Alice confirms the group settlement:**
```bash
curl -X PATCH http://localhost:8080/api/v1/settlements/$GROUP_SETTLEMENT_ID/confirm \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

**Verify group balances are now settled:**
```bash
curl -X GET http://localhost:8080/api/v1/groups/$GROUP_ID/settlement-suggestions \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

**Expected:**
```json
{
  "success": true,
  "message": "Group settlement suggestions generated successfully",
  "data": {
    "suggestions": [],
    "total_amount": 0,
    "message": "Group is all settled up!"
  }
}
```

---

## 🔍 Testing Edge Cases

### 1. Cannot create settlement with yourself
```bash
curl -X POST http://localhost:8080/api/v1/settlements \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -d '{
    "payee_id": "'$ALICE_ID'",
    "amount": 100.00,
    "payment_method": "cash"
  }'
```

**Expected:** Error "cannot create settlement with yourself"

---

### 2. Invalid payment method
```bash
curl -X POST http://localhost:8080/api/v1/settlements \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $BOB_TOKEN" \
  -d '{
    "payee_id": "'$ALICE_ID'",
    "amount": 100.00,
    "payment_method": "invalid_method"
  }'
```

**Expected:** Error "invalid payment method"

**Valid payment methods:**
- `cash`
- `bank_transfer`
- `upi`
- `paypal`
- `venmo`
- `credit_card`
- `debit_card`
- `other`

---

### 3. Only payer can update settlement
```bash
# Alice tries to update Bob's settlement
curl -X PATCH http://localhost:8080/api/v1/settlements/$SETTLEMENT_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -d '{
    "notes": "Updated notes"
  }'
```

**Expected:** Error "only the payer can update this settlement"

---

### 4. Only payee can confirm settlement
```bash
# Bob tries to confirm his own settlement
curl -X PATCH http://localhost:8080/api/v1/settlements/$SETTLEMENT_ID/confirm \
  -H "Authorization: Bearer $BOB_TOKEN"
```

**Expected:** Error "only the payee can confirm this settlement..."

---

### 5. Cannot update confirmed settlement
```bash
# Create and confirm a settlement first, then try to update
curl -X PATCH http://localhost:8080/api/v1/settlements/$SETTLEMENT_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $BOB_TOKEN" \
  -d '{
    "notes": "Try to update"
  }'
```

**Expected:** Error "cannot update a confirmed settlement"

---

### 6. Cannot delete confirmed settlement
```bash
curl -X DELETE http://localhost:8080/api/v1/settlements/$SETTLEMENT_ID \
  -H "Authorization: Bearer $BOB_TOKEN"
```

**Expected:** Error "only the payer can delete this settlement, and it must not be confirmed"

---

### 7. Only payer can delete settlement
```bash
# Alice tries to delete Bob's settlement
curl -X DELETE http://localhost:8080/api/v1/settlements/$SETTLEMENT_ID \
  -H "Authorization: Bearer $ALICE_TOKEN"
```

**Expected:** Error "only the payer can delete this settlement..."

---

### 8. Non-group members cannot create group settlement
```bash
# Create a third user Charlie
# Charlie tries to create settlement for a group he's not in
curl -X POST http://localhost:8080/api/v1/settlements \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CHARLIE_TOKEN" \
  -d '{
    "payee_id": "'$ALICE_ID'",
    "amount": 100.00,
    "payment_method": "cash",
    "group_id": "'$GROUP_ID'"
  }'
```

**Expected:** Error "payer must be a member of the group"

---

### 9. Cannot access another user's settlement
```bash
# Create a third user and try to access Alice-Bob settlement
curl -X GET http://localhost:8080/api/v1/settlements/$SETTLEMENT_ID \
  -H "Authorization: Bearer $CHARLIE_TOKEN"
```

**Expected:** Error "you are not involved in this settlement"

---

## 🗄️ Database Verification

**Check settlements table:**
```sql
psql -h localhost -U postgres -d personal-ess

-- View all settlements
SELECT * FROM settlements WHERE deleted_at IS NULL;

-- View with user details
SELECT 
    s.id,
    s.amount,
    s.payment_method,
    s.is_confirmed,
    s.confirmed_at,
    payer.email as payer_email,
    payee.email as payee_email,
    g.name as group_name
FROM settlements s
JOIN users payer ON s.payer_id = payer.id
JOIN users payee ON s.payee_id = payee.id
LEFT JOIN groups g ON s.group_id = g.id
WHERE s.deleted_at IS NULL;

-- Check balance calculation
SELECT 
    u.email,
    COALESCE(SUM(ep.paid_amount), 0) as total_paid,
    COALESCE(SUM(ep.owed_amount), 0) as total_owed,
    COALESCE(SUM(ep.paid_amount), 0) - COALESCE(SUM(ep.owed_amount), 0) as net_from_expenses
FROM users u
LEFT JOIN expense_participants ep ON u.id = ep.user_id
GROUP BY u.id, u.email;
```

---

## ✅ Validation Checklist

- [x] Migration 000008 applied successfully
- [x] Settlements table created with proper constraints
- [x] All 13 endpoints registered and accessible
- [x] Create settlement works
- [x] Update settlement works (payer only, before confirmation)
- [x] Confirm settlement works (payee only)
- [x] Delete settlement works (payer only, before confirmation)
- [x] Get settlement history works
- [x] Get settlements between users works
- [x] Get group settlements works
- [x] Get balance summary works
- [x] Get detailed balances works
- [x] Get group balances works
- [x] Get settlement suggestions works (greedy algorithm)
- [x] Get group settlement suggestions works
- [x] Balance calculations accurate (expenses - settlements)
- [x] Edge cases handled properly
- [x] User authentication required for all endpoints
- [x] Proper error messages returned
- [x] Database constraints enforced
- [x] Soft delete implemented
- [x] Auto-timestamps working (confirmed_at)

---

## 📊 Phase 4 Status

**Completed:** 13/13 endpoints ✅

**Files Created:**
- ✅ migrations/000008_create_settlements_table.up.sql
- ✅ migrations/000008_create_settlements_table.down.sql
- ✅ internal/models/settlement.go
- ✅ internal/dto/settlement_dto.go
- ✅ internal/repository/settlement_repository.go
- ✅ internal/repository/balance_repository.go
- ✅ internal/services/settlement_service.go
- ✅ internal/services/balance_service.go
- ✅ internal/handler/settlement_handler.go
- ✅ internal/handler/balance_handler.go

**Files Updated:**
- ✅ internal/router/router.go (added settlement and balance routes)
- ✅ cmd/api/main.go (initialized settlement and balance components)

---

## 🎯 Key Features

### Settlement Workflow
1. **Payer creates** settlement (unconfirmed state)
2. **Payer can update/delete** before confirmation
3. **Payee confirms** receipt of payment
4. **Confirmed settlements** are immutable and count toward balance

### Payment Methods Supported
- `cash` - Cash payment
- `bank_transfer` - Direct bank transfer
- `upi` - UPI payment
- `paypal` - PayPal
- `venmo` - Venmo
- `credit_card` - Credit card
- `debit_card` - Debit card
- `other` - Other payment methods

### Balance Calculation
**Formula:** `(Expenses Owed to You) - (Expenses You Owe) - (Settlements You Paid) + (Settlements You Received)`

**Interpretation:**
- Positive balance: Others owe you money
- Negative balance: You owe others money
- Zero balance: All settled up!

### Debt Simplification Algorithm
The settlement suggestions use a **greedy algorithm** to minimize the number of transactions:

1. Calculate all net balances
2. Separate into creditors (positive) and debtors (negative)
3. Sort both by absolute amount (largest first)
4. Match largest creditor with largest debtor
5. Settle the minimum of the two amounts
6. Repeat until all balanced

**Example:**
- **Before:** A owes B $50, B owes C $30, A owes C $20
- **After:** A pays C $50 (1 transaction instead of 3)

---

## 🧪 Advanced Testing Scenarios

### Scenario 1: Complex Multi-User Balance
Create expenses involving 3+ users and verify settlement suggestions minimize transactions.

### Scenario 2: Mixed Personal and Group Settlements
Create both personal and group settlements, verify they're correctly separated.

### Scenario 3: Partial Settlements
Create settlement for less than full amount owed, verify balances update correctly.

### Scenario 4: Settlement Workflow
1. Create settlement
2. Update payment method
3. Add notes
4. Delete and recreate with different amount
5. Confirm
6. Verify immutability

---

## 📈 Performance Tips

### For Large Settlement History
Use pagination:
```bash
curl -X GET "http://localhost:8080/api/v1/settlements?limit=20&offset=0" \
  -H "Authorization: Bearer $TOKEN"
```

### For Specific User Balances
Use the between endpoint instead of fetching all:
```bash
curl -X GET "http://localhost:8080/api/v1/settlements/between?user_id=$OTHER_USER_ID" \
  -H "Authorization: Bearer $TOKEN"
```

### For Group Operations
Always specify group_id to keep settlements organized:
```bash
# Good - linked to group
{
  "payee_id": "...",
  "amount": 100,
  "payment_method": "upi",
  "group_id": "group-uuid"
}
```

---

## 🎉 Next Steps

**Ready for Phase 5: Notifications Module**

When ready, we'll implement:
- Notification system for settlements
- Email notifications
- In-app notifications
- Notification preferences

Just say "Let's start Phase 5" when you're ready! 🚀

---

## 💡 Tips & Tricks

### Quick Balance Check
```bash
# One-liner to see if you're settled up
curl -s http://localhost:8080/api/v1/users/me/balance-summary \
  -H "Authorization: Bearer $TOKEN" | jq '.data.net_balance'
```

### Get Smart Settlement Plan
```bash
# See how to settle all debts with minimum transactions
curl -s http://localhost:8080/api/v1/settlements/suggestions \
  -H "Authorization: Bearer $TOKEN" | jq '.data.suggestions'
```

### Monitor Group Trip Expenses
```bash
# Check who owes what in a group
curl -s http://localhost:8080/api/v1/groups/$GROUP_ID/balances \
  -H "Authorization: Bearer $TOKEN" | jq '.data.balances'
```

---

## 🐛 Troubleshooting

### Balance doesn't match expected value
- Verify all expenses are created correctly
- Check if settlements are confirmed (unconfirmed don't count)
- Ensure no soft-deleted expenses

### Settlement confirmation fails
- Verify you're the payee (only payee can confirm)
- Check if already confirmed

### Cannot update settlement
- Verify you're the payer
- Check if settlement is already confirmed
- Confirmed settlements are immutable

### Group settlement fails
- Verify both payer and payee are group members
- Check group membership with GET /api/v1/groups/:id

---

**Happy Testing! 🎊**

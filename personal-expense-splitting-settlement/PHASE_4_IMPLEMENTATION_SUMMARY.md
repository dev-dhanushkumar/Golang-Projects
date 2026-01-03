# Phase 4: Balance & Settlement Module - Implementation Summary

## Overview
Phase 4 has been successfully implemented, adding balance calculation and settlement tracking functionality to the expense splitting application.

## Files Created

### 1. Database Migration
- **migrations/000008_create_settlements_table.up.sql**
  - Creates `settlements` table for tracking debt settlements between users
  - Fields: payer_id, payee_id, amount, payment_method, is_confirmed, confirmed_at, notes, group_id
  - Constraints: amount > 0, payer != payee, valid payment methods
  - Indexes on payer_id, payee_id, group_id, is_confirmed, created_at
  - Triggers: auto-update updated_at, auto-set confirmed_at on confirmation

- **migrations/000008_create_settlements_table.down.sql**
  - Rollback migration for settlements table

### 2. Models
- **internal/models/settlement.go**
  - Settlement model with PaymentMethod enum (8 types: cash, bank_transfer, upi, paypal, venmo, credit_card, debit_card, other)
  - Business logic methods:
    - `IsValidPaymentMethod()` - validates payment method
    - `CanConfirm(userID)` - checks if user can confirm (must be payee)
    - `CanDelete(userID)` - checks if user can delete (must be payer, before confirmation)
    - `Confirm()` - confirms the settlement
    - `Unconfirm()` - unconfirms the settlement

### 3. DTOs
- **internal/dto/settlement_dto.go**
  - `CreateSettlementRequest` - for creating new settlements
  - `UpdateSettlementRequest` - for partial updates
  - `SettlementResponse` - settlement response with user names
  - `BalanceItem` - individual balance (positive = they owe you, negative = you owe them)
  - `BalanceSummaryResponse` - total owed, owing, and net balance
  - `GroupBalanceItem` - per-member group balance breakdown
  - `GroupBalancesResponse` - group balances response
  - `SettlementSuggestion` - smart settlement suggestion
  - `SettlementSuggestionsResponse` - collection of suggestions

### 4. Repositories
- **internal/repository/settlement_repository.go**
  - CRUD operations for settlements
  - Methods:
    - `Create()` - create new settlement
    - `FindByID()` - find with preloaded relationships
    - `FindByUserID()` - all settlements involving user (as payer or payee)
    - `FindByUsers()` - settlements between two specific users
    - `FindByGroupID()` - group settlements
    - `ConfirmSettlement()` - confirm settlement (only payee)
    - `Update()` - update settlement details
    - `Delete()` - soft delete settlement
    - `CountByUserID()` - count user's settlements

- **internal/repository/balance_repository.go**
  - Complex SQL queries for balance calculations
  - Methods:
    - `CalculateUserBalance(userID, otherUserID)` - calculates balance between two users
      - Formula: (expenses other owes - expenses user owes) - (settlements paid - settlements received)
    - `CalculateAllUserBalances(userID)` - all non-zero balances for a user
      - Uses CTEs for efficient querying
    - `CalculateGroupBalances(groupID)` - per-member group balances
      - Returns total_paid, total_owed, net_balance for each member

### 5. Services
- **internal/services/settlement_service.go**
  - Business logic for settlement operations
  - Validations:
    - Prevents self-settlement
    - Verifies payee exists
    - Verifies group membership for group settlements
    - Only payer can update before confirmation
    - Only payee can confirm
    - Only payer can delete before confirmation
  - Methods:
    - `CreateSettlement()` - create with validations
    - `GetSettlement()` - get with access control
    - `GetUserSettlements()` - user's settlement history
    - `GetSettlementsBetweenUsers()` - settlements between two users
    - `GetGroupSettlements()` - group settlements
    - `UpdateSettlement()` - update (payer only, before confirmation)
    - `ConfirmSettlement()` - confirm (payee only)
    - `DeleteSettlement()` - delete (payer only, before confirmation)

- **internal/services/balance_service.go**
  - Balance calculation and settlement suggestions
  - **Greedy Algorithm for Debt Simplification:**
    1. Separate balances into creditors (positive) and debtors (negative)
    2. Sort both by absolute amount (descending)
    3. Greedily match largest creditor with largest debtor
    4. Settle minimum of the two amounts
    5. Remove settled parties, repeat
    - Result: Minimizes number of transactions to settle all debts
  - Methods:
    - `GetBalanceSummary()` - total owed, owing, net balance
    - `GetUserBalances()` - detailed balances with each person
    - `GetGroupBalances()` - group member balances
    - `GetSettlementSuggestions()` - smart suggestions to minimize transactions
    - `GetGroupSettlementSuggestions()` - group settlement suggestions

### 6. Handlers
- **internal/handler/settlement_handler.go**
  - HTTP handlers for settlement endpoints
  - 8 endpoints total

- **internal/handler/balance_handler.go**
  - HTTP handlers for balance and suggestion endpoints
  - 5 endpoints total

### 7. Router Updates
- **internal/router/router.go**
  - Added SettlementHandler and BalanceHandler to RouterConfig
  - Registered all new routes

### 8. Main Updates
- **cmd/api/main.go**
  - Initialized SettlementRepository and BalanceRepository
  - Initialized SettlementService and BalanceService
  - Initialized SettlementHandler and BalanceHandler
  - Wired everything together

## API Endpoints (13 total)

### Settlement Endpoints (8)
1. **POST /api/v1/settlements**
   - Create a new settlement
   - Body: `{ "payee_id": "uuid", "amount": 100.00, "payment_method": "upi", "notes": "...", "group_id": "uuid" }`

2. **GET /api/v1/settlements**
   - Get user's settlement history
   - Query params: `limit`, `offset`

3. **GET /api/v1/settlements/between?user_id={uuid}**
   - Get settlements between authenticated user and another user

4. **GET /api/v1/settlements/:id**
   - Get specific settlement details

5. **PATCH /api/v1/settlements/:id**
   - Update settlement (only payer, before confirmation)
   - Body: `{ "payment_method": "cash", "notes": "..." }`

6. **PATCH /api/v1/settlements/:id/confirm**
   - Confirm settlement (only payee)

7. **DELETE /api/v1/settlements/:id**
   - Delete settlement (only payer, before confirmation)

8. **GET /api/v1/groups/:id/settlements**
   - Get all settlements for a group

### Balance & Suggestion Endpoints (5)
9. **GET /api/v1/users/me/balance-summary**
   - Get total balance summary (total owed, total owing, net balance)

10. **GET /api/v1/users/me/balances**
    - Get detailed balances with each person

11. **GET /api/v1/groups/:id/balances**
    - Get balance breakdown for group members

12. **GET /api/v1/settlements/suggestions**
    - Get smart settlement suggestions to minimize transactions

13. **GET /api/v1/groups/:id/settlement-suggestions**
    - Get group settlement suggestions

## Key Features

### 1. Settlement Workflow
- **Two-Phase Confirmation:**
  1. Payer creates settlement
  2. Payee confirms receipt
- **State Management:**
  - Before confirmation: Only payer can update/delete
  - After confirmation: Immutable, counts toward balance
- **Automatic Timestamps:**
  - `confirmed_at` auto-set when `is_confirmed` becomes true

### 2. Payment Methods
Supports 8 payment methods:
- cash
- bank_transfer
- upi
- paypal
- venmo
- credit_card
- debit_card
- other

### 3. Balance Calculation
- **Expense-Based:** Calculates what users owe based on expense splits
- **Settlement-Aware:** Subtracts confirmed settlements
- **Direction:**
  - Positive amount: They owe you
  - Negative amount: You owe them
- **Efficient Queries:** Uses CTEs and joins for performance

### 4. Debt Simplification Algorithm
**Greedy Algorithm Benefits:**
- Minimizes number of transactions
- O(n log n) complexity (dominated by sorting)
- Handles complex debt graphs
- **Example:**
  - Before: A owes B $50, B owes C $30, A owes C $20
  - After: A pays C $50 (instead of 3 transactions)

### 5. Group Settlements
- Optional: Can link settlement to a group
- Verification: Both payer and payee must be group members
- Tracking: Separate endpoint for group settlements

## Testing Recommendations

1. **Create Settlement**
   ```bash
   curl -X POST http://localhost:8080/api/v1/settlements \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "payee_id": "uuid-of-friend",
       "amount": 250.50,
       "payment_method": "upi",
       "notes": "Dinner settlement"
     }'
   ```

2. **Get Balance Summary**
   ```bash
   curl http://localhost:8080/api/v1/users/me/balance-summary \
     -H "Authorization: Bearer $TOKEN"
   ```

3. **Get Settlement Suggestions**
   ```bash
   curl http://localhost:8080/api/v1/settlements/suggestions \
     -H "Authorization: Bearer $TOKEN"
   ```

4. **Confirm Settlement** (as payee)
   ```bash
   curl -X PATCH http://localhost:8080/api/v1/settlements/{id}/confirm \
     -H "Authorization: Bearer $TOKEN"
   ```

5. **Get Group Balances**
   ```bash
   curl http://localhost:8080/api/v1/groups/{group_id}/balances \
     -H "Authorization: Bearer $TOKEN"
   ```

## Next Steps

1. Test all 13 endpoints with curl or Postman
2. Verify migration 000008 was applied correctly
3. Test settlement workflow (create → confirm)
4. Test balance calculations with various expense scenarios
5. Test settlement suggestions algorithm
6. Move to Phase 5: Notifications Module

## Notes

- All endpoints require authentication (JWT token)
- Settlements are soft-deleted (deleted_at timestamp)
- Balance calculations exclude soft-deleted expenses and settlements
- Group settlements require group membership verification
- Settlement suggestions are recalculated on each request (no caching)

## Migration Status

```
✅ Migration 000008 applied successfully
✅ Settlements table created
✅ Triggers and constraints active
✅ Server running on port 8080
```

## Phase Completion

**Phase 4 Status: ✅ COMPLETE**

- ✅ Database migrations
- ✅ Models with business logic
- ✅ DTOs for all operations
- ✅ Repositories with complex SQL
- ✅ Services with greedy algorithm
- ✅ Handlers for 13 endpoints
- ✅ Router integration
- ✅ Server running successfully

Ready for testing!

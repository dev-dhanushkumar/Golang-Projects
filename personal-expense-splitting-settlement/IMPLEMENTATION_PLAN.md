# Implementation Plan: Complete API Development

## 📊 Current Status

### ✅ **Phase 0: Foundation (COMPLETED)**
- [x] Authentication Module (Register, Login, Logout, Refresh)
- [x] User Profile (Get, Update)
- [x] Session Management
- [x] File Logging
- [x] SQL Migrations System
- [x] Password Validation
- [x] JWT Authentication
- [x] Phone Number Encryption

**Endpoints Completed**: 8/8 auth endpoints

---

## 🎯 Implementation Roadmap

### **Phase 1: Friendships Module** (Priority: HIGH)
**Duration**: 1-2 days  
**Dependencies**: Users (✅ Done)

#### Database Changes
1. Create migration `000003_create_friendships_table`
   - friendships table (user_id_1, user_id_2, status, requested_by)
   - Indexes for performance

#### Models
- `Friendship` model with status enum (pending, accepted, blocked)

#### DTOs
- `FriendRequestDTO`
- `FriendResponseDTO`
- `FriendListResponseDTO`

#### Repository Layer
- `FriendshipRepository`
  - `SendFriendRequest()`
  - `AcceptFriendRequest()`
  - `RejectFriendRequest()`
  - `BlockFriend()`
  - `GetFriends()`
  - `GetPendingRequests()`
  - `RemoveFriend()`

#### Service Layer
- `FriendshipService`
  - Business logic for friend relationships
  - Validate bidirectional friendship
  - Calculate friend balances (placeholder for now)

#### Handler Layer
- `FriendshipHandler`
  - POST `/api/v1/friends/request`
  - POST `/api/v1/friends/:id/accept`
  - POST `/api/v1/friends/:id/reject`
  - POST `/api/v1/friends/:id/block`
  - DELETE `/api/v1/friends/:id`
  - GET `/api/v1/friends`
  - GET `/api/v1/friends/pending`

**API Endpoints**: 7 new endpoints

---

### **Phase 2: Groups Module** (Priority: HIGH)
**Duration**: 2-3 days  
**Dependencies**: Users (✅), Friendships (Phase 1)

#### Database Changes
1. Migration `000004_create_groups_table`
   - groups table (name, description, created_by, type, image_url)
   
2. Migration `000005_create_group_members_table`
   - group_members table (group_id, user_id, role, joined_at, left_at)
   - UNIQUE constraint on active memberships

#### Models
- `Group` model
- `GroupMember` model with role enum (admin, member)

#### DTOs
- `CreateGroupRequest`
- `UpdateGroupRequest`
- `GroupResponse`
- `GroupDetailResponse`
- `AddMemberRequest`
- `GroupMemberResponse`

#### Repository Layer
- `GroupRepository`
  - `Create()`, `Update()`, `Delete()`
  - `FindByID()`, `FindByUserID()`
  - `GetMembers()`, `AddMember()`, `RemoveMember()`
  - `UpdateMemberRole()`

#### Service Layer
- `GroupService`
  - Create/update/delete groups
  - Manage members (add, remove, update role)
  - Validate permissions (only admins can add/remove)
  - Get group statistics (total expenses, member count)

#### Handler Layer
- `GroupHandler`
  - POST `/api/v1/groups` - Create group
  - GET `/api/v1/groups` - List user's groups
  - GET `/api/v1/groups/:id` - Get group details
  - PATCH `/api/v1/groups/:id` - Update group
  - DELETE `/api/v1/groups/:id` - Delete group
  - POST `/api/v1/groups/:id/members` - Add member
  - DELETE `/api/v1/groups/:id/members/:user_id` - Remove member
  - PATCH `/api/v1/groups/:id/members/:user_id` - Update role

**API Endpoints**: 8 new endpoints

---

### **Phase 3: Expenses Module** (Priority: CRITICAL)
**Duration**: 3-4 days  
**Dependencies**: Users (✅), Groups (Phase 2)

#### Database Changes
1. Migration `000006_create_expenses_table`
   - expenses table (description, amount, category, created_by, group_id, receipt_url)
   
2. Migration `000007_create_expense_participants_table`
   - expense_participants table (expense_id, user_id, paid_amount, owed_amount, is_settled)

#### Models
- `Expense` model with category enum
- `ExpenseParticipant` model

#### DTOs
- `CreateExpenseRequest` with participants array
- `UpdateExpenseRequest`
- `ExpenseResponse`
- `ExpenseDetailResponse`
- `ExpenseParticipantDTO`
- `ExpenseListRequest` (filters: group, category, date range)

#### Repository Layer
- `ExpenseRepository`
  - `Create()`, `Update()`, `Delete()` (soft delete)
  - `FindByID()`, `FindByGroup()`, `FindByUser()`
  - `GetParticipants()`, `UpdateParticipant()`
  
- `ExpenseParticipantRepository`
  - `CreateBulk()` - Add multiple participants
  - `UpdateSettlementStatus()`
  - `GetByExpense()`, `GetByUser()`

#### Service Layer
- `ExpenseService`
  - Create expense with split calculations
  - Validate split amounts (must equal total)
  - Support split methods:
    - Equal split
    - Exact amounts
    - Percentages
    - By shares
  - Update/delete expenses (with validation)
  - Calculate user balances from expenses

#### Handler Layer
- `ExpenseHandler`
  - POST `/api/v1/expenses` - Create expense
  - GET `/api/v1/expenses` - List expenses (with filters)
  - GET `/api/v1/expenses/:id` - Get expense details
  - PATCH `/api/v1/expenses/:id` - Update expense
  - DELETE `/api/v1/expenses/:id` - Delete expense
  - GET `/api/v1/groups/:id/expenses` - Group expenses

**API Endpoints**: 6 new endpoints

---

### **Phase 4: Balance & Settlement Module** (Priority: CRITICAL)
**Duration**: 3-4 days  
**Dependencies**: Expenses (Phase 3)

#### Database Changes
1. Migration `000008_create_settlements_table`
   - settlements table (payer_id, payee_id, amount, payment_method, confirmed)
   
2. Migration `000009_create_account_balances_table`
   - account_balances table (materialized view for performance)
   - Indexes for user pairs

#### Models
- `Settlement` model with payment_method enum
- `AccountBalance` model

#### DTOs
- `CreateSettlementRequest`
- `SettlementResponse`
- `BalanceSummaryResponse`
- `SettlementSuggestionResponse`
- `ConfirmSettlementRequest`

#### Repository Layer
- `SettlementRepository`
  - `Create()`, `FindByID()`, `FindByUsers()`
  - `ConfirmSettlement()`, `GetHistory()`
  
- `BalanceRepository`
  - `CalculateBalance()` - Calculate balance between users
  - `GetUserBalances()` - All balances for a user
  - `UpdateBalance()` - Update materialized view

#### Service Layer
- `BalanceService`
  - Calculate balances from expenses and settlements
  - Generate settlement suggestions (minimize transactions)
  - Algorithm: Simplify debts (greedy algorithm)
  
- `SettlementService`
  - Record settlements
  - Update balances after settlement
  - Confirm settlements (by payee)

#### Handler Layer
- `BalanceHandler`
  - GET `/api/v1/users/me/balance-summary` - User's total balance
  - GET `/api/v1/users/me/balances` - Detailed balances per friend
  - GET `/api/v1/groups/:id/balances` - Group balances
  
- `SettlementHandler`
  - GET `/api/v1/settlements/suggestions` - Smart suggestions
  - POST `/api/v1/settlements` - Record settlement
  - GET `/api/v1/settlements` - Settlement history
  - PATCH `/api/v1/settlements/:id/confirm` - Confirm settlement

**API Endpoints**: 7 new endpoints

---

### **Phase 5: Notifications Module** (Priority: MEDIUM)
**Duration**: 2 days  
**Dependencies**: Expenses (Phase 3), Settlements (Phase 4), Friendships (Phase 1)

#### Database Changes
1. Migration `000010_create_notifications_table`
   - notifications table (user_id, type, title, message, related_entity, is_read)
   - Indexes for user and read status

#### Models
- `Notification` model with type enum

#### DTOs
- `NotificationResponse`
- `NotificationListResponse`

#### Repository Layer
- `NotificationRepository`
  - `Create()`, `CreateBulk()`
  - `MarkAsRead()`, `MarkAllAsRead()`
  - `GetByUser()`, `GetUnreadCount()`
  - `Delete()`

#### Service Layer
- `NotificationService`
  - Create notifications for events:
    - Expense added
    - Settlement received
    - Friend request
    - Group member added
  - Batch create for multiple users
  - Clean up old read notifications

#### Handler Layer
- `NotificationHandler`
  - GET `/api/v1/notifications` - List notifications
  - GET `/api/v1/notifications/unread` - Unread count
  - PATCH `/api/v1/notifications/:id/read` - Mark as read
  - POST `/api/v1/notifications/mark-all-read` - Mark all read
  - DELETE `/api/v1/notifications/:id` - Delete notification

**API Endpoints**: 5 new endpoints

---

### **Phase 6: Payment Methods** (Priority: LOW)
**Duration**: 1-2 days  
**Dependencies**: Users (✅)

#### Database Changes
1. Migration `000011_create_payment_methods_table`
   - payment_methods table (user_id, type, encrypted_details, display_name)

#### Models
- `PaymentMethod` model with type enum

#### DTOs
- `AddPaymentMethodRequest`
- `PaymentMethodResponse`
- `UpdatePaymentMethodRequest`

#### Repository Layer
- `PaymentMethodRepository`
  - `Create()`, `Update()`, `Delete()`
  - `FindByUser()`, `SetPrimary()`

#### Service Layer
- `PaymentMethodService`
  - Encrypt payment details (AES-256)
  - Validate payment method type
  - Manage primary payment method

#### Handler Layer
- `PaymentMethodHandler`
  - POST `/api/v1/payment-methods` - Add payment method
  - GET `/api/v1/payment-methods` - List methods
  - PATCH `/api/v1/payment-methods/:id` - Update method
  - DELETE `/api/v1/payment-methods/:id` - Delete method
  - POST `/api/v1/payment-methods/:id/set-primary` - Set primary

**API Endpoints**: 5 new endpoints

---

### **Phase 7: Audit & Analytics** (Priority: LOW)
**Duration**: 1-2 days  
**Dependencies**: All previous phases

#### Database Changes
1. Migration `000012_create_audit_logs_table`
   - audit_logs table (user_id, action, entity_type, entity_id, old/new values)

#### Models
- `AuditLog` model

#### Repository Layer
- `AuditLogRepository`
  - `Create()`, `GetByEntity()`, `GetByUser()`

#### Service Layer
- `AuditService`
  - Log all financial transactions
  - Track changes to expenses, settlements
  - Generate audit reports

#### Handler Layer
- `AuditHandler`
  - GET `/api/v1/audit/expenses/:id` - Expense history
  - GET `/api/v1/audit/settlements/:id` - Settlement history

**API Endpoints**: 2 new endpoints

---

## 📈 Summary by Phase

| Phase | Module | Endpoints | Duration | Priority |
|-------|--------|-----------|----------|----------|
| 0 | Authentication | 8 | ✅ Done | HIGH |
| 1 | Friendships | 7 | 1-2 days | HIGH |
| 2 | Groups | 8 | 2-3 days | HIGH |
| 3 | Expenses | 6 | 3-4 days | CRITICAL |
| 4 | Balance & Settlement | 7 | 3-4 days | CRITICAL |
| 5 | Notifications | 5 | 2 days | MEDIUM |
| 6 | Payment Methods | 5 | 1-2 days | LOW |
| 7 | Audit & Analytics | 2 | 1-2 days | LOW |
| **TOTAL** | **8 Modules** | **48 Endpoints** | **14-19 days** | - |

---

## 🎯 Recommended Implementation Order

### **Sprint 1: Core Social Features** (3-5 days)
1. Friendships Module (Phase 1)
2. Groups Module (Phase 2)

**Goal**: Users can connect with friends and create groups

### **Sprint 2: Financial Core** (6-8 days)
3. Expenses Module (Phase 3)
4. Balance & Settlement Module (Phase 4)

**Goal**: Users can add expenses, see balances, and settle debts

### **Sprint 3: User Experience** (2 days)
5. Notifications Module (Phase 5)

**Goal**: Users get notified of activities

### **Sprint 4: Optional Features** (2-4 days)
6. Payment Methods Module (Phase 6)
7. Audit & Analytics Module (Phase 7)

**Goal**: Enhanced features and tracking

---

## 🛠️ Technical Approach (Per Phase)

### Standard Implementation Pattern

For each phase, follow this pattern:

1. **Migration First**
   ```bash
   # Create migration files
   touch migrations/00000X_create_table.up.sql
   touch migrations/00000X_create_table.down.sql
   ```

2. **Models**
   ```go
   // internal/models/entity.go
   type Entity struct { ... }
   ```

3. **DTOs**
   ```go
   // internal/dto/entity_dto.go
   type CreateEntityRequest struct { ... }
   type EntityResponse struct { ... }
   ```

4. **Repository**
   ```go
   // internal/repository/entity_repository.go
   type EntityRepository interface { ... }
   type entityRepository struct { ... }
   ```

5. **Service**
   ```go
   // internal/services/entity_service.go
   type EntityService interface { ... }
   type entityService struct { ... }
   ```

6. **Handler**
   ```go
   // internal/handler/entity_handler.go
   type EntityHandler struct { ... }
   ```

7. **Router**
   ```go
   // internal/router/router.go
   // Add routes to SetupRouter
   ```

8. **Testing**
   ```bash
   # Test with curl or Postman
   curl -X POST http://localhost:8080/api/v1/endpoint
   ```

---

## 🔒 Security Considerations (All Phases)

### Authentication
- ✅ All endpoints require authentication (except auth endpoints)
- ✅ Use middleware for JWT validation
- ✅ Check user permissions (e.g., only group admins can remove members)

### Data Validation
- ✅ Validate all inputs with validator package
- ✅ Sanitize user inputs
- ✅ Validate amounts (positive, max 2 decimals)
- ✅ Validate dates (not future for expenses)

### Business Logic Validation
- ✅ Check user is group member before adding expense
- ✅ Check user is friend before seeing balances
- ✅ Validate split amounts equal total expense
- ✅ Prevent duplicate settlements

### Performance
- ✅ Add indexes on foreign keys
- ✅ Use pagination for list endpoints
- ✅ Cache frequently accessed data (consider Redis later)
- ✅ Use materialized views for balances

---

## 📝 Database Migration Checklist

For each new migration:

- [ ] Create `.up.sql` file
- [ ] Create `.down.sql` file
- [ ] Add proper UUID defaults
- [ ] Add foreign key constraints
- [ ] Add indexes for performance
- [ ] Add CHECK constraints where needed
- [ ] Add UNIQUE constraints where needed
- [ ] Add triggers for updated_at
- [ ] Test migration (up and down)
- [ ] Update documentation

---

## 🧪 Testing Strategy

### Per Endpoint
1. **Happy Path**: Test with valid data
2. **Validation**: Test with invalid data
3. **Authentication**: Test without token
4. **Authorization**: Test with wrong user
5. **Edge Cases**: Test boundary conditions

### Integration Testing
- Test complete flows (create group → add expense → settle)
- Test balance calculations
- Test settlement suggestions algorithm

---

## 📊 Progress Tracking

Create a checklist file to track progress:

```markdown
## Phase 1: Friendships
- [ ] Migration created and tested
- [ ] Models implemented
- [ ] DTOs created
- [ ] Repository implemented
- [ ] Service implemented
- [ ] Handler implemented
- [ ] Routes added
- [ ] All endpoints tested
- [ ] Documentation updated

## Phase 2: Groups
...
```

---

## 🚀 Getting Started

### Immediate Next Steps

1. **Review this plan** ✅ (You're here!)

2. **Start Phase 1: Friendships**
   ```bash
   # Create migration files
   touch migrations/000003_create_friendships_table.up.sql
   touch migrations/000003_create_friendships_table.down.sql
   
   # Start coding!
   ```

3. **Test as you go**
   ```bash
   # Keep server running
   go run cmd/api/main.go
   
   # Test each endpoint
   curl -X POST http://localhost:8080/api/v1/friends/request \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"email": "friend@example.com"}'
   ```

4. **Commit frequently**
   ```bash
   git add .
   git commit -m "feat: implement friendships module"
   git push
   ```

---

## 💡 Pro Tips

1. **Start Simple**: Implement basic CRUD first, optimize later
2. **Test Early**: Test each endpoint as you build it
3. **Use Logs**: Check `logs/app.log` for debugging
4. **Follow Patterns**: Maintain consistency with existing code
5. **Document**: Update README and API docs as you go
6. **Commit Often**: Small, focused commits are better
7. **Ask for Help**: Review plan.md when unsure

---

## 📚 Resources

- **Project Plan**: [plan.md](plan.md)
- **Auth Module**: Already implemented, use as reference
- **Migration Guide**: [migrations/README.md](migrations/README.md)
- **API Testing**: [API_TESTING.md](API_TESTING.md)
- **Quick Ref**: [QUICK_REFERENCE.md](QUICK_REFERENCE.md)

---

## 🎯 Success Metrics

By the end of implementation:

- ✅ All 48+ API endpoints working
- ✅ Complete expense splitting functionality
- ✅ Smart settlement suggestions
- ✅ Friend and group management
- ✅ Real-time balance calculations
- ✅ Comprehensive test coverage
- ✅ Production-ready code

---

**Total Estimated Time**: 14-19 days (working solo)  
**Can be accelerated with**: Focused development, parallel work on independent modules

---

## 🚦 Ready to Start?

**Your Next Command**:
```bash
# Start with Phase 1: Friendships
touch migrations/000003_create_friendships_table.up.sql
touch migrations/000003_create_friendships_table.down.sql
```

Then let me know when you're ready, and I'll help you implement the Friendships module! 🚀

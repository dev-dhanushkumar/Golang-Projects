# Personal Expense Splitting & Settlement System
## Project Overview
A group expense management application where users can track shared expenses with friends, split bills, and settle debts. Think "Splitwise" but simpler.

---

## Core Features
### 1. User Management

- User registration and authentication
- Profile management with payment preferences
- Friend/contact management

### 2. Group Management

- Create expense groups (e.g., "Roommates", "Trip to Bali")
- Add/remove members
- View group activity feed

### 3. Expense Tracking

- Add expenses with multiple split methods
- Attach receipts/images
- Categorize expenses
- View expense history

### 4. Balance & Settlement

- Real-time balance calculation per user
- Smart settlement suggestions (minimize transactions)
- Record settlements between users
- Balance history and audit trail

---

## UX Design

Information Architecture

```txt
├── Authentication
│   ├── Login
│   ├── Register
│   └── Forgot Password
│
├── Dashboard
│   ├── Overall Balance Summary
│   ├── Recent Activity Feed
│   ├── Quick Actions (Add Expense, Settle Up)
│   └── Active Groups Overview
│
├── Groups
│   ├── Group List
│   ├── Group Detail
│   │   ├── Expenses Tab
│   │   ├── Balances Tab
│   │   ├── Members Tab
│   │   └── Settings
│   └── Create/Edit Group
│
├── Expenses
│   ├── Add Expense
│   ├── Expense Detail
│   └── Edit/Delete Expense
│
├── Friends
│   ├── Friends List
│   ├── Add Friend
│   └── Friend Balance Detail
│
├── Settlements
│   ├── Settle Up (Smart Suggestions)
│   ├── Record Settlement
│   └── Settlement History
│
└── Profile
    ├── Personal Information
    ├── Payment Methods
    ├── Notification Preferences
    └── Security Settings
```

**Key User Flows**
**Flow 1: Adding an Expense**
```txt
1. User clicks "Add Expense"
2. Enter expense details:
   - Amount
   - Description
   - Category
   - Date
   - Group (optional)
3. Select split method:
   - Equal split
   - Exact amounts
   - Percentages
   - By shares
4. Select participants
5. Upload receipt (optional)
6. Confirm and save
7. System calculates balances
8. Notifications sent to participants
```

**Flow 2: Settling Up**
```txt
1. User views balances dashboard
2. System shows "You owe John $45"
3. User clicks "Settle Up"
4. System suggests optimal settlements
5. User confirms payment method
6. Records settlement
7. Balances updated
8. Both parties notified
```

**UI Screens (Key Elements)**
**Dashboard Screen**
- Balance Card: "You owe $125" / "You are owed $80" with color coding
- Activity Feed: Recent expenses and settlements
- Groups Grid: Cards with group name, member count, your balance
- Floating Action Button: "+ Add Expense"

**Add Expense Screen**
- Amount Input (large, prominent)
- Description Field
- Category Selector (icons: food, transport, entertainment, etc.)
- Date Picker
- Split Method Tabs: Equal | Exact | Percentage | Shares
- Participant Selector (with avatars)
- Receipt Upload Area
- "Save" and "Cancel" buttons

**Group Detail Screen**
- Header: Group name, total expenses, member avatars
- Tabs: Expenses | Balances | Activity
- Balance View: Visual graph showing who owes whom
- "Settle Up" button (prominent)
- Expense List: Scrollable cards with amount, description, date

---

## Database Schema

**Tables**
```sql
-- Users Table
users
├── id (UUID, PK)
├── email (VARCHAR, UNIQUE, INDEXED)
├── password_hash (VARCHAR) -- bcrypt hashed
├── first_name (VARCHAR)
├── last_name (VARCHAR)
├── phone_number (VARCHAR, ENCRYPTED)
├── profile_image_url (VARCHAR)
├── default_currency (VARCHAR, default: 'USD')
├── email_verified (BOOLEAN, default: false)
├── is_active (BOOLEAN, default: true)
├── created_at (TIMESTAMP)
├── updated_at (TIMESTAMP)
└── last_login_at (TIMESTAMP)

-- User Sessions (for auth token management)
user_sessions
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id)
├── token_hash (VARCHAR, INDEXED)
├── refresh_token_hash (VARCHAR)
├── ip_address (INET)
├── user_agent (TEXT)
├── expires_at (TIMESTAMP)
├── created_at (TIMESTAMP)
└── revoked_at (TIMESTAMP, NULLABLE)

-- Groups
groups
├── id (UUID, PK)
├── name (VARCHAR)
├── description (TEXT)
├── created_by (UUID, FK -> users.id)
├── group_type (ENUM: 'trip', 'home', 'couple', 'other')
├── image_url (VARCHAR)
├── is_active (BOOLEAN, default: true)
├── created_at (TIMESTAMP)
└── updated_at (TIMESTAMP)

-- Group Memberships
group_members
├── id (UUID, PK)
├── group_id (UUID, FK -> groups.id)
├── user_id (UUID, FK -> users.id)
├── role (ENUM: 'admin', 'member')
├── joined_at (TIMESTAMP)
└── left_at (TIMESTAMP, NULLABLE)
-- UNIQUE constraint on (group_id, user_id) where left_at IS NULL

-- Friendships
friendships
├── id (UUID, PK)
├── user_id_1 (UUID, FK -> users.id)
├── user_id_2 (UUID, FK -> users.id)
├── status (ENUM: 'pending', 'accepted', 'blocked')
├── requested_by (UUID, FK -> users.id)
├── created_at (TIMESTAMP)
└── updated_at (TIMESTAMP)
-- CHECK constraint: user_id_1 < user_id_2 (to prevent duplicates)

-- Expenses
expenses
├── id (UUID, PK)
├── description (VARCHAR)
├── amount (DECIMAL(12,2))
├── currency (VARCHAR, default: 'USD')
├── category (ENUM: 'food', 'transport', 'entertainment', 'utilities', 'other')
├── expense_date (DATE)
├── created_by (UUID, FK -> users.id)
├── group_id (UUID, FK -> groups.id, NULLABLE)
├── receipt_url (VARCHAR, NULLABLE)
├── notes (TEXT, NULLABLE)
├── is_deleted (BOOLEAN, default: false)
├── created_at (TIMESTAMP)
├── updated_at (TIMESTAMP)
└── deleted_at (TIMESTAMP, NULLABLE)

-- Expense Participants (who paid and who owes)
expense_participants
├── id (UUID, PK)
├── expense_id (UUID, FK -> expenses.id)
├── user_id (UUID, FK -> users.id)
├── paid_amount (DECIMAL(12,2), default: 0) -- what they paid
├── owed_amount (DECIMAL(12,2), default: 0) -- what they owe
├── is_settled (BOOLEAN, default: false)
└── created_at (TIMESTAMP)
-- UNIQUE constraint on (expense_id, user_id)

-- Settlements (recording payments between users)
settlements
├── id (UUID, PK)
├── payer_id (UUID, FK -> users.id) -- who paid
├── payee_id (UUID, FK -> users.id) -- who received
├── amount (DECIMAL(12,2))
├── currency (VARCHAR, default: 'USD')
├── settlement_date (DATE)
├── payment_method (ENUM: 'cash', 'bank_transfer', 'paypal', 'venmo', 'other')
├── notes (TEXT, NULLABLE)
├── group_id (UUID, FK -> groups.id, NULLABLE)
├── confirmed_by_payee (BOOLEAN, default: false)
├── created_at (TIMESTAMP)
└── updated_at (TIMESTAMP)

-- Account Balances (materialized view for performance)
account_balances
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id)
├── other_user_id (UUID, FK -> users.id)
├── balance (DECIMAL(12,2)) -- positive: they owe you, negative: you owe them
├── currency (VARCHAR, default: 'USD')
├── group_id (UUID, FK -> groups.id, NULLABLE)
├── last_updated_at (TIMESTAMP)
└── UNIQUE constraint on (user_id, other_user_id, group_id)

-- Notifications
notifications
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id)
├── type (ENUM: 'expense_added', 'expense_updated', 'settlement_received', 'friend_request')
├── title (VARCHAR)
├── message (TEXT)
├── related_entity_type (VARCHAR) -- 'expense', 'settlement', 'friendship'
├── related_entity_id (UUID)
├── is_read (BOOLEAN, default: false)
├── created_at (TIMESTAMP)
└── read_at (TIMESTAMP, NULLABLE)

-- Audit Log (for financial transactions)
audit_logs
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id, NULLABLE)
├── action (VARCHAR) -- 'create', 'update', 'delete'
├── entity_type (VARCHAR) -- 'expense', 'settlement', 'user'
├── entity_id (UUID)
├── old_values (JSONB, NULLABLE)
├── new_values (JSONB)
├── ip_address (INET)
├── user_agent (TEXT)
└── created_at (TIMESTAMP)

-- Payment Methods (encrypted PII)
payment_methods
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id)
├── type (ENUM: 'bank_account', 'paypal', 'venmo', 'upi')
├── encrypted_details (BYTEA) -- encrypted account numbers, etc.
├── display_name (VARCHAR) -- e.g., "Chase ****1234"
├── is_primary (BOOLEAN, default: false)
├── is_active (BOOLEAN, default: true)
├── created_at (TIMESTAMP)
└── updated_at (TIMESTAMP)
```

**Indexes**

```sql
-- Performance indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_user_sessions_token ON user_sessions(token_hash);
CREATE INDEX idx_expenses_created_by ON expenses(created_by);
CREATE INDEX idx_expenses_group_id ON expenses(group_id);
CREATE INDEX idx_expenses_date ON expenses(expense_date DESC);
CREATE INDEX idx_expense_participants_user ON expense_participants(user_id);
CREATE INDEX idx_settlements_payer ON settlements(payer_id);
CREATE INDEX idx_settlements_payee ON settlements(payee_id);
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read, created_at DESC);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_account_balances_user ON account_balances(user_id, other_user_id);
```

---

## API Schema

### Authentication Endpoints

```bash
POST /api/v1/auth/register
Request:
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "+1234567890"
}
Response: 201 Created
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2025-01-15T10:00:00Z"
  },
  "tokens": {
    "access_token": "jwt_token",
    "refresh_token": "refresh_jwt",
    "expires_in": 3600
  }
}

POST /api/v1/auth/login
Request:
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
Response: 200 OK
{
  "user": {...},
  "tokens": {...}
}

POST /api/v1/auth/logout
Headers: Authorization: Bearer {token}
Response: 204 No Content

POST /api/v1/auth/refresh
Request:
{
  "refresh_token": "refresh_jwt"
}
Response: 200 OK
{
  "access_token": "new_jwt_token",
  "expires_in": 3600
}
```

### User Endpoints
```bash
GET /api/v1/users/me
Headers: Authorization: Bearer {token}
Response: 200 OK
{
  "id": "uuid",
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "profile_image_url": "url",
  "default_currency": "USD",
  "created_at": "2025-01-15T10:00:00Z"
}

PATCH /api/v1/users/me
Headers: Authorization: Bearer {token}
Request:
{
  "first_name": "Jonathan",
  "default_currency": "EUR"
}
Response: 200 OK

GET /api/v1/users/me/balance-summary
Response: 200 OK
{
  "total_owed_to_you": 125.50,
  "total_you_owe": 80.00,
  "net_balance": 45.50,
  "currency": "USD",
  "breakdown": [
    {
      "user": {
        "id": "uuid",
        "name": "Jane Smith"
      },
      "balance": 45.50,
      "status": "owed_to_you"
    }
  ]
}
```

### Group Endpoints
```bash
GET /api/v1/groups
Query params: ?page=1&limit=20&status=active
Response: 200 OK
{
  "groups": [
    {
      "id": "uuid",
      "name": "Roommates",
      "description": "Apartment expenses",
      "member_count": 4,
      "your_balance": -25.00,
      "total_expenses": 1250.00,
      "created_at": "2025-01-10T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 5
  }
}

POST /api/v1/groups
Request:
{
  "name": "Bali Trip 2025",
  "description": "Group trip expenses",
  "group_type": "trip",
  "member_ids": ["uuid1", "uuid2"]
}
Response: 201 Created

GET /api/v1/groups/{group_id}
Response: 200 OK
{
  "id": "uuid",
  "name": "Roommates",
  "description": "Apartment expenses",
  "members": [
    {
      "id": "uuid",
      "name": "John Doe",
      "role": "admin",
      "balance": -25.00
    }
  ],
  "total_expenses": 1250.00,
  "expense_count": 45,
  "created_at": "2025-01-10T10:00:00Z"
}

POST /api/v1/groups/{group_id}/members
Request:
{
  "user_id": "uuid",
  "role": "member"
}
Response: 201 Created
```

### Expense Endpoints
```bash
POST /api/v1/expenses
Request:
{
  "description": "Dinner at Italian Restaurant",
  "amount": 120.00,
  "currency": "USD",
  "category": "food",
  "expense_date": "2025-01-18",
  "group_id": "uuid",
  "split_method": "equal",
  "participants": [
    {
      "user_id": "uuid1",
      "paid_amount": 120.00,
      "owed_amount": 40.00
    },
    {
      "user_id": "uuid2",
      "paid_amount": 0,
      "owed_amount": 40.00
    },
    {
      "user_id": "uuid3",
      "paid_amount": 0,
      "owed_amount": 40.00
    }
  ],
  "notes": "Great meal!"
}
Response: 201 Created
{
  "id": "uuid",
  "description": "Dinner at Italian Restaurant",
  "amount": 120.00,
  "created_by": {...},
  "participants": [...],
  "created_at": "2025-01-18T19:30:00Z"
}

GET /api/v1/expenses
Query params: ?group_id=uuid&page=1&limit=20&category=food&from_date=2025-01-01
Response: 200 OK
{
  "expenses": [
    {
      "id": "uuid",
      "description": "Dinner",
      "amount": 120.00,
      "currency": "USD",
      "category": "food",
      "expense_date": "2025-01-18",
      "created_by": {
        "id": "uuid",
        "name": "John Doe"
      },
      "your_share": 40.00,
      "participants_count": 3,
      "is_settled": false
    }
  ],
  "pagination": {...}
}

GET /api/v1/expenses/{expense_id}
Response: 200 OK
{
  "id": "uuid",
  "description": "Dinner at Italian Restaurant",
  "amount": 120.00,
  "participants": [
    {
      "user": {"id": "uuid", "name": "John Doe"},
      "paid_amount": 120.00,
      "owed_amount": 40.00,
      "is_settled": false
    }
  ],
  "receipt_url": "url",
  "notes": "Great meal!",
  "created_at": "2025-01-18T19:30:00Z"
}

PUT /api/v1/expenses/{expense_id}
DELETE /api/v1/expenses/{expense_id}
```

### Settlement Endpoints
```bash
GET /api/v1/settlements/suggestions
Query params: ?group_id=uuid (optional)
Response: 200 OK
{
  "suggestions": [
    {
      "from_user": {"id": "uuid", "name": "You"},
      "to_user": {"id": "uuid", "name": "Jane Smith"},
      "amount": 45.50,
      "currency": "USD"
    }
  ],
  "total_transactions": 1,
  "simplified_from": 5
}

POST /api/v1/settlements
Request:
{
  "payee_id": "uuid",
  "amount": 45.50,
  "currency": "USD",
  "settlement_date": "2025-01-19",
  "payment_method": "bank_transfer",
  "group_id": "uuid",
  "notes": "Paid via Venmo"
}
Response: 201 Created

GET /api/v1/settlements
Query params: ?page=1&limit=20&from_date=2025-01-01
Response: 200 OK
{
  "settlements": [
    {
      "id": "uuid",
      "payer": {"id": "uuid", "name": "John Doe"},
      "payee": {"id": "uuid", "name": "Jane Smith"},
      "amount": 45.50,
      "payment_method": "bank_transfer",
      "settlement_date": "2025-01-19",
      "confirmed": true,
      "created_at": "2025-01-19T14:30:00Z"
    }
  ],
  "pagination": {...}
}

PATCH /api/v1/settlements/{settlement_id}/confirm
Response: 200 OK
```

### Friendship Endpoints
```bash
GET /api/v1/friends
Response: 200 OK
{
  "friends": [
    {
      "user": {
        "id": "uuid",
        "name": "Jane Smith",
        "email": "jane@example.com"
      },
      "balance": 45.50,
      "status": "owed_to_you",
      "friendship_since": "2025-01-01T10:00:00Z"
    }
  ]
}

POST /api/v1/friends/request
Request:
{
  "email": "friend@example.com"
}
Response: 201 Created

POST /api/v1/friends/{friendship_id}/accept
PATCH /api/v1/friends/{friendship_id}/reject
DELETE /api/v1/friends/{friendship_id}
```

### Notification Endpoints
```bash
GET /api/v1/notifications
Query params: ?unread=true&page=1&limit=20
Response: 200 OK
{
  "notifications": [
    {
      "id": "uuid",
      "type": "expense_added",
      "title": "New expense added",
      "message": "John added 'Dinner' for $120.00",
      "is_read": false,
      "related_entity": {
        "type": "expense",
        "id": "uuid"
      },
      "created_at": "2025-01-18T19:30:00Z"
    }
  ],
  "unread_count": 5
}

PATCH /api/v1/notifications/{notification_id}/read
POST /api/v1/notifications/mark-all-read
```
---
## Security Considerations
### Authentication & Authorization
- JWT tokens with access (1hr) and refresh tokens (30 days)
- Password requirements: Min 8 chars, uppercase, lowercase, number, special char
- Bcrypt hashing for passwords (cost factor: 12)
- Rate limiting: 5 failed login attempts = 15 min lockout
- Session management: Track active sessions, allow revocation
- Role-based access: Group admins vs members

### Data Protection (PII)
- Encryption at rest: AES-256 for sensitive fields
    - Phone numbers
    - Payment method details
    - Personal notes
- Encryption in transit: TLS 1.3
- Database encryption: Column-level encryption for PII
- Secure key management: Use KMS (Key Management Service)

### API Security
**HTTPS only** - no HTTP allowed
**CORS configuration:** Whitelist frontend domains
**Rate limiting:**
    - Authentication: 5 req/min per IP
    - API calls: 100 req/min per user
    - Expense creation: 10 req/min per user
**Input validation:**
    - Sanitize all inputs
    - Validate amounts (positive, max 2 decimals)
    - Validate dates (not future for expenses)
- **SQL injection prevention:** Use parameterized queries
- **CSRF protection:** CSRF tokens for state-changing operations

### Financial Transaction Security
- **Idempotency keys:** Prevent duplicate expense/settlement creation
- **Audit trail:** Log all financial transactions
- **Balance validation:** Server-side recalculation, never trust client
- **Transaction atomicity:** Use database transactions
- **Soft deletes:** Never hard delete financial records
- **Immutability:**Expenses can't be edited after 24hrs, only deleted and recreated

### Privacy & Compliance
- **Data minimization**: Only collect necessary data
- **User consent**: Clear terms for data usage
- **Data retention**: Auto-delete inactive accounts after 2 years (with notice)
- **Right to deletion:** API endpoint to delete user data
- **Export capability**: Users can download their data (GDPR)
- **Anonymization:** Remove PII from analytics

### Infrastructure Security
- **Environment variables:** Never hardcode secrets
- **Secrets management:** Use vault for API keys, DB credentials
- **Database security:**
    - Separate read/write users
    - Principle of least privilege
    - Regular backups (encrypted)
- **Logging:** Log security events (failed logins, unusual activity)
- **Monitoring:** Alert on suspicious patterns

---

## Technical Stack Recommendations
### Backend (Golang)
- **Framework:** Gin or Echo
- **Database:** PostgreSQL 15+
- **ORM:** sqlx or GORM
- **Authentication:** golang-jwt/jwt
- **Encryption:** crypto/aes, crypto/bcrypt
- **Image storage:** AWS S3 or similar
- **Caching:** Redis (session, balance cache)
- **Queue:** For async notifications (RabbitMQ/Redis)

### Additional Tools

**API Documentation**: OpenAPI/Swagger
**Testing:** testify, gomock
**Migration:** golang-migrate
**Logging:** zerolog or zap
**Monitoring:** Prometheus + Grafana


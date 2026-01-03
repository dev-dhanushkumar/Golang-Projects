# Implementation Guide: Logging & Migrations

## 🎯 What's New

Two major improvements have been implemented:

1. **File-based Logging** - Logs are now saved to files AND displayed in console
2. **Manual SQL Migrations** - Database schema managed via SQL files instead of AutoMigrate

---

## 📝 1. File-based Logging

### How It Works

The logger now writes to **both console and file**:
- **Console**: Human-readable format for development
- **File**: JSON format in `logs/app.log` for production/analysis

### Configuration

**Location**: [pkg/logger/logger.go](pkg/logger/logger.go)

```go
// Logs are written to:
// 1. Console (human-readable, colored)
// 2. File: logs/app.log (JSON format)

logger, err := logger.InitLogger()
```

### Log Levels

Based on `ENVIRONMENT` variable:
- **development**: DEBUG level and above
- **production**: INFO level and above

### Log Files

```
logs/
└── app.log          # Main application log (JSON format)
```

**Features**:
- Auto-creates `logs/` directory
- Appends to existing log file
- JSON format for easy parsing
- Includes caller information (file:line)
- Stack traces on ERROR level

### Example Log Output

**Console** (development):
```
2026-01-03T10:15:23.456Z  INFO    cmd/api/main.go:25  Application starting up...
2026-01-03T10:15:23.789Z  DEBUG   internal/config/config.go:45  Configuration loaded successfully  {"db_host": "localhost"}
```

**File** (logs/app.log - JSON):
```json
{"level":"info","ts":"2026-01-03T10:15:23.456Z","caller":"cmd/api/main.go:25","msg":"Application starting up..."}
{"level":"debug","ts":"2026-01-03T10:15:23.789Z","caller":"internal/config/config.go:45","msg":"Configuration loaded successfully","db_host":"localhost"}
```

### Usage in Code

```go
// Already initialized in main.go
sugar.Info("Simple message")
sugar.Infow("Message with fields", "key", "value", "count", 42)
sugar.Debugw("Debug info", "user_id", userId)
sugar.Errorw("Error occurred", "error", err)
sugar.Fatalw("Critical error", "error", err) // Exits app
```

### Log Rotation (Future Enhancement)

For production, consider adding log rotation:

```go
// Using lumberjack
import "gopkg.in/natefinch/lumberjack.v2"

logFile := &lumberjack.Logger{
    Filename:   "logs/app.log",
    MaxSize:    100, // megabytes
    MaxBackups: 3,
    MaxAge:     28, // days
    Compress:   true,
}
```

### Viewing Logs

```bash
# View latest logs
tail -f logs/app.log

# View with formatting (requires jq)
tail -f logs/app.log | jq .

# Search logs
grep "error" logs/app.log | jq .

# Filter by level
jq 'select(.level=="error")' logs/app.log
```

---

## 🗄️ 2. Manual SQL Migrations

### Why Manual Migrations?

**Problems with AutoMigrate**:
- ❌ Limited control over schema changes
- ❌ UUID extension timing issues
- ❌ Can't handle complex many-to-many relationships well
- ❌ No rollback capability
- ❌ Difficult to version control

**Benefits of SQL Migrations**:
- ✅ Full control over schema
- ✅ Proper UUID extension setup
- ✅ Easy to handle complex relationships
- ✅ Rollback support
- ✅ Version controlled
- ✅ Production-safe

### Migration Structure

```
migrations/
├── README.md
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_user_sessions_table.up.sql
└── 000002_create_user_sessions_table.down.sql
```

### How It Works

1. **On App Start**: 
   - Checks `schema_migrations` table
   - Identifies pending migrations
   - Runs them in order
   - Records successful migrations

2. **Migration Tracking**:
   ```sql
   CREATE TABLE schema_migrations (
       version VARCHAR(255) PRIMARY KEY,
       name VARCHAR(255) NOT NULL,
       applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```

3. **Each Migration**:
   - Runs in a transaction
   - Either fully succeeds or fully rolls back
   - Recorded only if successful

### Using Migrations

#### Auto-run on App Start

Migrations run automatically when you start the app:

```bash
go run cmd/api/main.go
```

Output:
```
2026-01-03T10:15:24.123Z  INFO  Starting database migrations...
2026-01-03T10:15:24.234Z  INFO  Applying migration  {"version": "000001", "name": "create_users_table"}
2026-01-03T10:15:24.345Z  INFO  Migration applied successfully  {"version": "000001"}
2026-01-03T10:15:24.456Z  INFO  Applying migration  {"version": "000002", "name": "create_user_sessions_table"}
2026-01-03T10:15:24.567Z  INFO  Migration applied successfully  {"version": "000002"}
2026-01-03T10:15:24.678Z  INFO  All migrations completed successfully
```

#### Manual Migration Tool

Use the CLI tool for manual control:

```bash
# Check migration status
go run cmd/migrate/main.go status

# Output:
# Migration Status:
# ================
# 000001 - create_users_table [✅ Applied]
# 000002 - create_user_sessions_table [✅ Applied]

# Run migrations manually
go run cmd/migrate/main.go migrate

# Rollback last migration
go run cmd/migrate/main.go rollback
```

### Creating New Migrations

#### Step 1: Create Migration Files

Next version would be `000003`:

```bash
# Create the files
touch migrations/000003_create_groups_table.up.sql
touch migrations/000003_create_groups_table.down.sql
```

#### Step 2: Write UP Migration

**migrations/000003_create_groups_table.up.sql**:
```sql
-- Create groups table
CREATE TABLE IF NOT EXISTS groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID NOT NULL,
    group_type VARCHAR(50) DEFAULT 'other',
    image_url VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key to users
    CONSTRAINT fk_groups_created_by 
        FOREIGN KEY (created_by) 
        REFERENCES users(id) 
        ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_groups_created_by ON groups(created_by);
CREATE INDEX IF NOT EXISTS idx_groups_is_active ON groups(is_active);

-- Add auto-update trigger
CREATE TRIGGER update_groups_updated_at
    BEFORE UPDATE ON groups
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

#### Step 3: Write DOWN Migration

**migrations/000003_create_groups_table.down.sql**:
```sql
-- Drop trigger
DROP TRIGGER IF EXISTS update_groups_updated_at ON groups;

-- Drop indexes
DROP INDEX IF EXISTS idx_groups_is_active;
DROP INDEX IF EXISTS idx_groups_created_by;

-- Drop table
DROP TABLE IF EXISTS groups;
```

#### Step 4: Test

```bash
# Start app to run migration
go run cmd/api/main.go

# Or use migration tool
go run cmd/migrate/main.go migrate

# Verify in database
psql -h localhost -U postgres -d personal-ess
\dt groups
\d groups
```

### Existing Migrations

#### Migration 000001 - Users Table

Creates:
- `users` table with UUID primary key
- Email index
- `updated_at` auto-update trigger function
- PostgreSQL extensions (uuid-ossp, pgcrypto)

Features:
- UUID auto-generation
- Encrypted phone number field (500 chars for encryption)
- Automatic timestamp management

#### Migration 000002 - User Sessions Table

Creates:
- `user_sessions` table
- Foreign key to users (CASCADE delete)
- Multiple indexes for performance
- Composite index for active sessions

Features:
- Proper UUID relationships
- Optimized for token lookups
- Supports session revocation

### Database Setup

#### First Time Setup

1. **Create Database**:
   ```bash
   # Connect to PostgreSQL
   docker exec -it postgres_arch psql -U postgres
   
   # Create database
   CREATE DATABASE "personal-ess";
   
   # Exit
   \q
   ```

2. **Configure .env**:
   ```bash
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=personal-ess
   DB_SSL_MODE=disable
   ```

3. **Run Migrations**:
   ```bash
   # Start app (migrations run automatically)
   go run cmd/api/main.go
   
   # Or use migration tool
   go run cmd/migrate/main.go migrate
   ```

#### Verify Setup

```sql
-- Connect to database
psql -h localhost -U postgres -d personal-ess

-- Check tables
\dt

-- Should show:
-- public | schema_migrations | table | postgres
-- public | users            | table | postgres
-- public | user_sessions    | table | postgres

-- Check migrations
SELECT * FROM schema_migrations ORDER BY version;

-- Should show:
-- version | name                      | applied_at
-----------+---------------------------+------------
-- 000001  | create_users_table        | 2026-01-03...
-- 000002  | create_user_sessions_table| 2026-01-03...
```

### Common Patterns

#### Many-to-Many Relationship

Example: Groups and Members

```sql
-- migrations/000004_create_group_members.up.sql

CREATE TABLE IF NOT EXISTS group_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(50) DEFAULT 'member',
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP,
    
    -- Foreign keys
    CONSTRAINT fk_group_members_group 
        FOREIGN KEY (group_id) 
        REFERENCES groups(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_group_members_user 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE,
    
    -- Prevent duplicate active memberships
    CONSTRAINT unique_active_membership 
        UNIQUE (group_id, user_id, left_at)
);

-- Indexes
CREATE INDEX idx_group_members_group ON group_members(group_id);
CREATE INDEX idx_group_members_user ON group_members(user_id);
CREATE INDEX idx_group_members_active ON group_members(group_id, user_id) 
    WHERE left_at IS NULL;
```

#### Adding Columns

```sql
-- migrations/000005_add_bio_to_users.up.sql

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS bio TEXT;

ALTER TABLE users
ADD CONSTRAINT check_bio_length 
    CHECK (length(bio) <= 500);
```

#### Modifying Columns

```sql
-- migrations/000006_extend_description.up.sql

-- Make description longer
ALTER TABLE groups 
ALTER COLUMN description TYPE TEXT;

-- Add NOT NULL constraint (with default for existing rows)
UPDATE groups SET description = '' WHERE description IS NULL;
ALTER TABLE groups 
ALTER COLUMN description SET NOT NULL;
```

### Troubleshooting

#### Migration Fails

**Check logs**:
```bash
tail -n 50 logs/app.log | jq 'select(.level=="error")'
```

**Common issues**:
1. **Table already exists**: Check if migration was partially applied
2. **Foreign key violation**: Ensure parent tables exist
3. **Syntax error**: Validate SQL syntax

**Solution**:
```sql
-- Check what's in migrations table
SELECT * FROM schema_migrations;

-- If migration is recorded but failed, remove it
DELETE FROM schema_migrations WHERE version = 'XXXXX';

-- Fix the migration file
-- Run again
```

#### UUID Extension Missing

**Error**: `function uuid_generate_v4() does not exist`

**Solution**: Run migration 000001 which creates the extension:
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

#### Can't Rollback

**Error**: No down migration found

**Solution**: Create the `.down.sql` file that reverses the `.up.sql`:
- Drop tables in reverse order
- Remove added columns
- Remove constraints

---

## 📦 File Structure Changes

```
personal-expense-splitting-settlement/
├── cmd/
│   ├── api/
│   │   └── main.go              # Updated: Uses RunMigrationsFromFiles
│   └── migrate/
│       └── main.go              # NEW: Migration CLI tool
├── internal/
│   └── database/
│       ├── database.go          # Updated: Removed AutoMigrate
│       └── migrations.go        # NEW: Migration runner
├── logs/                        # NEW: Auto-created
│   └── app.log                  # NEW: Application logs (git-ignored)
├── migrations/                  # NEW: Migration files
│   ├── README.md               # NEW: Migration guide
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_user_sessions_table.up.sql
│   └── 000002_create_user_sessions_table.down.sql
├── pkg/
│   └── logger/
│       └── logger.go            # Updated: File + console logging
└── .gitignore                   # NEW: Ignores logs/
```

---

## ✅ Testing Checklist

### Logging
- [ ] Start app and check console output
- [ ] Verify `logs/app.log` is created
- [ ] Check log file contains JSON entries
- [ ] Test different log levels (Info, Debug, Error)
- [ ] Verify logs persist across app restarts

### Migrations
- [ ] Drop all tables if they exist
- [ ] Start app - migrations should run automatically
- [ ] Check `schema_migrations` table has 2 entries
- [ ] Verify `users` table exists with correct schema
- [ ] Verify `user_sessions` table exists
- [ ] Check foreign key constraint works
- [ ] Test migration status: `go run cmd/migrate/main.go status`
- [ ] Test rollback: `go run cmd/migrate/main.go rollback`
- [ ] Re-run migrations after rollback

---

## 🚀 Next Steps

### Recommended Actions

1. **Test Both Features**:
   ```bash
   # Clean start
   psql -h localhost -U postgres -d personal-ess -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
   
   # Run app
   go run cmd/api/main.go
   
   # Check logs
   tail -f logs/app.log | jq .
   ```

2. **Create Next Migration** (when ready for groups):
   ```bash
   # Create files
   touch migrations/000003_create_groups_table.up.sql
   touch migrations/000003_create_groups_table.down.sql
   
   # Edit files with your SQL
   # Run app or migration tool
   ```

3. **Monitor Logs**:
   ```bash
   # Development
   tail -f logs/app.log | jq .
   
   # Production (with filtering)
   jq 'select(.level | IN("error", "warn"))' logs/app.log
   ```

### Production Considerations

1. **Log Rotation**: Add lumberjack for automatic log rotation
2. **Log Aggregation**: Consider ELK stack or similar
3. **Migration Monitoring**: Set up alerts for failed migrations
4. **Backup Strategy**: Always backup before running migrations in production
5. **Rollback Plan**: Test rollback procedures before production deployment

---

## 📚 Resources

- **Migrations Guide**: See [migrations/README.md](migrations/README.md)
- **Zap Logger Docs**: https://pkg.go.dev/go.uber.org/zap
- **PostgreSQL Docs**: https://www.postgresql.org/docs/

---

**Summary**: Your application now has production-ready logging and database migration management! 🎉

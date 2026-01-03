# Database Migrations Guide

## Overview

This project uses **manual SQL migrations** instead of GORM's AutoMigrate. This approach provides:
- ✅ Better control over database schema changes
- ✅ Proper handling of UUID extensions
- ✅ Support for complex relationships (many-to-many, etc.)
- ✅ Version control for database changes
- ✅ Easy rollback capability

## Migration Files Structure

```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_user_sessions_table.up.sql
├── 000002_create_user_sessions_table.down.sql
└── ... (future migrations)
```

### Naming Convention

Format: `{version}_{description}.{direction}.sql`

- **version**: 6-digit number (e.g., `000001`, `000002`)
- **description**: Snake_case description (e.g., `create_users_table`)
- **direction**: Either `up` (apply) or `down` (rollback)

Examples:
- `000001_create_users_table.up.sql` - Creates users table
- `000001_create_users_table.down.sql` - Drops users table
- `000003_add_groups_table.up.sql` - Creates groups table

## How Migrations Work

### 1. Tracking Table

A `schema_migrations` table tracks applied migrations:

```sql
CREATE TABLE schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### 2. Running Migrations

Migrations run automatically when the application starts:

```go
// In cmd/api/main.go
if err := database.RunMigrationsFromFiles(sugar); err != nil {
    sugar.Fatalw("Failed to run migrations", "error", err)
}
```

The system:
1. Checks which migrations are already applied
2. Runs only pending migrations in order
3. Records each successful migration
4. Stops on first error

## Creating New Migrations

### Step 1: Determine Next Version Number

Check existing migrations and increment:

```bash
ls migrations/ | grep up.sql | tail -1
# If last is 000002, your new version is 000003
```

### Step 2: Create UP Migration

Create file: `migrations/000003_your_description.up.sql`

Example - Adding groups table:

```sql
-- migrations/000003_create_groups_table.up.sql

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
    
    CONSTRAINT fk_groups_created_by 
        FOREIGN KEY (created_by) 
        REFERENCES users(id) 
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_groups_created_by ON groups(created_by);
CREATE INDEX IF NOT EXISTS idx_groups_is_active ON groups(is_active);

-- Add trigger for updated_at
CREATE TRIGGER update_groups_updated_at
    BEFORE UPDATE ON groups
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

### Step 3: Create DOWN Migration

Create file: `migrations/000003_your_description.down.sql`

```sql
-- migrations/000003_create_groups_table.down.sql

DROP TRIGGER IF EXISTS update_groups_updated_at ON groups;
DROP INDEX IF EXISTS idx_groups_is_active;
DROP INDEX IF EXISTS idx_groups_created_by;
DROP TABLE IF EXISTS groups;
```

### Step 4: Test Migration

```bash
# Start the app - migrations run automatically
go run cmd/api/main.go

# Check logs for:
# "Applying migration" - Shows migration being applied
# "Migration applied successfully" - Confirms success
```

### Step 5: Verify in Database

```sql
-- Check if table was created
\dt groups

-- Check migration was recorded
SELECT * FROM schema_migrations ORDER BY version;

-- Check table structure
\d groups
```

## Managing Migrations

### Check Migration Status

```sql
SELECT version, name, applied_at 
FROM schema_migrations 
ORDER BY version;
```

### Manual Rollback (If Needed)

You can create a CLI tool or run SQL directly:

```sql
-- Get last migration
SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;

-- Execute the down migration manually
-- Then remove from tracking
DELETE FROM schema_migrations WHERE version = '000003';
```

### Reset All Migrations (Development Only)

**⚠️ WARNING: This deletes all data!**

```sql
-- Drop all tables
DROP TABLE IF EXISTS user_sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS schema_migrations CASCADE;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;

-- Drop extensions (optional)
DROP EXTENSION IF EXISTS "uuid-ossp" CASCADE;
DROP EXTENSION IF EXISTS "pgcrypto" CASCADE;
```

Then restart the app to run all migrations from scratch.

## Best Practices

### 1. Never Edit Applied Migrations

Once a migration is applied (especially in production), **never edit it**. Instead:
- Create a new migration to make changes
- Example: If you forgot a column, create `000004_add_column_to_users.up.sql`

### 2. Test Rollbacks

Always test that your `.down.sql` file correctly reverses the `.up.sql`:

```bash
# 1. Apply migration (start app)
# 2. Verify changes in DB
# 3. Manually run down migration
# 4. Verify everything is reverted
```

### 3. Use Transactions Implicitly

Each migration runs in a transaction. If any statement fails, the entire migration rolls back.

### 4. Handle UUID Properly

For UUID primary keys:

```sql
-- Good: Let PostgreSQL generate UUIDs
CREATE TABLE example (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ...
);

-- Bad: Relying on application to generate IDs
CREATE TABLE example (
    id UUID PRIMARY KEY,  -- No default, app must provide
    ...
);
```

### 5. Foreign Keys

Always include proper foreign key constraints:

```sql
CREATE TABLE child_table (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id UUID NOT NULL,
    
    CONSTRAINT fk_child_parent
        FOREIGN KEY (parent_id)
        REFERENCES parent_table(id)
        ON DELETE CASCADE  -- or RESTRICT, SET NULL based on needs
);
```

### 6. Indexes

Add indexes for:
- Foreign keys
- Frequently queried columns
- Columns used in WHERE clauses
- Columns used in ORDER BY

```sql
CREATE INDEX IF NOT EXISTS idx_table_column ON table_name(column_name);

-- Composite index for multiple columns
CREATE INDEX IF NOT EXISTS idx_table_multi 
    ON table_name(column1, column2);

-- Partial index for common filtered queries
CREATE INDEX IF NOT EXISTS idx_active_users 
    ON users(email) 
    WHERE is_active = TRUE;
```

## Existing Migrations

### 000001 - Users Table

Creates the `users` table with:
- UUID primary key
- Email (unique, indexed)
- Password hash
- Encrypted phone number
- Profile information
- Timestamps with auto-update trigger

### 000002 - User Sessions Table

Creates the `user_sessions` table with:
- UUID primary key
- Foreign key to users
- Token hashes for authentication
- IP and user agent tracking
- Expiration and revocation support
- Indexes for performance

## Troubleshooting

### Migration Fails to Apply

**Error**: "migration X failed to apply"

**Solution**:
1. Check logs for specific SQL error
2. Verify SQL syntax in migration file
3. Check if table/column already exists
4. Ensure foreign key references exist

### UUID Extension Not Found

**Error**: "function uuid_generate_v4() does not exist"

**Solution**:
Migration `000001` should create this extension. If not:

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

### Migration Already Applied

**Error**: Migration shows as applied but table doesn't exist

**Solution**:
```sql
-- Check what's actually in schema_migrations
SELECT * FROM schema_migrations;

-- If migration record exists but table doesn't, manually remove record
DELETE FROM schema_migrations WHERE version = 'XXXXX';

-- Then restart app to reapply
```

### Can't Drop Table (Dependencies)

**Error**: "cannot drop table X because other objects depend on it"

**Solution**:
```sql
-- Use CASCADE to drop dependent objects
DROP TABLE table_name CASCADE;

-- Or in migration, drop in correct order (children before parents)
```

## Future Enhancements

When you need to add:

### Many-to-Many Relationships

Example: Users and Groups

```sql
-- migrations/000005_create_group_members.up.sql

CREATE TABLE IF NOT EXISTS group_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(50) DEFAULT 'member',
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP,
    
    CONSTRAINT fk_group_members_group_id
        FOREIGN KEY (group_id)
        REFERENCES groups(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_group_members_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    
    -- Ensure user can't be added to same group twice (while active)
    CONSTRAINT unique_active_membership 
        UNIQUE (group_id, user_id, left_at)
);

CREATE INDEX IF NOT EXISTS idx_group_members_group ON group_members(group_id);
CREATE INDEX IF NOT EXISTS idx_group_members_user ON group_members(user_id);
```

### Adding Columns

```sql
-- migrations/000006_add_avatar_to_users.up.sql

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(500);

-- Create index if needed
CREATE INDEX IF NOT EXISTS idx_users_avatar ON users(avatar_url) 
    WHERE avatar_url IS NOT NULL;
```

### Modifying Columns

```sql
-- migrations/000007_extend_user_bio.up.sql

-- Increase bio length
ALTER TABLE users 
ALTER COLUMN bio TYPE TEXT;

-- Add constraint
ALTER TABLE users
ADD CONSTRAINT check_bio_length 
    CHECK (length(bio) <= 1000);
```

## CLI Commands Reference

### Check Migration Status

```bash
# Connect to database
psql -h localhost -U postgres -d personal-ess

# Check applied migrations
SELECT version, name, applied_at FROM schema_migrations ORDER BY version;
```

### Manual Migration Operations

```bash
# Apply specific migration manually
psql -h localhost -U postgres -d personal-ess < migrations/000003_create_groups.up.sql

# Rollback specific migration manually
psql -h localhost -U postgres -d personal-ess < migrations/000003_create_groups.down.sql
```

---

## Summary

✅ **Migrations run automatically** on app start
✅ **Never edit applied migrations** - create new ones
✅ **Always create both up and down** migrations
✅ **Test migrations** before committing
✅ **Use proper constraints** and indexes
✅ **Follow naming conventions** for consistency

For questions or issues, check the troubleshooting section or review existing migration files as examples.

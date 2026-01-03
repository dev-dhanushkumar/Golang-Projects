# Quick Start: Logging & Migrations

## 🚀 Get Started in 3 Steps

### Step 1: Clean Database (First Time Only)

```bash
# Connect to PostgreSQL
docker exec -it postgres_arch psql -U postgres

# Drop and recreate database
DROP DATABASE IF EXISTS "personal-ess";
CREATE DATABASE "personal-ess";

# Exit
\q
```

### Step 2: Run the Application

```bash
# Start the app (migrations run automatically)
go run cmd/api/main.go
```

**Expected Output**:
```
2026-01-03T10:15:24.123Z  INFO  Application starting up...
2026-01-03T10:15:24.234Z  INFO  Database connection established successfully
2026-01-03T10:15:24.345Z  INFO  Starting database migrations...
2026-01-03T10:15:24.456Z  INFO  Applying migration  {"version": "000001", "name": "create_users_table"}
2026-01-03T10:15:24.567Z  INFO  Migration applied successfully  {"version": "000001"}
2026-01-03T10:15:24.678Z  INFO  Applying migration  {"version": "000002", "name": "create_user_sessions_table"}
2026-01-03T10:15:24.789Z  INFO  Migration applied successfully  {"version": "000002"}
2026-01-03T10:15:24.890Z  INFO  Database migrations completed successfully
2026-01-03T10:15:25.000Z  INFO  Server Starting  {"address": ":8080"}
```

### Step 3: Verify Everything Works

```bash
# In another terminal, check the logs file
tail -f logs/app.log | jq .

# Test the API
curl http://localhost:8080/health

# Check database
psql -h localhost -U postgres -d personal-ess -c "\dt"
```

**Expected Tables**:
```
               List of relations
 Schema |       Name        | Type  |  Owner   
--------+-------------------+-------+----------
 public | schema_migrations | table | postgres
 public | user_sessions     | table | postgres
 public | users             | table | postgres
```

---

## ✅ What Changed?

### 1. **Logging** 📝

**Before**: Console only
```
INFO  Application starting...
```

**After**: Console + File
```
Console:  INFO  Application starting...
File:     {"level":"info","ts":"2026-01-03T10:15:24.123Z","msg":"Application starting..."}
```

**Log File Location**: `logs/app.log`

### 2. **Migrations** 🗄️

**Before**: GORM AutoMigrate
```go
DB.AutoMigrate(&models.User{}, &models.UserSession{})
```

**After**: SQL Migration Files
```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_user_sessions_table.up.sql
└── 000002_create_user_sessions_table.down.sql
```

---

## 🔧 Using the Migration Tool

### Check Status

```bash
go run cmd/migrate/main.go status
```

**Output**:
```
Migration Status:
================
000001 - create_users_table [✅ Applied]
000002 - create_user_sessions_table [✅ Applied]
```

### Run Migrations Manually

```bash
go run cmd/migrate/main.go migrate
```

### Rollback Last Migration

```bash
go run cmd/migrate/main.go rollback
```

**⚠️ Warning**: This will drop the last table (user_sessions)!

---

## 📋 Common Tasks

### View Logs

```bash
# Live tail (formatted JSON)
tail -f logs/app.log | jq .

# View all logs
cat logs/app.log | jq .

# Filter by level
jq 'select(.level=="error")' logs/app.log

# Search for specific text
grep "migration" logs/app.log | jq .
```

### Check Database

```bash
# Connect
psql -h localhost -U postgres -d personal-ess

# List tables
\dt

# Check users table
\d users

# Check migrations
SELECT * FROM schema_migrations ORDER BY version;

# Check if UUID extension is enabled
\dx
```

### Reset Everything

```bash
# Drop all tables
psql -h localhost -U postgres -d personal-ess << EOF
DROP TABLE IF EXISTS user_sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS schema_migrations CASCADE;
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;
EOF

# Restart app (migrations will run again)
go run cmd/api/main.go
```

---

## 🐛 Troubleshooting

### Issue: "Failed to connect to database"

**Solution**:
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Start it if not running
docker container start postgres_arch

# Verify connection
psql -h localhost -U postgres -c "SELECT version();"
```

### Issue: "Table already exists"

**Solution**:
```sql
-- Check what migrations were applied
SELECT * FROM schema_migrations;

-- If migration is recorded but table doesn't exist, remove record
DELETE FROM schema_migrations WHERE version = '000001';

-- Restart app
```

### Issue: "logs/app.log not created"

**Solution**:
```bash
# Check permissions
ls -la logs/

# Create manually if needed
mkdir -p logs
touch logs/app.log
chmod 644 logs/app.log

# Restart app
```

### Issue: "Migration fails with SQL error"

**Check logs for details**:
```bash
tail -n 100 logs/app.log | jq 'select(.level=="error")'
```

**Common causes**:
- Missing PostgreSQL extensions
- Foreign key references to non-existent tables
- SQL syntax errors

---

## 📊 What to Expect

### On First Run

1. ✅ Creates `logs/` directory
2. ✅ Creates `logs/app.log` file
3. ✅ Connects to database
4. ✅ Creates `schema_migrations` table
5. ✅ Runs migration 000001 (users table)
6. ✅ Runs migration 000002 (user_sessions table)
7. ✅ Starts server on port 8080

### On Subsequent Runs

1. ✅ Appends to existing `logs/app.log`
2. ✅ Connects to database
3. ✅ Checks for pending migrations
4. ✅ Skips already-applied migrations
5. ✅ Runs only new migrations (if any)
6. ✅ Starts server

---

## 🎯 Next Steps

### 1. Test the Authentication API

```bash
# Register a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### 2. Check Logs

```bash
# You should see the request logged
tail -f logs/app.log | jq .
```

### 3. Verify Database

```sql
-- Check if user was created
SELECT id, email, first_name, last_name, created_at FROM users;

-- Check if session was created
SELECT id, user_id, ip_address, created_at FROM user_sessions;
```

### 4. Create Your First Migration

When you're ready to add groups:

```bash
# Create files
touch migrations/000003_create_groups_table.up.sql
touch migrations/000003_create_groups_table.down.sql

# Edit the files with your SQL
# See migrations/README.md for examples

# Run app to apply
go run cmd/api/main.go
```

---

## 📚 Full Documentation

- **Detailed Guide**: [LOGGING_MIGRATIONS_GUIDE.md](LOGGING_MIGRATIONS_GUIDE.md)
- **Migrations Guide**: [migrations/README.md](migrations/README.md)
- **API Testing**: [API_TESTING.md](API_TESTING.md)

---

## ✨ Summary

✅ **Logging**: Both console and file (`logs/app.log`)
✅ **Migrations**: SQL files in `migrations/` folder
✅ **Auto-run**: Migrations run on app start
✅ **CLI Tool**: Manual control via `cmd/migrate/main.go`
✅ **Production-ready**: Proper schema version control

**You're all set!** 🎉

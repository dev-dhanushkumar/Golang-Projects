-- Drop indexes
DROP INDEX IF EXISTS idx_user_sessions_active;
DROP INDEX IF EXISTS idx_user_sessions_expire_at;
DROP INDEX IF EXISTS idx_user_sessions_refresh_token_hash;
DROP INDEX IF EXISTS idx_user_sessions_token_hash;
DROP INDEX IF EXISTS idx_user_sessions_user_id;

-- Drop user_sessions table
DROP TABLE IF EXISTS user_sessions;

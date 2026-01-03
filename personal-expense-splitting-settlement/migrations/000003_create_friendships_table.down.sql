-- Drop trigger
DROP TRIGGER IF EXISTS trigger_update_friendships_updated_at ON friendships;

-- Drop function
DROP FUNCTION IF EXISTS update_friendships_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_friendships_users_status;
DROP INDEX IF EXISTS idx_friendships_requested_by;
DROP INDEX IF EXISTS idx_friendships_status;
DROP INDEX IF EXISTS idx_friendships_user_id_2;
DROP INDEX IF EXISTS idx_friendships_user_id_1;

-- Drop table
DROP TABLE IF EXISTS friendships;

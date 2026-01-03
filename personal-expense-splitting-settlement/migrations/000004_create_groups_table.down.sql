DROP TRIGGER IF EXISTS trigger_update_groups_updated_at ON groups;
DROP FUNCTION IF EXISTS update_groups_updated_at();
DROP TABLE IF EXISTS groups CASCADE;

DROP TRIGGER IF EXISTS trigger_check_group_has_admin ON group_members;
DROP TRIGGER IF EXISTS trigger_update_group_members_updated_at ON group_members;
DROP FUNCTION IF EXISTS check_group_has_admin();
DROP FUNCTION IF EXISTS update_group_members_updated_at();
DROP TABLE IF EXISTS group_members CASCADE;

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_set_settlement_confirmed_at ON settlements;
DROP TRIGGER IF EXISTS trigger_update_settlements_updated_at ON settlements;

-- Drop functions
DROP FUNCTION IF EXISTS set_settlement_confirmed_at();
DROP FUNCTION IF EXISTS update_settlements_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_settlements_created_at;
DROP INDEX IF EXISTS idx_settlements_is_confirmed;
DROP INDEX IF EXISTS idx_settlements_group_id;
DROP INDEX IF EXISTS idx_settlements_payee_id;
DROP INDEX IF EXISTS idx_settlements_payer_id;

-- Drop table
DROP TABLE IF EXISTS settlements;

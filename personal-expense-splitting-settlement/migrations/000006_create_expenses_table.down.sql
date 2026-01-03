DROP TRIGGER IF EXISTS trigger_update_expenses_updated_at ON expenses;
DROP FUNCTION IF EXISTS update_expenses_updated_at();
DROP TABLE IF EXISTS expenses CASCADE;

DROP TRIGGER IF EXISTS trigger_set_expense_participant_settled_at ON expense_participants;
DROP TRIGGER IF EXISTS trigger_update_expense_participants_updated_at ON expense_participants;
DROP FUNCTION IF EXISTS set_expense_participant_settled_at();
DROP FUNCTION IF EXISTS update_expense_participants_updated_at();
DROP TABLE IF EXISTS expense_participants CASCADE;

CREATE TABLE IF NOT EXISTS expense_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expense_id UUID NOT NULL REFERENCES expenses(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    paid_amount DECIMAL(12, 2) NOT NULL DEFAULT 0,
    owed_amount DECIMAL(12, 2) NOT NULL DEFAULT 0,
    is_settled BOOLEAN NOT NULL DEFAULT FALSE,
    settled_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT expense_participants_paid_amount_non_negative CHECK (paid_amount >= 0),
    CONSTRAINT expense_participants_owed_amount_non_negative CHECK (owed_amount >= 0),
    CONSTRAINT expense_participants_unique_user_expense UNIQUE (expense_id, user_id)
);

CREATE INDEX idx_expense_participants_expense_id ON expense_participants(expense_id);
CREATE INDEX idx_expense_participants_user_id ON expense_participants(user_id);
CREATE INDEX idx_expense_participants_is_settled ON expense_participants(is_settled);
CREATE INDEX idx_expense_participants_unsettled ON expense_participants(user_id, is_settled) WHERE is_settled = FALSE;

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_expense_participants_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_expense_participants_updated_at
    BEFORE UPDATE ON expense_participants
    FOR EACH ROW
    EXECUTE FUNCTION update_expense_participants_updated_at();

-- Trigger to set settled_at when is_settled becomes true
CREATE OR REPLACE FUNCTION set_expense_participant_settled_at()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_settled = TRUE AND OLD.is_settled = FALSE THEN
        NEW.settled_at = CURRENT_TIMESTAMP;
    ELSIF NEW.is_settled = FALSE THEN
        NEW.settled_at = NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_expense_participant_settled_at
    BEFORE UPDATE ON expense_participants
    FOR EACH ROW
    EXECUTE FUNCTION set_expense_participant_settled_at();

-- Create settlements table for tracking debt settlements between users
CREATE TABLE IF NOT EXISTS settlements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payer_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    payee_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    amount DECIMAL(12, 2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL DEFAULT 'cash',
    notes TEXT,
    is_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    confirmed_at TIMESTAMP WITH TIME ZONE,
    group_id UUID REFERENCES groups(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT settlements_amount_positive CHECK (amount > 0),
    CONSTRAINT settlements_different_users CHECK (payer_id != payee_id),
    CONSTRAINT settlements_payment_method_valid CHECK (payment_method IN (
        'cash', 'bank_transfer', 'upi', 'paypal', 'venmo', 
        'credit_card', 'debit_card', 'other'
    ))
);

-- Indexes for performance
CREATE INDEX idx_settlements_payer_id ON settlements(payer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_settlements_payee_id ON settlements(payee_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_settlements_group_id ON settlements(group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_settlements_is_confirmed ON settlements(is_confirmed) WHERE deleted_at IS NULL;
CREATE INDEX idx_settlements_created_at ON settlements(created_at DESC) WHERE deleted_at IS NULL;

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_settlements_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_settlements_updated_at
    BEFORE UPDATE ON settlements
    FOR EACH ROW
    EXECUTE FUNCTION update_settlements_updated_at();

-- Trigger to set confirmed_at when is_confirmed becomes true
CREATE OR REPLACE FUNCTION set_settlement_confirmed_at()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_confirmed = TRUE AND OLD.is_confirmed = FALSE THEN
        NEW.confirmed_at = CURRENT_TIMESTAMP;
    ELSIF NEW.is_confirmed = FALSE THEN
        NEW.confirmed_at = NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_settlement_confirmed_at
    BEFORE UPDATE ON settlements
    FOR EACH ROW
    EXECUTE FUNCTION set_settlement_confirmed_at();

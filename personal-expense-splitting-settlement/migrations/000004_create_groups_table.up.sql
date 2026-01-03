CREATE TABLE IF NOT EXISTS groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL DEFAULT 'general',
    image_url TEXT,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT groups_name_not_empty CHECK (TRIM(name) <> ''),
    CONSTRAINT groups_type_valid CHECK (type IN ('general', 'trip', 'home', 'couple', 'event', 'project', 'other'))
);

CREATE INDEX idx_groups_created_by ON groups(created_by) WHERE deleted_at IS NULL;
CREATE INDEX idx_groups_type ON groups(type) WHERE deleted_at IS NULL;
CREATE INDEX idx_groups_created_at ON groups(created_at DESC) WHERE deleted_at IS NULL;

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_groups_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_groups_updated_at
    BEFORE UPDATE ON groups
    FOR EACH ROW
    EXECUTE FUNCTION update_groups_updated_at();

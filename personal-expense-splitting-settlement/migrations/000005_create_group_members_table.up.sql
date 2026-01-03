CREATE TABLE IF NOT EXISTS group_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT group_members_role_valid CHECK (role IN ('admin', 'member')),
    CONSTRAINT group_members_unique_active_membership UNIQUE (group_id, user_id) 
        DEFERRABLE INITIALLY DEFERRED
);

-- Note: The unique constraint allows only one active membership per user per group
-- When left_at IS NOT NULL, the membership is considered inactive

CREATE INDEX idx_group_members_group_id ON group_members(group_id) WHERE left_at IS NULL;
CREATE INDEX idx_group_members_user_id ON group_members(user_id) WHERE left_at IS NULL;
CREATE INDEX idx_group_members_role ON group_members(role) WHERE left_at IS NULL;
CREATE INDEX idx_group_members_joined_at ON group_members(joined_at DESC);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_group_members_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_group_members_updated_at
    BEFORE UPDATE ON group_members
    FOR EACH ROW
    EXECUTE FUNCTION update_group_members_updated_at();

-- Ensure at least one admin exists in each group
CREATE OR REPLACE FUNCTION check_group_has_admin()
RETURNS TRIGGER AS $$
DECLARE
    admin_count INTEGER;
BEGIN
    -- If we're removing the last admin, prevent it
    IF (TG_OP = 'UPDATE' AND OLD.role = 'admin' AND NEW.role = 'member') OR
       (TG_OP = 'UPDATE' AND OLD.left_at IS NULL AND NEW.left_at IS NOT NULL AND OLD.role = 'admin') OR
       (TG_OP = 'DELETE' AND OLD.role = 'admin' AND OLD.left_at IS NULL) THEN
        
        SELECT COUNT(*) INTO admin_count
        FROM group_members
        WHERE group_id = OLD.group_id
          AND role = 'admin'
          AND left_at IS NULL
          AND id != OLD.id;
        
        IF admin_count = 0 THEN
            RAISE EXCEPTION 'Cannot remove the last admin from the group';
        END IF;
    END IF;
    
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_group_has_admin
    BEFORE UPDATE OR DELETE ON group_members
    FOR EACH ROW
    EXECUTE FUNCTION check_group_has_admin();

ALTER TABLE refresh_tokens
    ADD COLUMN family_id TEXT,
    ADD COLUMN rotated_from_token_id INTEGER REFERENCES refresh_tokens(id) ON DELETE SET NULL,
    ADD COLUMN replaced_by_token_id INTEGER REFERENCES refresh_tokens(id) ON DELETE SET NULL,
    ADD COLUMN revoked_at TIMESTAMPTZ,
    ADD COLUMN revoke_reason TEXT;

UPDATE refresh_tokens
SET family_id = token_hash
WHERE family_id IS NULL OR family_id = '';

ALTER TABLE refresh_tokens
    ALTER COLUMN family_id SET NOT NULL;

CREATE INDEX idx_refresh_tokens_family_id ON refresh_tokens(family_id);
CREATE INDEX idx_refresh_tokens_revoked_at ON refresh_tokens(revoked_at);
CREATE INDEX idx_refresh_tokens_active_user_id ON refresh_tokens(user_id) WHERE revoked_at IS NULL AND deleted_at IS NULL;

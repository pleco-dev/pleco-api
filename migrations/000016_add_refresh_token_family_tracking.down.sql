DROP INDEX IF EXISTS idx_refresh_tokens_active_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_revoked_at;
DROP INDEX IF EXISTS idx_refresh_tokens_family_id;

ALTER TABLE refresh_tokens
    DROP COLUMN IF EXISTS revoke_reason,
    DROP COLUMN IF EXISTS revoked_at,
    DROP COLUMN IF EXISTS replaced_by_token_id,
    DROP COLUMN IF EXISTS rotated_from_token_id,
    DROP COLUMN IF EXISTS family_id;

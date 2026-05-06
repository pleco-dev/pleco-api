DROP INDEX IF EXISTS idx_social_accounts_deleted_at;
DROP INDEX IF EXISTS idx_role_permissions_deleted_at;
DROP INDEX IF EXISTS idx_roles_deleted_at;
DROP INDEX IF EXISTS idx_permissions_deleted_at;
DROP INDEX IF EXISTS idx_email_verification_tokens_deleted_at;
DROP INDEX IF EXISTS idx_users_last_password_change;
DROP INDEX IF EXISTS idx_users_last_login_at;

ALTER TABLE social_accounts DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE role_permissions
DROP COLUMN IF EXISTS deleted_at,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE roles
DROP COLUMN IF EXISTS deleted_at,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE permissions
DROP COLUMN IF EXISTS deleted_at,
DROP COLUMN IF EXISTS updated_at,
DROP COLUMN IF EXISTS created_at;

ALTER TABLE email_verification_tokens
DROP COLUMN IF EXISTS deleted_at,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE users
DROP COLUMN IF EXISTS last_password_change,
DROP COLUMN IF EXISTS last_login_at;

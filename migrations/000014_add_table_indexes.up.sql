CREATE INDEX IF NOT EXISTS idx_users_role ON users (role);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users (role_id);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_user_id ON email_verification_tokens (user_id);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_token ON email_verification_tokens (token);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_expires_at ON email_verification_tokens (expires_at);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_created_at ON refresh_tokens (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_device ON refresh_tokens (user_id, device_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_deleted_at ON refresh_tokens (deleted_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expired_at ON refresh_tokens (expired_at);

CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions (role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission ON role_permissions (permission);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_permission ON role_permissions (role_id, permission);

CREATE INDEX IF NOT EXISTS idx_social_accounts_user_id ON social_accounts (user_id);
CREATE INDEX IF NOT EXISTS idx_social_accounts_user_provider ON social_accounts (user_id, provider);

CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_status ON audit_logs (status);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource_id ON audit_logs (resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_actor_created_at ON audit_logs (actor_user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource_status_created_at ON audit_logs (resource, status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_investigations_status ON audit_investigations (status);
CREATE INDEX IF NOT EXISTS idx_audit_investigations_ai_provider ON audit_investigations (ai_provider);
CREATE INDEX IF NOT EXISTS idx_audit_investigations_ai_model ON audit_investigations (ai_model);
CREATE INDEX IF NOT EXISTS idx_audit_investigations_created_by_created_at ON audit_investigations (created_by_user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_error_analyses_deleted_at ON error_analyses (deleted_at);
CREATE INDEX IF NOT EXISTS idx_error_analyses_created_at ON error_analyses (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_error_analyses_error_type ON error_analyses (error_type);
CREATE INDEX IF NOT EXISTS idx_error_analyses_severity ON error_analyses (severity);

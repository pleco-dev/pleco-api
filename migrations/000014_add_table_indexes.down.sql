DROP INDEX IF EXISTS idx_error_analyses_severity;
DROP INDEX IF EXISTS idx_error_analyses_error_type;
DROP INDEX IF EXISTS idx_error_analyses_created_at;
DROP INDEX IF EXISTS idx_error_analyses_deleted_at;

DROP INDEX IF EXISTS idx_audit_investigations_created_by_created_at;
DROP INDEX IF EXISTS idx_audit_investigations_ai_model;
DROP INDEX IF EXISTS idx_audit_investigations_ai_provider;
DROP INDEX IF EXISTS idx_audit_investigations_status;

DROP INDEX IF EXISTS idx_audit_logs_resource_status_created_at;
DROP INDEX IF EXISTS idx_audit_logs_actor_created_at;
DROP INDEX IF EXISTS idx_audit_logs_resource_id;
DROP INDEX IF EXISTS idx_audit_logs_status;
DROP INDEX IF EXISTS idx_audit_logs_created_at;

DROP INDEX IF EXISTS idx_social_accounts_user_provider;
DROP INDEX IF EXISTS idx_social_accounts_user_id;

DROP INDEX IF EXISTS idx_role_permissions_role_permission;
DROP INDEX IF EXISTS idx_role_permissions_permission;
DROP INDEX IF EXISTS idx_role_permissions_role_id;

DROP INDEX IF EXISTS idx_refresh_tokens_expired_at;
DROP INDEX IF EXISTS idx_refresh_tokens_deleted_at;
DROP INDEX IF EXISTS idx_refresh_tokens_user_device;
DROP INDEX IF EXISTS idx_refresh_tokens_user_created_at;

DROP INDEX IF EXISTS idx_email_verification_tokens_expires_at;
DROP INDEX IF EXISTS idx_email_verification_tokens_token;
DROP INDEX IF EXISTS idx_email_verification_tokens_user_id;

DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_role_id;
DROP INDEX IF EXISTS idx_users_role;

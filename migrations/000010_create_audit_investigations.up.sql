CREATE TABLE IF NOT EXISTS audit_investigations (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    created_by_user_id BIGINT NULL,
    action TEXT NOT NULL DEFAULT '',
    resource TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT '',
    actor_user_id BIGINT NULL,
    search TEXT NOT NULL DEFAULT '',
    date_from TIMESTAMPTZ NULL,
    date_to TIMESTAMPTZ NULL,
    limit_value INTEGER NOT NULL DEFAULT 50,
    log_count INTEGER NOT NULL DEFAULT 0,
    ai_provider TEXT NOT NULL DEFAULT '',
    ai_model TEXT NOT NULL DEFAULT '',
    summary TEXT NOT NULL DEFAULT '',
    timeline_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    suspicious_signals_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    recommendations_json JSONB NOT NULL DEFAULT '[]'::jsonb
);

CREATE INDEX IF NOT EXISTS idx_audit_investigations_deleted_at ON audit_investigations (deleted_at);
CREATE INDEX IF NOT EXISTS idx_audit_investigations_created_at ON audit_investigations (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_investigations_created_by_user_id ON audit_investigations (created_by_user_id);
CREATE INDEX IF NOT EXISTS idx_audit_investigations_resource_status ON audit_investigations (resource, status);

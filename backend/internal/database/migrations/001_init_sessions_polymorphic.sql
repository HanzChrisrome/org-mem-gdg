-- Migration: 001_init_sessions_polymorphic.sql
-- Description: Standardize sessions table to support both members and executives

CREATE TABLE IF NOT EXISTS sessions (
    session_id VARCHAR(64) PRIMARY KEY,
    owner_id VARCHAR(255) NOT NULL, -- UUID or string depending on source
    owner_type VARCHAR(20) NOT NULL CHECK (owner_type IN ('member', 'executive')),
    refresh_token_hash VARCHAR(64) NOT NULL,
    user_agent TEXT,
    ip_address VARCHAR(45),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ,

    -- Ensure indexed lookups for user/executive sessions
    CONSTRAINT idx_sessions_owner UNIQUE(owner_id, owner_type, created_at)
);

CREATE INDEX IF NOT EXISTS idx_sessions_owner_lookup ON sessions (owner_id, owner_type);
CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions (expires_at);

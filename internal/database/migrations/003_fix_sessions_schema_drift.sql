-- Migration: 003_fix_sessions_schema_drift.sql
-- Description: Standardize sessions table to resolve 'integer = uuid' and 'string = integer' mismatches

-- Drop conflicting session definitions
DROP TABLE IF EXISTS "public"."sessions";

-- Recreate sessions table with consistent owner_id type (VARCHAR) to match Go config.Session struct
CREATE TABLE "public"."sessions" (
    "session_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "refresh_token_id" VARCHAR(64) NOT NULL UNIQUE,
    "owner_id" VARCHAR(255) NOT NULL, -- Keep as VARCHAR to match Go string ID and handle polymorphism
    "owner_type" VARCHAR(20) NOT NULL CHECK (owner_type IN ('member', 'executive')),
    "refresh_token_hash" VARCHAR(64) NOT NULL,
    "user_agent" TEXT,
    "ip_address" VARCHAR(45),
    "expires_at" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "revoked_at" TIMESTAMPTZ
);

-- Indices for performance
CREATE INDEX IF NOT EXISTS idx_sessions_refresh_token_lookup ON public.sessions (refresh_token_id);
CREATE INDEX IF NOT EXISTS idx_sessions_owner_lookup ON public.sessions (owner_id, owner_type);
CREATE INDEX IF NOT EXISTS idx_sessions_expires ON public.sessions (expires_at);

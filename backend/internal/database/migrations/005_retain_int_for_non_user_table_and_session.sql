-- Sequences for integer IDs
CREATE SEQUENCE IF NOT EXISTS roles_role_id_seq;
CREATE SEQUENCE IF NOT EXISTS permissions_permission_id_seq;
CREATE SEQUENCE IF NOT EXISTS audit_log_audit_id_seq;
CREATE SEQUENCE IF NOT EXISTS payments_payment_id_seq;

-- 1. Correct Executives & Members & Sessions to ensure UUIDs (if not already)
ALTER TABLE IF EXISTS "public"."executives" 
    ALTER COLUMN "executive_id" SET DEFAULT gen_random_uuid(),
    ALTER COLUMN "created_at" SET DEFAULT now();

ALTER TABLE IF EXISTS "public"."members" 
    ALTER COLUMN "member_id" SET DEFAULT gen_random_uuid(),
    ALTER COLUMN "created_at" SET DEFAULT now();

ALTER TABLE IF EXISTS "public"."sessions" 
    ALTER COLUMN "session_id" SET DEFAULT gen_random_uuid(),
    ALTER COLUMN "created_at" SET DEFAULT now();

-- 2. Refactor non-user tables to use INT4
DELETE FROM "public"."role_permissions";
DELETE FROM "public"."audit_log";

-- PERMISSIONS
ALTER TABLE IF EXISTS "public"."permissions" DROP CONSTRAINT IF EXISTS "permissions_pkey" CASCADE;
ALTER TABLE IF EXISTS "public"."permissions" 
    ALTER COLUMN "permission_id" TYPE int4 USING (nextval('permissions_permission_id_seq'::regclass)),
    ALTER COLUMN "created_by" TYPE int4 USING NULL,
    ALTER COLUMN "permission_id" SET DEFAULT nextval('permissions_permission_id_seq'::regclass);
ALTER TABLE IF EXISTS "public"."permissions" ADD PRIMARY KEY ("permission_id");

-- ROLES
ALTER TABLE IF EXISTS "public"."roles" DROP CONSTRAINT IF EXISTS "roles_pkey" CASCADE;
ALTER TABLE IF EXISTS "public"."roles" 
    ALTER COLUMN "role_id" TYPE int4 USING (nextval('roles_role_id_seq'::regclass)),
    ALTER COLUMN "role_id" SET DEFAULT nextval('roles_role_id_seq'::regclass);
ALTER TABLE IF EXISTS "public"."roles" ADD PRIMARY KEY ("role_id");

-- Update EXECUTIVES role_id to int4
ALTER TABLE IF EXISTS "public"."executives" DROP CONSTRAINT IF EXISTS "executive_data_role_id_fkey";
ALTER TABLE IF EXISTS "public"."executives" 
    ALTER COLUMN "role_id" TYPE int4 USING NULL;
ALTER TABLE IF EXISTS "public"."executives" 
    ADD CONSTRAINT "executive_data_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles"("role_id");

-- RE-CREATE ROLE_PERMISSIONS
ALTER TABLE IF EXISTS "public"."role_permissions" 
    ALTER COLUMN "role_id" TYPE int4,
    ALTER COLUMN "permission_id" TYPE int4;
ALTER TABLE IF EXISTS "public"."role_permissions" 
    ADD CONSTRAINT "role_permissions_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles"("role_id") ON DELETE CASCADE,
    ADD CONSTRAINT "role_permissions_permission_id_fkey" FOREIGN KEY ("permission_id") REFERENCES "public"."permissions"("permission_id") ON DELETE CASCADE,
    ADD PRIMARY KEY ("role_id", "permission_id");

-- AUDIT LOG
ALTER TABLE IF EXISTS "public"."audit_log" 
    ALTER COLUMN "audit_id" TYPE int4 USING (nextval('audit_log_audit_id_seq'::regclass)),
    ALTER COLUMN "audit_id" SET DEFAULT nextval('audit_log_audit_id_seq'::regclass),
    ALTER COLUMN "actor_id" TYPE uuid USING NULL,
    ALTER COLUMN "entity_id" TYPE int4 USING NULL;
ALTER TABLE IF EXISTS "public"."audit_log" ADD PRIMARY KEY ("audit_id");

-- PAYMENTS
ALTER TABLE IF EXISTS "public"."payments" DROP CONSTRAINT IF EXISTS "payments_member_id_fkey";
ALTER TABLE IF EXISTS "public"."payments" 
    ALTER COLUMN "payment_id" TYPE int4 USING (nextval('payments_payment_id_seq'::regclass)),
    ALTER COLUMN "payment_id" SET DEFAULT nextval('payments_payment_id_seq'::regclass),
    ALTER COLUMN "approved_by" TYPE uuid USING NULL,
    ALTER COLUMN "member_id" TYPE uuid USING "member_id";
ALTER TABLE IF EXISTS "public"."payments" 
    ADD CONSTRAINT "payments_member_id_fkey" FOREIGN KEY ("member_id") REFERENCES "public"."members"("member_id"),
    ADD CONSTRAINT "payments_approved_by_fkey" FOREIGN KEY ("approved_by") REFERENCES "public"."executives"("executive_id"),
    ADD PRIMARY KEY ("payment_id");

-- RE-SEED DEFAULT DATA
INSERT INTO "public"."roles" ("role_id", "role_name", "description") 
VALUES (1, 'System Administrator', 'Default executive role')
ON CONFLICT (role_id) DO NOTHING;

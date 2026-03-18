ALTER TABLE IF EXISTS "public"."executives" DROP CONSTRAINT IF EXISTS "executive_data_role_id_fkey";
ALTER TABLE IF EXISTS "public"."executives" DROP COLUMN IF EXISTS "role_id";
ALTER TABLE IF EXISTS "public"."executives" ADD COLUMN "role_id" uuid;

DROP TABLE IF EXISTS "public"."role_permissions";
DROP TABLE IF EXISTS "public"."payments";
DROP TABLE IF EXISTS "public"."audit_log";
DROP TABLE IF EXISTS "public"."permissions";
DROP TABLE IF EXISTS "public"."roles";

-- Remove legacy integer sequences after UUID conversion.
DROP SEQUENCE IF EXISTS roles_role_id_seq;
DROP SEQUENCE IF EXISTS permissions_permission_id_seq;
DROP SEQUENCE IF EXISTS audit_log_audit_id_seq;
DROP SEQUENCE IF EXISTS payments_payment_id_seq;

CREATE TABLE "public"."roles" (
    "role_id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "role_name" varchar(100) NOT NULL,
    "description" text,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("role_id")
);

ALTER TABLE IF EXISTS "public"."executives"
    ADD CONSTRAINT "executive_data_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles"("role_id");

CREATE UNIQUE INDEX roles_role_name_key ON public.roles USING btree (role_name);

CREATE TABLE "public"."permissions" (
    "permission_id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "permission_key" varchar(100) NOT NULL,
    "resource" varchar(100) NOT NULL,
    "action" varchar(50) NOT NULL,
    "description" text,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_by" uuid,
    PRIMARY KEY ("permission_id")
);

CREATE UNIQUE INDEX permissions_permission_key_key ON public.permissions USING btree (permission_key);

CREATE TABLE "public"."role_permissions" (
    "role_id" uuid NOT NULL,
    "permission_id" uuid NOT NULL,
    CONSTRAINT "role_permissions_permission_id_fkey" FOREIGN KEY ("permission_id") REFERENCES "public"."permissions"("permission_id") ON DELETE CASCADE,
    CONSTRAINT "role_permissions_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles"("role_id") ON DELETE CASCADE,
    PRIMARY KEY ("role_id", "permission_id")
);

CREATE TABLE "public"."audit_log" (
    "audit_id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "actor_id" uuid,
    "actor_role" varchar(50),
    "action" varchar(100),
    "entity_type" varchar(100),
    "entity_id" uuid,
    "details" jsonb,
    "timestamp" timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("audit_id")
);

CREATE TABLE "public"."payments" (
    "payment_id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "payment_proof_image" text,
    "payment_status" varchar(50) DEFAULT 'pending'::character varying,
    "submission_date" timestamp DEFAULT CURRENT_TIMESTAMP,
    "approval_date" timestamp,
    "approved_by" uuid,
    "member_id" uuid,
    CONSTRAINT "payments_approved_by_fkey" FOREIGN KEY ("approved_by") REFERENCES "public"."executives"("executive_id") ON DELETE SET NULL,
    CONSTRAINT "payments_member_id_fkey" FOREIGN KEY ("member_id") REFERENCES "public"."members"("member_id"),
    PRIMARY KEY ("payment_id")
);

INSERT INTO "public"."roles" ("role_id", "role_name", "description", "created_at") VALUES
('00000000-0000-0000-0000-000000000001', 'System Administrator', 'Default executive role', NOW());

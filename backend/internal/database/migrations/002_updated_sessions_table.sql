-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS roles_role_id_seq;
CREATE SEQUENCE IF NOT EXISTS permissions_permission_id_seq;
CREATE SEQUENCE IF NOT EXISTS payments_payment_id_seq;
CREATE SEQUENCE IF NOT EXISTS audit_log_audit_id_seq;
CREATE SEQUENCE IF NOT EXISTS members_member_id_seq;

DROP TABLE IF EXISTS "public"."roles";
-- Table Definition
CREATE TABLE "public"."roles" (
    "role_id" int4 NOT NULL DEFAULT nextval('roles_role_id_seq'::regclass),
    "role_name" varchar(100) NOT NULL,
    "description" text,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("role_id")
);


-- Indices
CREATE UNIQUE INDEX roles_role_name_key ON public.roles USING btree (role_name);

DROP TABLE IF EXISTS "public"."executives";
-- Table Definition
CREATE TABLE "public"."executives" (
    "executive_id" int4 NOT NULL,
    "name" varchar(150) NOT NULL,
    "email" varchar(150) NOT NULL,
    "student_id" varchar(50) NOT NULL,
    "course" varchar(100),
    "contact_number" varchar(30),
    "role_id" int4,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "last_updated" timestamp DEFAULT CURRENT_TIMESTAMP,
    "password_hash" varchar(60) NOT NULL,
    CONSTRAINT "executives_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles"("role_id") ON DELETE SET NULL,
    PRIMARY KEY ("executive_id")
);


-- Indices
CREATE UNIQUE INDEX executives_email_key ON public.executives USING btree (email);
CREATE UNIQUE INDEX executives_student_id_key ON public.executives USING btree (student_id);

DROP TABLE IF EXISTS "public"."role_permissions";
-- Table Definition
CREATE TABLE "public"."role_permissions" (
    "role_id" int4 NOT NULL,
    "permission_id" int4 NOT NULL,
    CONSTRAINT "role_permissions_permission_id_fkey" FOREIGN KEY ("permission_id") REFERENCES "public"."permissions"("permission_id") ON DELETE CASCADE,
    CONSTRAINT "role_permissions_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles"("role_id") ON DELETE CASCADE,
    PRIMARY KEY ("role_id","permission_id")
);

DROP TABLE IF EXISTS "public"."permissions";
-- Table Definition
CREATE TABLE "public"."permissions" (
    "permission_id" int4 NOT NULL DEFAULT nextval('permissions_permission_id_seq'::regclass),
    "permission_key" varchar(100) NOT NULL,
    "resource" varchar(100) NOT NULL,
    "action" varchar(50) NOT NULL,
    "description" text,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_by" int4,
    PRIMARY KEY ("permission_id")
);


-- Indices
CREATE UNIQUE INDEX permissions_permission_key_key ON public.permissions USING btree (permission_key);

DROP TABLE IF EXISTS "public"."payments";
-- Table Definition
CREATE TABLE "public"."payments" (
    "payment_id" int4 NOT NULL DEFAULT nextval('payments_payment_id_seq'::regclass),
    "member_id" int4 NOT NULL,
    "payment_proof_image" text,
    "payment_status" varchar(50) DEFAULT 'pending'::character varying,
    "submission_date" timestamp DEFAULT CURRENT_TIMESTAMP,
    "approval_date" timestamp,
    "approved_by" int4,
    CONSTRAINT "payments_approved_by_fkey" FOREIGN KEY ("approved_by") REFERENCES "public"."executives"("executive_id"),
    CONSTRAINT "payments_member_id_fkey" FOREIGN KEY ("member_id") REFERENCES "public"."members"("member_id") ON DELETE CASCADE,
    PRIMARY KEY ("payment_id")
);

DROP TABLE IF EXISTS "public"."audit_log";
-- Table Definition
CREATE TABLE "public"."audit_log" (
    "audit_id" int4 NOT NULL DEFAULT nextval('audit_log_audit_id_seq'::regclass),
    "actor_id" int4,
    "actor_role" varchar(50),
    "action" varchar(100),
    "entity_type" varchar(100),
    "entity_id" int4,
    "details" jsonb,
    "timestamp" timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("audit_id")
);

DROP TABLE IF EXISTS "public"."members";
-- Table Definition
CREATE TABLE "public"."members" (
    "member_id" int4 NOT NULL DEFAULT nextval('members_member_id_seq'::regclass),
    "name" varchar(150) NOT NULL,
    "email" varchar(150) NOT NULL,
    "student_id" varchar(50) NOT NULL,
    "course" varchar(100),
    "contact_number" varchar(30),
    "registration_status" varchar(50) DEFAULT 'pending'::character varying,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "last_updated" timestamp DEFAULT CURRENT_TIMESTAMP,
    "password_hash" varchar(60) NOT NULL,
    PRIMARY KEY ("member_id")
);


-- Indices
CREATE UNIQUE INDEX members_email_key ON public.members USING btree (email);
CREATE UNIQUE INDEX members_student_id_key ON public.members USING btree (student_id);

DROP TABLE IF EXISTS "public"."sessions";
-- Table Definition
CREATE TABLE "public"."sessions" (
    "session_id" text NOT NULL,
    "refresh_token_hash" varchar(64),
    "user_agent" text,
    "ip_address" text,
    "expires_at" timestamptz,
    "created_at" timestamptz,
    "revoked_at" timestamptz,
    "owner_id" int4,
    "owner_type" varchar(64)
);


-- Indices
CREATE INDEX idx_sessions_expires ON public.sessions USING btree (expires_at);

INSERT INTO "public"."roles" ("role_id", "role_name", "description", "created_at") VALUES
(1, 'System Administrator', 'Default executive role', '2026-03-17 04:18:35.178826');





INSERT INTO "public"."members" ("member_id", "name", "email", "student_id", "course", "contact_number", "registration_status", "created_at", "last_updated", "password_hash") VALUES
(1, 'Test Member', 'member@example.com', '', NULL, NULL, 'pending', '2026-03-16 11:55:48.644652', '2026-03-16 11:55:48.644652', '$2a$12$YBVX6.jzwunRfM147PNxMuVp0CfIma3iNektdVh1Mdnym195pNjbq'),
(2, 'Member 20260316115955', 'member20260316115955@test.com', 'S20260316115955', NULL, NULL, 'pending', '2026-03-16 11:59:56.256106', '2026-03-16 11:59:56.256106', '$2a$12$jFWcSANWPP3DZ7fFlQ/dCe6MIyznUz3Ykx4hhyBm7vXE6AS83coti'),
(3, 'Aeron Sarondo', 'a3rune@gmail.com', '2022-160122', NULL, NULL, 'pending', '2026-03-17 09:54:42.936389', '2026-03-17 09:54:42.936389', '$2a$12$Y3Matl1jeMck.L3aJf3GduE4p8Hdv.AiDvUugcHAQw9/MpzJccoku');


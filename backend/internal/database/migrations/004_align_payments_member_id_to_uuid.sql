-- Migration: 004_align_payments_member_id_to_uuid.sql
-- Description: Align payments.member_id type with members.member_id UUID refactor.

-- Payments is currently empty, so we can safely recreate it with the correct FK type.
DROP TABLE IF EXISTS "public"."payments";

CREATE TABLE "public"."payments" (
    "payment_id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "member_id" uuid NOT NULL,
    "payment_proof_image" text,
    "payment_status" varchar(50) DEFAULT 'pending'::character varying,
    "submission_date" timestamp DEFAULT CURRENT_TIMESTAMP,
    "approval_date" timestamp,
    "approved_by" uuid,
    CONSTRAINT "payments_approved_by_fkey" FOREIGN KEY ("approved_by") REFERENCES "public"."executives"("executive_id"),
    CONSTRAINT "payments_member_id_fkey" FOREIGN KEY ("member_id") REFERENCES "public"."members"("member_id") ON DELETE CASCADE,
    PRIMARY KEY ("payment_id")
);

CREATE INDEX IF NOT EXISTS "payments_member_id_idx" ON "public"."payments" ("member_id");

DROP TABLE IF EXISTS "public"."members";
-- Table Definition
CREATE TABLE "public"."members" (
    "member_id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "name" text,
    "email" text,
    "student_id" text,
    "course" text,
    "contact_number" text,
    "registration_status" text,
    "last_updated" timestamp,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "password_hash" text,
    PRIMARY KEY ("member_id")
);


-- Indices
CREATE UNIQUE INDEX members_data_pkey ON public.members USING btree (member_id);



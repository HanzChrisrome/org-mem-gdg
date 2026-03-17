-- Add password_hash column to executives table to support login
ALTER TABLE "public"."executives"
ADD COLUMN IF NOT EXISTS "password_hash" varchar(60) NOT NULL DEFAULT '';

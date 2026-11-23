-- +goose Up
CREATE SCHEMA IF NOT EXISTS public;
CREATE SCHEMA IF NOT EXISTS migrations;

-- A database created by the old bootstrap needs a verified baseline.
-- +goose StatementBegin
DO $$
BEGIN
    IF to_regclass('public.users') IS NOT NULL THEN
        RAISE EXCEPTION 'Existing BiOME schema detected. Use a fresh database or baseline the existing schema before migrating.';
    END IF;
END
$$;
-- +goose StatementEnd

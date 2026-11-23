-- +goose Up
CREATE TABLE password_reset_tokens (
    token_hash TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    CONSTRAINT password_reset_tokens_expiry CHECK (expires_at > created_at)
);

CREATE INDEX password_reset_tokens_user_idx ON password_reset_tokens (user_id);

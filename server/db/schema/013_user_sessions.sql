CREATE TABLE user_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    refresh_token_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL CHECK (expires_at > created_at),
    last_used_at TIMESTAMPTZ CHECK (last_used_at >= created_at),
    revoked_at TIMESTAMPTZ CHECK (revoked_at >= created_at),
    user_agent TEXT,
    ip_address INET
);

CREATE INDEX user_sessions_user_id_idx ON user_sessions (user_id);
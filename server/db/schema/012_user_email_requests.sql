-- User account requests and email-change requests.
-- Both flows require explicit email verification; invitation-based accounts bypass these tables.
CREATE TYPE user_account_request_status AS ENUM('pending', 'verified', 'cancelled', 'expired');
CREATE TYPE user_email_change_request_status AS ENUM(
	'pending',
	'verified',
	'applied',
	'cancelled',
	'expired'
);

CREATE TABLE user_account_requests (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	email CITEXT NOT NULL,
	name TEXT NOT NULL,
	motivations TEXT,
	status user_account_request_status NOT NULL DEFAULT 'pending',
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	expires_at TIMESTAMPTZ NOT NULL,
	verified_at TIMESTAMPTZ,
	cancelled_at TIMESTAMPTZ,
	CONSTRAINT user_account_requests_dates CHECK (expires_at > created_at),
	CONSTRAINT user_account_requests_verified_time CHECK (
		verified_at IS NULL
		OR verified_at >= created_at
	),
	CONSTRAINT user_account_requests_cancelled_time CHECK (
		cancelled_at IS NULL
		OR cancelled_at >= created_at
	)
);

CREATE UNIQUE INDEX user_account_requests_one_active_per_email_idx ON user_account_requests (email)
WHERE status = 'pending';

CREATE INDEX user_account_requests_expires_idx ON user_account_requests (expires_at);

-- Tokens table: store only a hash of the raw token for verification.
CREATE TABLE user_account_request_tokens (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	user_account_request_id UUID NOT NULL REFERENCES user_account_requests (id) ON DELETE CASCADE,
	token_hash TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	consumed BOOLEAN NOT NULL DEFAULT false,
	consumed_at TIMESTAMPTZ,
	CONSTRAINT user_account_request_tokens_hash_unique UNIQUE (token_hash),
	CONSTRAINT user_account_request_tokens_consumed_time CHECK (
		(consumed = false)
		OR (consumed_at IS NOT NULL)
	)
);

CREATE INDEX user_account_request_tokens_hash_idx ON user_account_request_tokens (token_hash);
CREATE INDEX user_account_request_tokens_request_idx ON user_account_request_tokens (user_account_request_id);


CREATE TABLE user_email_change_requests (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	email CITEXT NOT NULL,
	status user_email_change_request_status NOT NULL DEFAULT 'pending',
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	expires_at TIMESTAMPTZ NOT NULL,
	verified_at TIMESTAMPTZ,
	applied_at TIMESTAMPTZ,
	cancelled_at TIMESTAMPTZ,
	CONSTRAINT user_email_change_requests_dates CHECK (expires_at > created_at),
	CONSTRAINT user_email_change_requests_verified_time CHECK (
		verified_at IS NULL
		OR verified_at >= created_at
	),
	CONSTRAINT user_email_change_requests_applied_time CHECK (
		applied_at IS NULL
		OR (
			verified_at IS NOT NULL
			AND applied_at >= verified_at
		)
	),
	CONSTRAINT user_email_change_requests_cancelled_time CHECK (
		cancelled_at IS NULL
		OR cancelled_at >= created_at
	)
);

CREATE UNIQUE INDEX user_email_change_requests_one_active_per_user_idx ON user_email_change_requests (user_id)
WHERE status = 'pending';

CREATE UNIQUE INDEX user_email_change_requests_one_active_per_email_idx ON user_email_change_requests (email)
WHERE status = 'pending';

CREATE INDEX user_email_change_requests_expires_idx ON user_email_change_requests (expires_at);
CREATE INDEX user_email_change_requests_user_idx ON user_email_change_requests (user_id);
CREATE INDEX user_email_change_requests_email_idx ON user_email_change_requests (email);

-- Tokens table: store only a hash of the raw token for verification.
CREATE TABLE user_email_change_request_tokens (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	user_email_change_request_id UUID NOT NULL REFERENCES user_email_change_requests (id) ON DELETE CASCADE,
	token_hash TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	consumed BOOLEAN NOT NULL DEFAULT false,
	consumed_by UUID REFERENCES users (id) ON DELETE
	SET NULL,
		consumed_at TIMESTAMPTZ,
		CONSTRAINT user_email_change_request_tokens_hash_unique UNIQUE (token_hash),
		CONSTRAINT user_email_change_request_tokens_consumed_time CHECK (
			(consumed = false)
			OR (consumed_at IS NOT NULL)
		)
);

CREATE INDEX user_email_change_request_tokens_hash_idx ON user_email_change_request_tokens (token_hash);
CREATE INDEX user_email_change_request_tokens_request_idx ON user_email_change_request_tokens (user_email_change_request_id);
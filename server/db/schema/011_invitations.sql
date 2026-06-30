-- Invitations and tokens (separate metadata and token storage)
-- This design separates invitation metadata from token hashes and allows multiple tokens per invitation (rotation).
CREATE TYPE invitation_status AS ENUM('pending', 'redeemed', 'cancelled', 'expired');
CREATE TABLE invitations (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	email CITEXT NOT NULL,
	invitee_name TEXT NOT NULL,
	role user_role NOT NULL DEFAULT 'Visitor',
	message TEXT,
	inviter_id UUID REFERENCES users (id) ON DELETE
	SET NULL,
		status invitation_status NOT NULL DEFAULT 'pending',
		created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		expires_at TIMESTAMPTZ NOT NULL,
		redeemed_at TIMESTAMPTZ,
		revoked_at TIMESTAMPTZ,
		revoked_by UUID REFERENCES users (id) ON DELETE
	SET NULL,
		CONSTRAINT invitations_dates CHECK (expires_at > created_at)
);

-- Only one pending invitation per email
CREATE UNIQUE INDEX invitations_one_active_per_email_idx ON invitations (email)
WHERE status = 'pending';

CREATE INDEX invitations_expires_idx ON invitations (expires_at);
CREATE INDEX invitations_inviter_idx ON invitations (inviter_id);


-- Tokens table: store only a hash of the raw token for verification.
-- Redeeming an invitation implicitly verifies the invited email.
CREATE TABLE invitation_tokens (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	invitation_id UUID NOT NULL REFERENCES invitations (id) ON DELETE CASCADE,
	token_hash TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	consumed BOOLEAN NOT NULL DEFAULT false,
	consumed_by UUID REFERENCES users (id) ON DELETE
	SET NULL,
		consumed_at TIMESTAMPTZ,
		CONSTRAINT invitation_tokens_hash_unique UNIQUE (token_hash),
		CONSTRAINT invitation_tokens_consumed_time CHECK (
			(consumed = false)
			OR (consumed_at IS NOT NULL)
		)
);
CREATE INDEX invitation_tokens_hash_idx ON invitation_tokens (token_hash);
CREATE INDEX invitation_tokens_invitation_idx ON invitation_tokens (invitation_id);
-- Utility hashing function (uses pgcrypto.digest). Applications should pass only the raw token to the hashing function
CREATE OR REPLACE FUNCTION token_sha256 (txt TEXT) RETURNS TEXT LANGUAGE SQL IMMUTABLE AS $$
SELECT encode(digest(txt, 'sha256'), 'hex');
$$;
-- Example usage (application side):
-- INSERT INTO invitations (email, invitee_name, organisation, role, message, expires_at)
--   VALUES (..., now() + interval '7 days') RETURNING id;
-- INSERT INTO invitation_tokens (invitation_id, token_hash)
--   VALUES (:inv_id, token_sha256(:raw_token));
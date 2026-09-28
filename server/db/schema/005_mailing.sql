CREATE TABLE IF NOT EXISTS mailing (
    -- enforce singleton, including inserts that omit id
    id INTEGER PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    mail_from_address TEXT NOT NULL,
    mail_from_name TEXT NOT NULL,
    -- NULL means SMTP has not been configured (or authentication is unused).
    smtp_host TEXT,
    smtp_port INTEGER CHECK (
        smtp_port BETWEEN 1 AND 65535
    ),
    smtp_user TEXT,
    smtp_password TEXT,
    last_updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
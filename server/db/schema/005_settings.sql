CREATE TABLE IF NOT EXISTS settings (
    -- enforce singleton
    id SERIAL PRIMARY KEY CHECK (id = 1),
    app_name TEXT NOT NULL DEFAULT 'BiOME',
    app_subtitle TEXT,
    app_description TEXT,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    account_requests_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    -- feature flags
    molecular_data_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    -- dashboard
    frontpage_message_md TEXT,
    -- contact
    admin_email TEXT NOT NULL
);
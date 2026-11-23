-- +goose Up
CREATE TABLE IF NOT EXISTS service_settings (
    -- enforce singleton, including inserts that omit id
    id INTEGER PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    -- NULL means Geoapify has not been configured.
    geoapify_api_key TEXT
);

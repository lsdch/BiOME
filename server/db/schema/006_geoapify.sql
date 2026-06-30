CREATE TABLE IF NOT EXISTS geoapify_usage (
    id SERIAL PRIMARY KEY,
    usage_date DATE NOT NULL,
    requests_count INTEGER NOT NULL,
    UNIQUE (usage_date)
)
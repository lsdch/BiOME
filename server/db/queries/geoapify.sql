-- name: ListGeoapifyUsage :many 
SELECT id,
    usage_date,
    requests_count
FROM geoapify_usage
ORDER BY usage_date DESC;

-- name: GetTodayGeoapifyUsage :one
SELECT id,
    usage_date,
    requests_count
FROM geoapify_usage
WHERE usage_date = DATE(NOW())
LIMIT 1;

-- name: IncrementTodayGeoapifyUsage :one
INSERT INTO geoapify_usage (usage_date, requests_count)
VALUES (DATE(NOW()), @requests_count) ON CONFLICT (usage_date) DO
UPDATE
SET requests_count = geoapify_usage.requests_count + @requests_count
RETURNING id,
    usage_date,
    requests_count;
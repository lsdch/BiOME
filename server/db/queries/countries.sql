-- name: InsertCountry :one
INSERT INTO countries (code, name, continent, subcontinent)
VALUES ($1, $2, $3, $4) ON CONFLICT (code) DO
UPDATE
SET name = EXCLUDED.name,
    continent = EXCLUDED.continent,
    subcontinent = EXCLUDED.subcontinent
RETURNING *;

-- name: ListCountries :many
SELECT *
FROM countries
ORDER BY name,
    code;
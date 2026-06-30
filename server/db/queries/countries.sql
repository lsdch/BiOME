-- name: InsertCountry :exec
INSERT INTO countries (code, name, continent, subcontinent, geom)
VALUES (
        @code,
        @name,
        @continent,
        @subcontinent,
        ST_SetSRID(ST_GeomFromGeoJSON(@geom::json), 4326)
    );

-- name: ListCountries :many
SELECT name,
    code,
    continent,
    subcontinent
FROM countries
ORDER BY name;

-- name: ListCountriesSummary :many
SELECT c.name,
    c.code,
    c.continent,
    c.subcontinent,
    COUNT(s.id) AS sampling_count,
    COUNT(o.id) AS occurrence_count
FROM countries c
    LEFT JOIN samplings s ON s.site_country_code = c.code
    LEFT JOIN occurrences o ON o.site_country_code = c.code
GROUP BY c.code
ORDER BY c.name;
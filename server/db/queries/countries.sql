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
SELECT *
FROM countries
ORDER BY name;

-- name: ListCountriesSummary :many
SELECT c.name,
    c.code,
    c.continent,
    c.subcontinent,
    COUNT(DISTINCT s.id) AS sampling_count,
    COUNT(DISTINCT o.id) AS occurrence_count
FROM countries c
    LEFT JOIN samplings s ON s.site_country_code = c.code
    LEFT JOIN occurrences o ON o.sampling_id = s.id
GROUP BY c.code
ORDER BY c.name;

-- name: CoordinatesToCountry :many
SELECT c.*
FROM countries c
WHERE ST_Contains(
        c.geom,
        ST_SetSRID(
            ST_Point(
                @latitude::real,
                @longitude::real
            ),
            4326
        )
    )
LIMIT 1;
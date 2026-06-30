

-- name: ListOccurrencesSortByEventDateASC :many
SELECT sqlc.embed(o),
    sqlc.embed(s),
    sqlc.embed(t),
    sqlc.embed(c)
FROM occurrences o
INNER JOIN samplings s ON s.id = o.sampling_id
LEFT JOIN countries c ON c.code = s.site_country_code
INNER JOIN taxa t ON t.id = o.taxon_id
WHERE (
    sqlc.arg(occurrence_ids)::uuid[] IS NULL
    OR o.id = ANY(sqlc.arg(occurrence_ids))
)
ORDER BY s.event_date ASC, t.name ASC, c.name ASC, o.id ASC
LIMIT CASE
    WHEN sqlc.narg(page_limit)::INTEGER IS NULL THEN NULL
    ELSE sqlc.narg(page_limit)::INTEGER
END
OFFSET COALESCE(sqlc.narg(page_offset)::INTEGER, 0);

-- name: ListOccurrencesSortByEventDateDESC :many
SELECT sqlc.embed(o),
    sqlc.embed(s),
    sqlc.embed(t),
    sqlc.embed(c)
FROM occurrences o
INNER JOIN samplings s ON s.id = o.sampling_id
LEFT JOIN countries c ON c.code = s.site_country_code
INNER JOIN taxa t ON t.id = o.taxon_id
WHERE (
    sqlc.arg(occurrence_ids)::uuid[] IS NULL
    OR o.id = ANY(sqlc.arg(occurrence_ids))
)
ORDER BY s.event_date DESC, t.name ASC, c.name ASC, o.id ASC
LIMIT CASE
    WHEN sqlc.narg(page_limit)::INTEGER IS NULL THEN NULL
    ELSE sqlc.narg(page_limit)::INTEGER
END
OFFSET COALESCE(sqlc.narg(page_offset)::INTEGER, 0);

-- name: ListOccurrencesByTaxon :many
SELECT sqlc.embed(o),
    sqlc.embed(s),
    sqlc.embed(t),
    sqlc.embed(c)
FROM occurrences o
INNER JOIN samplings s ON s.id = o.sampling_id
LEFT JOIN countries c ON c.code = s.site_country_code
INNER JOIN taxa t ON t.id = o.taxon_id
WHERE (
    sqlc.arg(occurrence_ids)::uuid[] IS NULL
    OR o.id = ANY(sqlc.arg(occurrence_ids))
)
ORDER BY t.name ASC, s.event_date DESC, c.name ASC, o.id ASC
LIMIT CASE
    WHEN sqlc.narg(page_limit)::INTEGER IS NULL THEN NULL
    ELSE sqlc.narg(page_limit)::INTEGER
END
OFFSET COALESCE(sqlc.narg(page_offset)::INTEGER, 0);

-- name: ListOccurrencesByCountry :many
SELECT sqlc.embed(o),
    sqlc.embed(s),
    sqlc.embed(t),
    sqlc.embed(c)
FROM occurrences o
INNER JOIN samplings s ON s.id = o.sampling_id
LEFT JOIN countries c ON c.code = s.site_country_code
INNER JOIN taxa t ON t.id = o.taxon_id
WHERE (
    sqlc.arg(occurrence_ids)::uuid[] IS NULL
    OR o.id = ANY(sqlc.arg(occurrence_ids))
)
ORDER BY c.name ASC, s.event_date DESC, t.name ASC, o.id ASC
LIMIT CASE
    WHEN sqlc.narg(page_limit)::INTEGER IS NULL THEN NULL
    ELSE sqlc.narg(page_limit)::INTEGER
END
OFFSET COALESCE(sqlc.narg(page_offset)::INTEGER, 0);
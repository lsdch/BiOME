-- name: InsertTaxon :one
INSERT INTO taxa (
        gbif_id,
        name,
        authorship,
        rank,
        status,
        accepted_taxon_id,
        parent_id,
        comments
    ) -- params: gbif_id, name, authorship, rank, status, accepted_scientific_or_name, parent_scientific_or_name, comments
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        (
            SELECT id
            FROM taxa t
            WHERE t.scientific_name = sqlc.arg (synonym_scientific_name)::text
            LIMIT 1
        ), (
            SELECT id
            FROM taxa t
            WHERE t.scientific_name = sqlc.arg (parent_scientific_name)::text
            LIMIT 1
        ), $6
    ) ON CONFLICT ON CONSTRAINT taxon_name_authorship_uidx DO
UPDATE
SET gbif_id = EXCLUDED.gbif_id,
    rank = EXCLUDED.rank,
    status = EXCLUDED.status,
    accepted_taxon_id = EXCLUDED.accepted_taxon_id,
    parent_id = EXCLUDED.parent_id,
    comments = EXCLUDED.comments
RETURNING *;

-- name: GetTaxonByName :one
SELECT *
FROM taxa t
WHERE t.scientific_name = $1
LIMIT 1;

-- name: GetTaxonByID :one
SELECT *
FROM taxa t
WHERE t.id = $1
LIMIT 1;

-- name: DeleteTaxonByScientificName :exec
DELETE FROM taxa
WHERE scientific_name = $1;

-- name: DeleteTaxonByID :exec
DELETE FROM taxa
WHERE id = $1;

-- name: CheckMissingTaxonScientificNamesBulk :many
WITH input AS (
    SELECT unnest(sqlc.arg(names)::text []) AS input_name
)
SELECT input.input_name::text AS input_name
FROM input
WHERE NOT EXISTS (
        SELECT 1
        FROM taxa t
        WHERE t.scientific_name = input.input_name
    )
ORDER BY input.input_name;

-- name: CheckMissingOrDuplicateTaxonNamesBulk :many
-- This query checks for missing or ambiguous taxon names in the database, from a list of input names. 
-- It returns the input names along with matching scientific names from the taxon table. 
-- If a name is missing, it will have an empty array of matching names. 
-- If a name has duplicates, it will return all matching scientific names.
WITH input AS (
    SELECT unnest(sqlc.arg(names)::text []) AS input_name
)
SELECT i.input_name::text AS input_name,
    COALESCE(
        array_agg(t.scientific_name)::text [],
        ARRAY []::text []
    )::text [] AS matching_names
FROM input i
    LEFT JOIN taxa t ON t.name = i.input_name
GROUP BY i.input_name
HAVING COUNT(t.id) != 1
ORDER BY i.input_name;

-- name: GetTaxaByRank :many
SELECT id,
    gbif_id,
    name,
    scientific_name,
    authorship,
    rank,
    status,
    accepted_taxon_id,
    parent_id,
    comments
FROM taxa
WHERE rank = $1
ORDER BY scientific_name ASC,
    name ASC;
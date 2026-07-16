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
    )
VALUES (
        @gbif_id,
        @name,
        @authorship,
        @rank,
        @status,
        @accepted_id,
        @parent_id,
        @comments
    )
RETURNING *;

-- name: GetTaxonByScientificName :one
SELECT *
FROM taxa t
WHERE t.scientific_name = @scientific_name
LIMIT 1;

-- name: GetTaxonByID :one
SELECT *
FROM taxa t
WHERE t.id = @taxon_id
LIMIT 1;

-- name: DeleteTaxonByScientificName :exec
DELETE FROM taxa
WHERE scientific_name = @scientific_name;

-- name: DeleteTaxonByID :exec
DELETE FROM taxa
WHERE id = @taxon_id;

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
SELECT *
FROM taxa
WHERE rank = @rank
ORDER BY scientific_name ASC,
    name ASC;

-- name: GetTaxonLineage :many
With parents AS (
    SELECT ancestor_id AS id
    FROM taxa_closure
    WHERE descendant_id = @taxon_id
)
SELECT t.*
FROM taxa t
    JOIN parents p ON p.id = t.id
ORDER BY t.rank DESC;

-- name: GetTaxonDescendants :many
With descendants AS (
    SELECT descendant_id AS id
    FROM taxa_closure
    WHERE ancestor_id = @taxon_id
        AND depth = 1
)
SELECT t.*
FROM taxa t
    JOIN descendants d ON d.id = t.id
ORDER BY t.name ASC;
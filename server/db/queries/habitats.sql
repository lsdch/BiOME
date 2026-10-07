-- name: InsertHabitat :one
INSERT INTO habitats (label, facet, description)
VALUES (@label, sqlc.narg('facet'), sqlc.narg('description'))
RETURNING *;

-- name: DeleteHabitatByID :exec
DELETE FROM habitats WHERE id = @id;

-- name: ListHabitats :many
SELECT * FROM habitats ORDER BY facet NULLS LAST, label;

-- name: UpdateHabitat :exec
UPDATE habitats
SET label = COALESCE(sqlc.narg('label'), label),
    facet = CASE WHEN @set_facet::boolean THEN sqlc.narg('facet') ELSE facet END,
    description = CASE WHEN @set_description::boolean THEN sqlc.narg('description') ELSE description END
WHERE id = @habitat_id;

-- name: InsertHabitatGroup :one
INSERT INTO habitat_groups (
        id,
        label,
        description,
        parent_habitat_id
    )
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: InsertHabitatInGroup :one
INSERT INTO habitats (
        id,
        label,
        description,
        habitat_group_id
    )
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteHabitatByName :exec
DELETE FROM habitats
WHERE label = $1;

-- name: ListHabitatGroups :many
SELECT id,
    label,
    description,
    parent_habitat_id
FROM habitat_groups
ORDER BY label;

-- name: ListHabitats :many
SELECT id,
    label,
    description,
    habitat_group_id
FROM habitats
ORDER BY habitat_group_id,
    label;
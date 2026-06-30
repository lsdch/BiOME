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

-- name: DeleteHabitatGroupByName :exec
DELETE FROM habitat_groups
WHERE label = $1;

-- name: ListHabitatGroups :many
SELECT g.*
FROM habitat_groups g
ORDER BY label;

-- name: ListHabitats :many
SELECT *
FROM habitats
ORDER BY habitat_group_id,
    label;
-- name: InsertHabitatGroup :one
INSERT INTO habitat_groups (label, exclusive_elements, parent_habitat_id)
VALUES (@label, @exclusive_elements, @parent_habitat_id)
RETURNING *;

-- name: InsertHabitatInGroup :one
INSERT INTO habitats (
        label,
        description,
        habitat_group_id
    )
VALUES (@label, @description, @habitat_group_id)
RETURNING *;

-- name: DeleteHabitatByID :exec
DELETE FROM habitats
WHERE id = @id;

-- name: DeleteHabitatGroup :exec
DELETE FROM habitat_groups
WHERE id = @id;

-- name: ListHabitatGroups :many
SELECT g.*
FROM habitat_groups g
ORDER BY label;

-- name: ListHabitats :many
SELECT *
FROM habitats
ORDER BY habitat_group_id,
    label;

-- name: UpdateHabitat :exec
UPDATE habitats
SET label = COALESCE(sqlc.narg('label'), label),
    description = CASE
        WHEN @set_description::boolean THEN sqlc.narg('description')
        ELSE description
    END
WHERE id = @habitat_id;

-- name: UpdateHabitatGroupInfo :exec
UPDATE habitat_groups
SET label = COALESCE(sqlc.narg('label'), label),
    exclusive_elements = COALESCE(
        sqlc.narg('exclusive_elements'),
        exclusive_elements
    ),
    parent_habitat_id = CASE
        WHEN @set_parent_habitat_id::boolean THEN sqlc.narg('parent_habitat_id')
        ELSE parent_habitat_id
    END
WHERE id = @group_id;
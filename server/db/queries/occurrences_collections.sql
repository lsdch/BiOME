-- name: ListCollectionNames :many
SELECT DISTINCT c.name
FROM occurrence_collections c;

-- name: GetOccurrenceCollections :many
SELECT c.*
FROM occurrence_collections c
WHERE c.occurrence_id = @occurrence_id;

-- name: AddOccurrenceCollection :one
INSERT INTO occurrence_collections (
        occurrence_id,
        name,
        vouchers
    )
VALUES (
        @occurrence_id::ulid,
        @name::TEXT,
        @vouchers::TEXT []
    )
RETURNING *;

-- name: RemoveOccurrenceCollection :exec
DELETE FROM occurrence_collections
WHERE collection_id = @collection_id::uuid;

-- name: UpdateOccurrenceCollection :one
UPDATE occurrence_collections
SET name = @name::TEXT,
    vouchers = @vouchers::TEXT []
WHERE collection_id = @collection_id::uuid
RETURNING *;
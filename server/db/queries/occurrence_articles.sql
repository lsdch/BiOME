-- name: GetOccurrencePublications :many
SELECT a.*
FROM publications a
    JOIN occurrences_publications oa ON oa.publication_id = a.id
WHERE oa.occurrence_id = @occurrence_id;

-- name: AddPublicationToOccurrence :exec
INSERT INTO occurrences_publications (occurrence_id, publication_id)
VALUES (
        @occurrence_id::ulid,
        @publications_id::uuid
    );

-- name: RemovePublicationFromOccurrence :exec
DELETE FROM occurrences_publications
WHERE occurrence_id = @occurrence_id::ulid
    AND publication_id = @publication_id::uuid;
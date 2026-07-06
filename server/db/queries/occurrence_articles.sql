-- name: GetOccurrenceArticles :many
SELECT a.*
FROM articles a
    JOIN occurrences_articles oa ON oa.article_id = a.id
WHERE oa.occurrence_id = @occurrence_id;

-- name: AddArticleToOccurrence :exec
INSERT INTO occurrences_articles (occurrence_id, article_id)
VALUES (
        @occurrence_id::ulid,
        @article_id::uuid
    );

-- name: RemoveArticleFromOccurrence :exec
DELETE FROM occurrences_articles
WHERE occurrence_id = @occurrence_id::ulid
    AND article_id = @article_id::uuid;
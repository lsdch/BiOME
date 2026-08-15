-- name: ListPublications :many
SELECT *
FROM publications
ORDER BY authors [1] ASC,
    year DESC;
-- name: GetPublicationByID :one
SELECT *
FROM publications
WHERE id = @id
LIMIT 1;


-- name: GetPublicationByDOI :one
SELECT *
FROM publications
WHERE doi = @doi
LIMIT 1;

-- name: CreatePublication :one
INSERT INTO publications (
        authors,
        year,
        title,
        journal,
        verbatim,
        doi,
        comments
    )
VALUES (
        @authors,
        @year,
        @title,
        @journal,
        @verbatim,
        @doi,
        @comments
    )
RETURNING *;

-- name: UpdatePublicationByID :one
UPDATE publications
SET authors = COALESCE(sqlc.narg('authors'), authors),
    year = COALESCE(sqlc.narg('year'), year),
    title = CASE
        WHEN @set_title::bool THEN @title
        ELSE title
    END,
    journal = CASE
        WHEN @set_journal::bool THEN @journal
        ELSE journal
    END,
    verbatim = CASE
        WHEN @set_verbatim::bool THEN @verbatim
        ELSE verbatim
    END,
    doi = CASE
        WHEN @set_doi::bool THEN @doi
        ELSE doi
    END,
    comments = CASE
        WHEN @set_comments::bool THEN @comments
        ELSE comments
    END
WHERE id = @id
RETURNING *;


-- name: DeletePublicationByID :one
DELETE FROM publications
WHERE id = @id
RETURNING *;
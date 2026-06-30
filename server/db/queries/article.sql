-- name: ListArticles :many
SELECT *
FROM articles
ORDER BY authors [1] ASC,
    year DESC;
-- name: GetArticleByID :one
SELECT *
FROM articles
WHERE id = @id
LIMIT 1;


-- name: GetArticleByDOI :one
SELECT *
FROM articles
WHERE doi = @doi
LIMIT 1;

-- name: CreateArticle :one
INSERT INTO articles (
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

-- name: UpdateArticleByID :one
UPDATE articles
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


-- name: DeleteArticleByID :one
DELETE FROM articles
WHERE id = @id
RETURNING *;
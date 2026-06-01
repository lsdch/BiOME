-- name: ListArticles :many
SELECT *
FROM articles
ORDER BY authors [1] ASC,
    year DESC,
    code ASC;

-- name: GetArticleByID :one
SELECT *
FROM articles
WHERE id = sqlc.arg(id)
LIMIT 1;

-- name: GetArticleByCode :one
SELECT *
FROM articles
WHERE code = sqlc.arg(code)
LIMIT 1;

-- name: CreateArticle :one
INSERT INTO articles (
        authors,
        year,
        title,
        journal,
        verbatim,
        doi,
        comments,
        code
    )
VALUES (
        sqlc.arg(authors),
        sqlc.arg(year),
        sqlc.narg(title),
        sqlc.narg(journal),
        sqlc.narg(verbatim),
        sqlc.narg(doi),
        sqlc.narg(comments),
        sqlc.arg(code)
    )
RETURNING *;

-- name: UpdateArticleByID :one
UPDATE articles
SET authors = sqlc.arg(authors),
    year = sqlc.arg(year),
    title = sqlc.narg(title),
    journal = sqlc.narg(journal),
    verbatim = sqlc.narg(verbatim),
    doi = sqlc.narg(doi),
    comments = sqlc.narg(comments),
    code = sqlc.arg(code)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateArticleByCode :one
UPDATE articles
SET authors = sqlc.arg(authors),
    year = sqlc.arg(year),
    title = sqlc.narg(title),
    journal = sqlc.narg(journal),
    verbatim = sqlc.narg(verbatim),
    doi = sqlc.narg(doi),
    comments = sqlc.narg(comments),
    code = sqlc.arg(code)
WHERE code = sqlc.arg(old_code)
RETURNING *;


-- name: DeleteArticleByID :one
DELETE FROM articles
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteArticleByCode :one
DELETE FROM articles
WHERE code = sqlc.arg(code)
RETURNING *;
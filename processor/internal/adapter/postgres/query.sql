-- name: GetRepositories :many
SELECT * FROM repositories;

-- name: GetRepository :one
SELECT * FROM repositories
WHERE owner = $1 AND repo = $2;

-- name: UpdateRepository :exec
UPDATE repositories
SET
    description = $3,
    stars = $4,
    forks = $5,
    created_at = $6
WHERE owner = $1 AND repo = $2;

-- name: IsExistsRepository :one
SELECT EXISTS(
    SELECT 1 FROM repositories
    WHERE owner = $1 AND repo = $2
);

-- name: CreateRepository :one
INSERT INTO repositories (
  owner, repo, description, stars, forks, created_at
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: DeleteSubscription :exec
DELETE FROM repositories
WHERE owner = $1 AND repo = $2;

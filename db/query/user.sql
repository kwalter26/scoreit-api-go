-- name: CreateUser :one
INSERT INTO users (username, first_name, last_name, created_at, updated_at)
VALUES ($1, $2, $3, now(), now())
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: UpdateUser :one
UPDATE Users
SET username   = COALESCE(sqlc.narg(username), username),
    first_name = COALESCE(sqlc.narg(first_name), first_name),
    last_name  = COALESCE(sqlc.narg(last_name), last_name),
    updated_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;
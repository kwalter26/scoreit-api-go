-- name: CreateTeam :one
INSERT INTO teams (name)
VALUES ($1)
RETURNING *;

-- name: GetTeam :one
SELECT *
FROM teams
WHERE id = $1
LIMIT 1;

-- name: UpdateTeam :one
UPDATE teams
SET name   = COALESCE(sqlc.narg(name), name),
    updated_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: ListTeams :many
SELECT *
FROM teams
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: DeleteTeam :exec
DELETE
FROM teams
WHERE id = $1;
-- name: CreateRole :one
INSERT INTO user_roles (name, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetRole :one
SELECT *
FROM user_roles
WHERE id = $1;

-- name: GetRoleByName :one
SELECT *
FROM user_roles
WHERE name = $1;

-- name: GetRoles :many
SELECT *
FROM user_roles
WHERE user_id = $1;
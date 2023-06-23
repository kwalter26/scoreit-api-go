-- name: CreateUser :one
INSERT INTO users (username, first_name, last_name, email, hashed_password)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: UpdateUser :one
UPDATE Users
SET username          = COALESCE(sqlc.narg(username), username),
    first_name        = COALESCE(sqlc.narg(first_name), first_name),
    last_name         = COALESCE(sqlc.narg(last_name), last_name),
    email             = COALESCE(sqlc.narg(email), email),
    is_email_verified = COALESCE(sqlc.narg(is_email_verified), is_email_verified),
    hashed_password   = COALESCE(sqlc.narg(hashed_password), hashed_password),
    updated_at        = now()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;
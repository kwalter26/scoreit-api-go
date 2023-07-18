-- name: CreateSession :one
INSERT INTO sessions (id,
                      user_id,
                      refresh_token,
                      user_agent,
                      client_ip,
                      is_blocked,
                      expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetSession :one
SELECT *
FROM sessions
WHERE id = $1
LIMIT 1;

-- name: UpdateSession :one
UPDATE sessions
SET is_blocked = $2
WHERE id = $1
RETURNING *;
-- name: CreateGame :one
INSERT INTO game (home_team_id, away_team_id, home_score, away_score)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetGame :one
SELECT *
FROM game
WHERE id = $1;

-- name: ListGames :many
SELECT *
FROM game g
WHERE (sqlc.narg(home_team_id)::UUID IS NULL OR g.home_team_id = sqlc.narg(home_team_id)::UUID)
  AND (sqlc.narg(away_team_id)::UUID IS NULL OR g.away_team_id = sqlc.narg(away_team_id)::UUID)
ORDER BY g.created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateGame :one
UPDATE game
SET home_score = $1,
    away_score = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING *;
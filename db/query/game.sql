-- name: CreateGame :one
INSERT INTO game (home_team_id, away_team_id, home_score, away_score)
VALUES ($1, $2, $3, $4)
RETURNING *;
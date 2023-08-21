-- name: CreateGameParticipant :one
INSERT INTO game_participant (game_id, player_id, team_id, bat_position)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetGameParticipant :one
SELECT *
FROM game_participant
WHERE id = $1;

-- name: ListGameParticipants :many
SELECT *
FROM game_participant
WHERE game_id = $1
LIMIT $2 OFFSET $3;

-- name: ListGameParticipantsForPlayer :many
SELECT *
FROM game_participant
WHERE player_id = $1
LIMIT $2 OFFSET $3;

-- name: UpdateGameParticipant :one
UPDATE game_participant
SET bat_position = COALESCE(sqlc.narg(bat_position), bat_position)
WHERE id = $1
RETURNING *;

-- name: DeleteGameParticipant :exec
DELETE
FROM game_participant
WHERE id = $1;





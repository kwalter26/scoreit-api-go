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

-- name: ListTeamsOfUser :many
SELECT *
FROM teams
WHERE id IN (
    SELECT team_id
    FROM user_teams
    WHERE user_id = $1
)
LIMIT $2 OFFSET $3;

-- name: ListUsersOfTeam :many
SELECT u.id, u.first_name, u.last_name, ut.primary_position, ut.number, t.name as team_name
FROM users u
 JOIN user_teams ut on u.id = ut.user_id
 JOIN teams t ON ut.team_id = t.id
WHERE u.id IN (
    SELECT user_id
    FROM user_teams
    WHERE ut.team_id = $1
)
LIMIT $2 OFFSET $3;

-- name: AddUserToTeam :one
INSERT INTO user_teams (user_id, team_id,number,primary_position)
VALUES ($1, $2,$3,$4)
RETURNING *;

-- name: DeleteTeam :exec
DELETE
FROM teams
WHERE id = $1;
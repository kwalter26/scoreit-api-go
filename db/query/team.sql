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
SET name       = COALESCE(sqlc.narg(name), name),
    updated_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: ListTeams :many
SELECT *
FROM teams
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: ListTeamMembers :many
SELECT u.id, u.first_name, u.last_name, tm.primary_position, tm.number, t.name as team_name
FROM users u
         JOIN team_members tm on u.id = tm.user_id
         JOIN teams t ON tm.team_id = t.id
WHERE u.id IN (SELECT user_id
               FROM team_members
               WHERE tm.team_id = $1)
LIMIT $2 OFFSET $3;

-- name: ListTeamsOfUser :many
SELECT t.id, t.name
FROM teams t
WHERE t.id IN (SELECT team_id
               FROM team_members tm
               WHERE user_id = $1)
LIMIT $2 OFFSET $3;

-- name: AddTeamMember :one
INSERT INTO team_members (user_id, team_id, number, primary_position)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteTeam :exec
DELETE
FROM teams
WHERE id = $1;
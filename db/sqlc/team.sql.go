// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: team.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const addUserToTeam = `-- name: AddUserToTeam :one
INSERT INTO user_teams (user_id, team_id,number,primary_position)
VALUES ($1, $2,$3,$4)
RETURNING id, number, primary_position, user_id, team_id, created_at
`

type AddUserToTeamParams struct {
	UserID          uuid.UUID `json:"user_id"`
	TeamID          uuid.UUID `json:"team_id"`
	Number          int64     `json:"number"`
	PrimaryPosition string    `json:"primary_position"`
}

func (q *Queries) AddUserToTeam(ctx context.Context, arg AddUserToTeamParams) (UserTeam, error) {
	row := q.db.QueryRowContext(ctx, addUserToTeam,
		arg.UserID,
		arg.TeamID,
		arg.Number,
		arg.PrimaryPosition,
	)
	var i UserTeam
	err := row.Scan(
		&i.ID,
		&i.Number,
		&i.PrimaryPosition,
		&i.UserID,
		&i.TeamID,
		&i.CreatedAt,
	)
	return i, err
}

const createTeam = `-- name: CreateTeam :one
INSERT INTO teams (name)
VALUES ($1)
RETURNING id, name, created_at, updated_at
`

func (q *Queries) CreateTeam(ctx context.Context, name string) (Team, error) {
	row := q.db.QueryRowContext(ctx, createTeam, name)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteTeam = `-- name: DeleteTeam :exec
DELETE
FROM teams
WHERE id = $1
`

func (q *Queries) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteTeam, id)
	return err
}

const getTeam = `-- name: GetTeam :one
SELECT id, name, created_at, updated_at
FROM teams
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetTeam(ctx context.Context, id uuid.UUID) (Team, error) {
	row := q.db.QueryRowContext(ctx, getTeam, id)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listTeams = `-- name: ListTeams :many
SELECT id, name, created_at, updated_at
FROM teams
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListTeamsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListTeams(ctx context.Context, arg ListTeamsParams) ([]Team, error) {
	rows, err := q.db.QueryContext(ctx, listTeams, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Team{}
	for rows.Next() {
		var i Team
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTeamsOfUser = `-- name: ListTeamsOfUser :many
SELECT id, name, created_at, updated_at
FROM teams
WHERE id IN (
    SELECT team_id
    FROM user_teams
    WHERE user_id = $1
)
LIMIT $2 OFFSET $3
`

type ListTeamsOfUserParams struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

func (q *Queries) ListTeamsOfUser(ctx context.Context, arg ListTeamsOfUserParams) ([]Team, error) {
	rows, err := q.db.QueryContext(ctx, listTeamsOfUser, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Team{}
	for rows.Next() {
		var i Team
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsersOfTeam = `-- name: ListUsersOfTeam :many
SELECT u.id, u.first_name, u.last_name, ut.primary_position, ut.number, t.name as team_name
FROM users u
 JOIN user_teams ut on u.id = ut.user_id
 JOIN teams t ON ut.team_id = t.id
WHERE u.id IN (
    SELECT user_id
    FROM user_teams
    WHERE ut.team_id = $1
)
LIMIT $2 OFFSET $3
`

type ListUsersOfTeamParams struct {
	TeamID uuid.UUID `json:"team_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

type ListUsersOfTeamRow struct {
	ID              uuid.UUID `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	PrimaryPosition string    `json:"primary_position"`
	Number          int64     `json:"number"`
	TeamName        string    `json:"team_name"`
}

func (q *Queries) ListUsersOfTeam(ctx context.Context, arg ListUsersOfTeamParams) ([]ListUsersOfTeamRow, error) {
	rows, err := q.db.QueryContext(ctx, listUsersOfTeam, arg.TeamID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListUsersOfTeamRow{}
	for rows.Next() {
		var i ListUsersOfTeamRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.PrimaryPosition,
			&i.Number,
			&i.TeamName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTeam = `-- name: UpdateTeam :one
UPDATE teams
SET name   = COALESCE($1, name),
    updated_at = now()
WHERE id = $2
RETURNING id, name, created_at, updated_at
`

type UpdateTeamParams struct {
	Name sql.NullString `json:"name"`
	ID   uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateTeam(ctx context.Context, arg UpdateTeamParams) (Team, error) {
	row := q.db.QueryRowContext(ctx, updateTeam, arg.Name, arg.ID)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: game_participant.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createGameParticipant = `-- name: CreateGameParticipant :one
INSERT INTO game_participant (game_id, player_id, team_id, bat_position)
VALUES ($1, $2, $3, $4)
RETURNING id, game_id, player_id, team_id, bat_position, created_at, updated_at
`

type CreateGameParticipantParams struct {
	GameID      uuid.UUID `json:"game_id"`
	PlayerID    uuid.UUID `json:"player_id"`
	TeamID      uuid.UUID `json:"team_id"`
	BatPosition int64     `json:"bat_position"`
}

func (q *Queries) CreateGameParticipant(ctx context.Context, arg CreateGameParticipantParams) (GameParticipant, error) {
	row := q.db.QueryRowContext(ctx, createGameParticipant,
		arg.GameID,
		arg.PlayerID,
		arg.TeamID,
		arg.BatPosition,
	)
	var i GameParticipant
	err := row.Scan(
		&i.ID,
		&i.GameID,
		&i.PlayerID,
		&i.TeamID,
		&i.BatPosition,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteGameParticipant = `-- name: DeleteGameParticipant :exec
DELETE
FROM game_participant
WHERE id = $1
`

func (q *Queries) DeleteGameParticipant(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteGameParticipant, id)
	return err
}

const getGameParticipant = `-- name: GetGameParticipant :one
SELECT id, game_id, player_id, team_id, bat_position, created_at, updated_at
FROM game_participant
WHERE id = $1
`

func (q *Queries) GetGameParticipant(ctx context.Context, id uuid.UUID) (GameParticipant, error) {
	row := q.db.QueryRowContext(ctx, getGameParticipant, id)
	var i GameParticipant
	err := row.Scan(
		&i.ID,
		&i.GameID,
		&i.PlayerID,
		&i.TeamID,
		&i.BatPosition,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listGameParticipants = `-- name: ListGameParticipants :many
SELECT id, game_id, player_id, team_id, bat_position, created_at, updated_at
FROM game_participant
WHERE game_id = $1
LIMIT $2 OFFSET $3
`

type ListGameParticipantsParams struct {
	GameID uuid.UUID `json:"game_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

func (q *Queries) ListGameParticipants(ctx context.Context, arg ListGameParticipantsParams) ([]GameParticipant, error) {
	rows, err := q.db.QueryContext(ctx, listGameParticipants, arg.GameID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GameParticipant{}
	for rows.Next() {
		var i GameParticipant
		if err := rows.Scan(
			&i.ID,
			&i.GameID,
			&i.PlayerID,
			&i.TeamID,
			&i.BatPosition,
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

const listGameParticipantsForPlayer = `-- name: ListGameParticipantsForPlayer :many
SELECT id, game_id, player_id, team_id, bat_position, created_at, updated_at
FROM game_participant
WHERE player_id = $1
LIMIT $2 OFFSET $3
`

type ListGameParticipantsForPlayerParams struct {
	PlayerID uuid.UUID `json:"player_id"`
	Limit    int32     `json:"limit"`
	Offset   int32     `json:"offset"`
}

func (q *Queries) ListGameParticipantsForPlayer(ctx context.Context, arg ListGameParticipantsForPlayerParams) ([]GameParticipant, error) {
	rows, err := q.db.QueryContext(ctx, listGameParticipantsForPlayer, arg.PlayerID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GameParticipant{}
	for rows.Next() {
		var i GameParticipant
		if err := rows.Scan(
			&i.ID,
			&i.GameID,
			&i.PlayerID,
			&i.TeamID,
			&i.BatPosition,
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

const updateGameParticipant = `-- name: UpdateGameParticipant :one
UPDATE game_participant
SET bat_position = COALESCE($2, bat_position)
WHERE id = $1
RETURNING id, game_id, player_id, team_id, bat_position, created_at, updated_at
`

type UpdateGameParticipantParams struct {
	ID          uuid.UUID     `json:"id"`
	BatPosition sql.NullInt64 `json:"bat_position"`
}

func (q *Queries) UpdateGameParticipant(ctx context.Context, arg UpdateGameParticipantParams) (GameParticipant, error) {
	row := q.db.QueryRowContext(ctx, updateGameParticipant, arg.ID, arg.BatPosition)
	var i GameParticipant
	err := row.Scan(
		&i.ID,
		&i.GameID,
		&i.PlayerID,
		&i.TeamID,
		&i.BatPosition,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

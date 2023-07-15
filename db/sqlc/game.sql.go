// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: game.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createGame = `-- name: CreateGame :one
INSERT INTO game (home_team_id, away_team_id, home_score, away_score)
VALUES ($1, $2, $3, $4)
RETURNING id, home_team_id, away_team_id, home_score, away_score
`

type CreateGameParams struct {
	HomeTeamID uuid.UUID `json:"home_team_id"`
	AwayTeamID uuid.UUID `json:"away_team_id"`
	HomeScore  int64     `json:"home_score"`
	AwayScore  int64     `json:"away_score"`
}

func (q *Queries) CreateGame(ctx context.Context, arg CreateGameParams) (Game, error) {
	row := q.db.QueryRowContext(ctx, createGame,
		arg.HomeTeamID,
		arg.AwayTeamID,
		arg.HomeScore,
		arg.AwayScore,
	)
	var i Game
	err := row.Scan(
		&i.ID,
		&i.HomeTeamID,
		&i.AwayTeamID,
		&i.HomeScore,
		&i.AwayScore,
	)
	return i, err
}
// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateTeam(ctx context.Context, name string) (Team, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteTeam(ctx context.Context, id uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetTeam(ctx context.Context, id uuid.UUID) (Team, error)
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	ListTeams(ctx context.Context, arg ListTeamsParams) ([]Team, error)
	UpdateTeam(ctx context.Context, arg UpdateTeamParams) (Team, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)

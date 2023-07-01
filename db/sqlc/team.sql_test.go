package db

import (
	"context"
	"database/sql"
	"github.com/kwalter26/scoreit-api-go/util"
	"github.com/stretchr/testify/require"

	"testing"
)

func createRandomTeam(t *testing.T) Team {
	teamName := util.RandomName()

	team, err := testQueries.CreateTeam(context.Background(), teamName)
	require.NoError(t, err)
	require.NotEmpty(t, team)
	require.Equal(t, teamName, team.Name)
	require.NotZero(t, team.CreatedAt)
	require.NotZero(t, team.UpdatedAt)
	return team
}

func TestQueries_CreateTeam(t *testing.T) {
	createRandomTeam(t)
}

func TestQueries_GetTeam(t *testing.T) {
	team := createRandomTeam(t)
	team2, err := testQueries.GetTeam(context.Background(), team.ID)
	require.NoError(t, err)
	require.NotEmpty(t, team2)

	require.Equal(t, team.Name, team2.Name)
	require.Equal(t, team.CreatedAt, team2.CreatedAt)
	require.Equal(t, team.UpdatedAt, team2.UpdatedAt)
}

func TestQueries_UpdateTeam(t *testing.T) {
	team := createRandomTeam(t)

	arg := UpdateTeamParams{
		ID:   team.ID,
		Name: sql.NullString{String: util.RandomName(), Valid: true},
	}
	updatedTeam, err := testQueries.UpdateTeam(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTeam)

	require.Equal(t, arg.Name.String, updatedTeam.Name)
	require.Equal(t, team.CreatedAt, updatedTeam.CreatedAt)
	require.NotEqual(t, team.UpdatedAt, updatedTeam.UpdatedAt)
}

func TestQueries_DeleteTeam(t *testing.T) {
	team := createRandomTeam(t)

	err := testQueries.DeleteTeam(context.Background(), team.ID)
	require.NoError(t, err)

	team2, err := testQueries.GetTeam(context.Background(), team.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, team2)
}

func TestQueries_ListTeams(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTeam(t)
	}

	arg := ListTeamsParams{
		Limit:  5,
		Offset: 0,
	}

	teams, err := testQueries.ListTeams(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, teams, 5)

	for _, team := range teams {
		require.NotEmpty(t, team)
	}
}

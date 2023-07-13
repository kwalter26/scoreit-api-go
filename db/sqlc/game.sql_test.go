package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomGame(t *testing.T) Game {
	team1 := createRandomTeam(t)
	team2 := createRandomTeam(t)
	arg := CreateGameParams{
		HomeTeamID: team1.ID,
		AwayTeamID: team2.ID,
		HomeScore:  0,
		AwayScore:  0,
	}
	game, err := testQueries.CreateGame(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, game)

	require.Equal(t, team1.ID, game.HomeTeamID)
	require.Equal(t, team2.ID, game.AwayTeamID)
	require.Equal(t, int64(0), game.HomeScore)
	require.Equal(t, int64(0), game.AwayScore)

	return game
}

func TestQueries_CreateGame(t *testing.T) {
	createRandomGame(t)
}

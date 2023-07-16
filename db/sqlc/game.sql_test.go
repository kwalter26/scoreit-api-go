package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomGame(t *testing.T, homeTeam *Team, awayTeam *Team) Game {
	var team1 Team
	var team2 Team
	if homeTeam == nil {
		team1 = createRandomTeam(t)
	} else {
		team1 = *homeTeam
	}
	if awayTeam == nil {
		team2 = createRandomTeam(t)
	} else {
		team2 = *awayTeam
	}
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
	createRandomGame(t, nil, nil)
}

func TestQueries_GetGame(t *testing.T) {
	game1 := createRandomGame(t, nil, nil)
	game2, err := testQueries.GetGame(context.Background(), game1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, game2)

	require.Equal(t, game1.ID, game2.ID)
	require.Equal(t, game1.HomeTeamID, game2.HomeTeamID)
	require.Equal(t, game1.AwayTeamID, game2.AwayTeamID)
	require.Equal(t, game1.HomeScore, game2.HomeScore)
	require.Equal(t, game1.AwayScore, game2.AwayScore)
	require.WithinDuration(t, game1.CreatedAt, game2.CreatedAt, time.Second)
	require.WithinDuration(t, game1.UpdatedAt, game2.UpdatedAt, time.Second)
}

func TestQueries_GetGames(t *testing.T) {
	games := make([]Game, 0)
	team := createRandomTeam(t)
	for i := 0; i < 2; i++ {
		games = append(games, createRandomGame(t, &team, nil))
	}
	for i := 0; i < 2; i++ {
		games = append(games, createRandomGame(t, nil, nil))
	}

	fmt.Println("team id: ", team.ID)

	arg := ListGamesParams{
		HomeTeamID: uuid.NullUUID{UUID: team.ID, Valid: true},
		AwayTeamID: uuid.NullUUID{UUID: uuid.Nil, Valid: false},
		Limit:      5,
		Offset:     0,
	}
	newgames, err := testQueries.ListGames(context.Background(), arg)

	for _, game := range newgames {
		fmt.Println(game)
	}
	require.NoError(t, err)
	require.Len(t, newgames, 2)

	for _, game := range newgames {
		require.NotEmpty(t, game)
	}
}

// TestQueries_ListGamesWithAwayTeamId tests the ListGames query with an away team id
func TestQueries_ListGamesWithAwayTeamId(t *testing.T) {
	games := make([]Game, 0)
	team := createRandomTeam(t)
	for i := 0; i < 2; i++ {
		games = append(games, createRandomGame(t, nil, &team))
	}
	for i := 0; i < 2; i++ {
		games = append(games, createRandomGame(t, nil, nil))
	}

	fmt.Println("team id: ", team.ID)

	arg := ListGamesParams{
		HomeTeamID: uuid.NullUUID{UUID: uuid.Nil, Valid: false},
		AwayTeamID: uuid.NullUUID{UUID: team.ID, Valid: true},
		Limit:      5,
		Offset:     0,
	}
	newgames, err := testQueries.ListGames(context.Background(), arg)

	for _, game := range newgames {
		fmt.Println(game)
	}
	require.NoError(t, err)
	require.Len(t, newgames, 2)

	for _, game := range newgames {
		require.NotEmpty(t, game)
	}
}

func TestQueries_GetGamesWithoutTeamId(t *testing.T) {
	games := make([]Game, 0)
	for i := 0; i < 10; i++ {
		games = append(games, createRandomGame(t, nil, nil))
	}
	arg := ListGamesParams{
		HomeTeamID: uuid.NullUUID{UUID: uuid.Nil, Valid: false},
		AwayTeamID: uuid.NullUUID{UUID: uuid.Nil, Valid: false},
		Limit:      5,
		Offset:     0,
	}
	newgames, err := testQueries.ListGames(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, newgames, 5)

	for _, game := range newgames {
		require.NotEmpty(t, game)
	}
}

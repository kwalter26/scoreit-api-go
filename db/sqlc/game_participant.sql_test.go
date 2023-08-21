package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueries_CreateGameParticipant(t *testing.T) {
	homeTeam := createRandomTeam(t)
	awayTeam := createRandomTeam(t)
	game := createRandomGame(t, &homeTeam, &awayTeam)
	player := createRandomUser(t)
	createRandomParticipant(t, homeTeam, awayTeam, game, player)
}

func createRandomParticipant(t *testing.T, homeTeam Team, awayTeam Team, game Game, player User) GameParticipant {

	arg := CreateGameParticipantParams{
		GameID:      game.ID,
		PlayerID:    player.ID,
		TeamID:      homeTeam.ID,
		BatPosition: 1,
	}
	participant, err := testQueries.CreateGameParticipant(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, participant)

	require.Equal(t, arg.GameID, participant.GameID)
	require.Equal(t, arg.PlayerID, participant.PlayerID)
	require.Equal(t, arg.TeamID, participant.TeamID)
	require.Equal(t, arg.BatPosition, participant.BatPosition)

	require.NotZero(t, participant.ID)
	require.NotZero(t, participant.CreatedAt)
	require.NotZero(t, participant.UpdatedAt)
	return participant
}

func TestQueries_GetGameParticipant(t *testing.T) {
	homeTeam := createRandomTeam(t)
	awayTeam := createRandomTeam(t)
	game := createRandomGame(t, &homeTeam, &awayTeam)
	player := createRandomUser(t)
	participant1 := createRandomParticipant(t, homeTeam, awayTeam, game, player)
	participant2, err := testQueries.GetGameParticipant(context.Background(), participant1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, participant2)

	require.Equal(t, participant1.ID, participant2.ID)
	require.Equal(t, participant1.GameID, participant2.GameID)
	require.Equal(t, participant1.PlayerID, participant2.PlayerID)
	require.Equal(t, participant1.TeamID, participant2.TeamID)
	require.Equal(t, participant1.BatPosition, participant2.BatPosition)
	require.WithinDuration(t, participant1.CreatedAt, participant2.CreatedAt, 0)
	require.WithinDuration(t, participant1.UpdatedAt, participant2.UpdatedAt, 0)
}

func TestQueries_ListGameParticipants(t *testing.T) {
	homeTeam := createRandomTeam(t)
	awayTeam := createRandomTeam(t)
	game := createRandomGame(t, &homeTeam, &awayTeam)
	for i := 0; i < 10; i++ {
		player := createRandomUser(t)
		createRandomParticipant(t, homeTeam, awayTeam, game, player)
	}
	arg := ListGameParticipantsParams{
		GameID: game.ID,
		Limit:  5,
		Offset: 0,
	}
	participants, err := testQueries.ListGameParticipants(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, participants, 5)
	for _, participant := range participants {
		require.NotEmpty(t, participant)
	}
}

func TestQueries_ListGameParticipantsByGameId(t *testing.T) {
	homeTeam := createRandomTeam(t)
	awayTeam := createRandomTeam(t)
	player := createRandomUser(t)
	for i := 0; i < 10; i++ {
		game := createRandomGame(t, &homeTeam, &awayTeam)
		createRandomParticipant(t, homeTeam, awayTeam, game, player)
	}
	arg := ListGameParticipantsForPlayerParams{
		PlayerID: player.ID,
		Limit:    5,
		Offset:   0,
	}
	participants, err := testQueries.ListGameParticipantsForPlayer(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, participants, 5)
	for _, participant := range participants {
		require.NotEmpty(t, participant)
	}
}

func TestQueries_UpdateGameParticipant(t *testing.T) {
	homeTeam := createRandomTeam(t)
	awayTeam := createRandomTeam(t)
	game := createRandomGame(t, &homeTeam, &awayTeam)
	player := createRandomUser(t)
	participant1 := createRandomParticipant(t, homeTeam, awayTeam, game, player)

	arg := UpdateGameParticipantParams{
		ID:          participant1.ID,
		BatPosition: sql.NullInt64{Int64: 2, Valid: true},
	}
	participant2, err := testQueries.UpdateGameParticipant(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, participant2)

	require.Equal(t, participant1.ID, participant2.ID)
	require.Equal(t, participant1.GameID, participant2.GameID)
	require.Equal(t, participant1.PlayerID, participant2.PlayerID)
	require.Equal(t, participant1.TeamID, participant2.TeamID)
	require.Equal(t, arg.BatPosition.Int64, participant2.BatPosition)
	require.WithinDuration(t, participant1.CreatedAt, participant2.CreatedAt, 0)
	require.WithinDuration(t, participant1.UpdatedAt, participant2.UpdatedAt, 0)
}

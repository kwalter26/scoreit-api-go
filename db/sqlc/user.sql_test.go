package db

import (
	"context"
	"database/sql"
	"github.com/kwalter26/scoreit-api-go/security"
	"github.com/kwalter26/scoreit-api-go/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := security.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomName(),
		FirstName:      util.RandomName(),
		LastName:       util.RandomName(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestQueries_CreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestQueries_GetUser(t *testing.T) {
	user := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.FirstName, user2.FirstName)
	require.Equal(t, user.LastName, user2.LastName)
}

func TestQueries_UpdateUser(t *testing.T) {
	user := createRandomUser(t)

	arg := UpdateUserParams{
		ID:        user.ID,
		Username:  sql.NullString{String: util.RandomName(), Valid: true},
		FirstName: sql.NullString{String: util.RandomName(), Valid: true},
		LastName:  sql.NullString{String: util.RandomName(), Valid: true},
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, arg.Username.String, updatedUser.Username)
	require.Equal(t, arg.FirstName.String, updatedUser.FirstName)
	require.Equal(t, arg.LastName.String, updatedUser.LastName)
}

func TestQueries_DeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestQueries_ListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}
	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

func TestQueries_UpdateGame(t *testing.T) {
	game := createRandomGame(t, nil, nil)

	arg := UpdateGameParams{
		ID:        game.ID,
		HomeScore: 2,
		AwayScore: 1,
	}
	updatedGame, err := testQueries.UpdateGame(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedGame)

	require.Equal(t, arg.HomeScore, updatedGame.HomeScore)
	require.Equal(t, arg.AwayScore, updatedGame.AwayScore)
}

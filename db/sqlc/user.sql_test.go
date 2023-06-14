package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"scoreit-api-go/util"
	"testing"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:  util.RandomName(),
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
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

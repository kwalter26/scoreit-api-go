package db

import (
	"context"
	"github.com/kwalter26/scoreit-api-go/security"
	"github.com/kwalter26/scoreit-api-go/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomUserTx(t *testing.T) {
	hashedPassword, err := security.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	createUserParams := CreateUserParams{
		Username:       util.RandomName(),
		FirstName:      util.RandomName(),
		LastName:       util.RandomName(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}
	createRoleParams := CreateRoleParams{
		Name: util.RandomName(),
	}

	createUserRequest := CreateUserTxParams{
		CreateUserParams: createUserParams,
		CreateRoleParams: createRoleParams,
	}

	createUserTxResult, err := testStore.CreateUserTx(context.Background(), createUserRequest)
	require.NoError(t, err)
	require.NotEmpty(t, createUserTxResult)

	require.Equal(t, createUserParams.Username, createUserTxResult.User.Username)
	require.Equal(t, createUserParams.FirstName, createUserTxResult.User.FirstName)
	require.Equal(t, createUserParams.LastName, createUserTxResult.User.LastName)
	require.Equal(t, createUserParams.Email, createUserTxResult.User.Email)
	require.Equal(t, createUserParams.HashedPassword, createUserTxResult.User.HashedPassword)

	require.Equal(t, createRoleParams.Name, createUserTxResult.UserRoles[0].Name)
	require.Equal(t, createUserTxResult.User.ID, createUserTxResult.UserRoles[0].UserID)

}

func TestQueries_CreateUserTx(t *testing.T) {
	createRandomUserTx(t)
}

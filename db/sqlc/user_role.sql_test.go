package db

import (
	"context"
	"github.com/kwalter26/scoreit-api-go/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomRole(t *testing.T) UserRole {
	user := createRandomUser(t)
	arg := CreateRoleParams{
		Name:   util.RandomName(),
		UserID: user.ID,
	}

	role, err := testQueries.CreateRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, role)

	require.Equal(t, arg.Name, role.Name)
	require.NotZero(t, role.CreatedAt)
	return role
}

func TestQueriesCreateRole(t *testing.T) {
	createRandomRole(t)
}

func TestQueriesGetRole(t *testing.T) {
	role := createRandomRole(t)
	role2, err := testQueries.GetRole(context.Background(), role.ID)
	require.NoError(t, err)
	require.NotEmpty(t, role2)

	require.Equal(t, role.Name, role2.Name)
}

func TestQueries_GetRoles(t *testing.T) {
	role := createRandomRole(t)
	roles, err := testQueries.GetRoles(context.Background(), role.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, roles)
	require.Equal(t, role.Name, roles[0].Name)
}

func TestQueries_GetRolesByName(t *testing.T) {
	role := createRandomRole(t)
	roles, err := testQueries.GetRolesByName(context.Background(), role.Name)
	require.NoError(t, err)
	require.NotEmpty(t, roles)
	require.Equal(t, role.Name, roles[0].Name)
}

func TestQueries_ListRoles(t *testing.T) {
	for i := 0; i < 3; i++ {
		createRandomRole(t)
	}
	arg := ListRolesParams{
		Limit:  3,
		Offset: 0,
	}
	roles, err := testQueries.ListRoles(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, roles)
	require.Equal(t, len(roles), 3)
}

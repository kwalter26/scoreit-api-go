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

func TestQueries_AddUserToTeam(t *testing.T) {
	user := createRandomUser(t)
	team := createRandomTeam(t)

	arg := AddTeamMemberParams{
		UserID:          user.ID,
		TeamID:          team.ID,
		Number:          util.RandomInt(1, 99),
		PrimaryPosition: string(util.RandomBaseballPosition()),
	}
	userTeam, err := testQueries.AddTeamMember(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, userTeam)
}

func TestQueries_ListUsersInTeam(t *testing.T) {

	team := createRandomTeam(t)
	var users []struct {
		user       User
		teamMember TeamMember
	}
	// Create 5 random users
	for i := 0; i < 5; i++ {
		user := createRandomUser(t)

		arg := AddTeamMemberParams{
			UserID:          user.ID,
			TeamID:          team.ID,
			Number:          util.RandomInt(1, 99),
			PrimaryPosition: string(util.RandomBaseballPosition()),
		}
		userTeam, err := testQueries.AddTeamMember(context.Background(), arg)
		require.NoError(t, err)
		users = append(users, struct {
			user       User
			teamMember TeamMember
		}{user: user, teamMember: userTeam})
	}

	listedUsers, err := testQueries.ListTeamMembers(context.Background(), ListTeamMembersParams{
		TeamID: team.ID,
		Limit:  5,
		Offset: 0,
	})

	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range listedUsers {
		require.NotEmpty(t, user)
		for _, u := range users {
			if user.ID == u.user.ID {
				require.Equal(t, user.ID, u.user.ID)
				require.Equal(t, user.FirstName, u.user.FirstName)
				require.Equal(t, user.LastName, u.user.LastName)
				require.Equal(t, user.Number, u.teamMember.Number)
				require.Equal(t, user.PrimaryPosition, u.teamMember.PrimaryPosition)
				require.Equal(t, user.TeamName, team.Name)
			}
		}
	}
}

func TestQueries_ListTeamsOfUser(t *testing.T) {
	user := createRandomUser(t)
	var teams []struct {
		team       Team
		teamMember TeamMember
	}
	// Create 5 random teams
	for i := 0; i < 5; i++ {
		team := createRandomTeam(t)

		arg := AddTeamMemberParams{
			UserID:          user.ID,
			TeamID:          team.ID,
			Number:          util.RandomInt(1, 99),
			PrimaryPosition: string(util.RandomBaseballPosition()),
		}
		userTeam, err := testQueries.AddTeamMember(context.Background(), arg)
		require.NoError(t, err)
		teams = append(teams, struct {
			team       Team
			teamMember TeamMember
		}{team: team, teamMember: userTeam})
	}

	listedTeams, err := testQueries.ListTeamsOfUser(context.Background(), ListTeamsOfUserParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 0,
	})

	require.NoError(t, err)
	require.Len(t, teams, 5)

	for _, team := range listedTeams {
		require.NotEmpty(t, team)
		for _, u := range teams {
			if team.ID == u.team.ID {
				require.Equal(t, team.ID, u.team.ID)
				require.Equal(t, team.Name, u.team.Name)
			}
		}
	}
}

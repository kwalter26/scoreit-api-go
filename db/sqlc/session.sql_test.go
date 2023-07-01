package db

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/security/token"
	"github.com/kwalter26/scoreit-api-go/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomSession(t *testing.T) Session {
	user := createRandomUser(t)

	// Create refreshToken
	maker, err := token.NewPasetoMaker(util.RandomString(32))

	refreshToken, p, err := maker.CreateToken(user.Username, 4*time.Minute)

	arg := CreateSessionParams{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    "userAgent",
		ClientIp:     "clientIp",
		IsBlocked:    false,
		ExpiresAt:    sql.NullTime{Time: p.ExpireAt, Valid: true},
	}
	session, err := testQueries.CreateSession(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, user.ID, session.UserID)
	require.Equal(t, refreshToken, session.RefreshToken)
	require.Equal(t, "userAgent", session.UserAgent)
	require.Equal(t, "clientIp", session.ClientIp)
	require.Equal(t, false, session.IsBlocked)
	require.Equal(t, p.ExpireAt.UTC(), session.ExpiresAt.Time)

	return session
}

func TestQueries_CreateSession(t *testing.T) {
	createRandomSession(t)
}

func TestQueries_GetSession(t *testing.T) {
	session := createRandomSession(t)
	session2, err := testQueries.GetSession(context.Background(), session.ID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, session.UserID, session2.UserID)
	require.Equal(t, session.RefreshToken, session2.RefreshToken)
	require.Equal(t, session.UserAgent, session2.UserAgent)
	require.Equal(t, session.ClientIp, session2.ClientIp)
	require.Equal(t, session.IsBlocked, session2.IsBlocked)
	require.Equal(t, session.ExpiresAt, session2.ExpiresAt)
	require.Equal(t, session.CreatedAt, session2.CreatedAt)
}

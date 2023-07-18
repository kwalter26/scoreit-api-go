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

	u, err := uuid.NewRandom()
	require.NoError(t, err)

	refreshToken, p, err := maker.CreateToken(u, 4*time.Minute)

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
	require.WithinDurationf(t, p.ExpireAt.UTC(), session.ExpiresAt.Time, time.Second, "expire at should be about 4 minutes")
	require.WithinDurationf(t, p.IssuedAt.UTC(), session.CreatedAt.UTC(), time.Second, "created at should be about now")

	return session
}

func TestQueriesCreateSession(t *testing.T) {
	createRandomSession(t)
}

func TestQueriesGetSession(t *testing.T) {
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

func TestQueriesUpdateSession(t *testing.T) {
	session := createRandomSession(t)

	arg := UpdateSessionParams{
		ID:        session.ID,
		IsBlocked: true,
	}

	session2, err := testQueries.UpdateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, session.ID, session2.ID)
	require.Equal(t, session.UserID, session2.UserID)
	require.Equal(t, session.RefreshToken, session2.RefreshToken)
	require.Equal(t, session.UserAgent, session2.UserAgent)
	require.Equal(t, session.ClientIp, session2.ClientIp)
	require.Equal(t, true, session2.IsBlocked)
	require.Equal(t, session.ExpiresAt, session2.ExpiresAt)
	require.Equal(t, session.CreatedAt, session2.CreatedAt)
}

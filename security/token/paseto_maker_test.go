package token

import (
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/security"
	"github.com/kwalter26/udemy-simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	u, err := uuid.NewRandom()
	require.NoError(t, err)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(u, security.UserRoles, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, u, payload.UserID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpireAt, time.Second)
}

// TestExpiredPasetoToken tests the case when a token is expired
func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	u, err := uuid.NewRandom()
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(u, security.UserRoles, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

// TestInvalidPasetoTokenLength tests the case when a token is invalid
func TestInvalidPasetoTokenLength(t *testing.T) {
	_, err := NewPasetoMaker(util.RandomString(31))
	require.Error(t, err)
}

// TestInvalidToken
func TestInvalidToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	u, err := uuid.NewRandom()
	duration := time.Minute

	token, payload, err := maker.CreateToken(u, security.UserRoles, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	maker2, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	_, err = maker2.VerifyToken(token)
	require.Error(t, err)
}

package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload is the output of the token creation process.
type Payload struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpireAt    time.Time `json:"expire_at"`
	NotBefore   time.Time `json:"not_before"`
	Audience    string
	Issuer      string
	Subject     string
	Permissions []string
}

func (p *Payload) Valid() error {
	now := time.Now()
	if now.After(p.ExpireAt) {
		return ErrExpiredToken
	}
	return nil
}

// NewPayload creates a new payload with a specific username and duration.
func NewPayload(userID uuid.UUID, permissions []string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:          tokenID,
		UserID:      userID,
		IssuedAt:    time.Now(),
		ExpireAt:    time.Now().Add(duration),
		NotBefore:   time.Now().Add(duration - 5*time.Minute),
		Audience:    "scoreit-app",
		Issuer:      "scoreit-project",
		Subject:     "user-token",
		Permissions: permissions,
	}, nil
}

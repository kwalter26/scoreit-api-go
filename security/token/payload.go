package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload is the output of the token creation process.
type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpireAt  time.Time `json:"expire_at"`
	NotBefore time.Time `json:"not_before"`
	Audience  string
	Issuer    string
	Subject   string
}

func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.ExpireAt), nil
}

func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.NotBefore), nil
}

func (p *Payload) GetIssuer() (string, error) {
	return p.Issuer, nil
}

func (p *Payload) GetSubject() (string, error) {
	return p.Subject, nil
}

func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{p.Audience}, nil
}

func (p *Payload) Valid() error {
	now := time.Now()
	if now.After(p.ExpireAt) {
		return ErrExpiredToken
	}
	return nil
}

// NewPayload creates a new payload with a specific username and duration.
func NewPayload(userID uuid.UUID, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Payload{
		ID:        tokenID,
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpireAt:  time.Now().Add(duration),
		NotBefore: time.Now().Add(duration - 5*time.Minute),
		Audience:  "simplebank-app",
		Issuer:    "simplebank-project",
		Subject:   "user-token",
	}, nil
}

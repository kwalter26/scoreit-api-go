package token

import (
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/security"
	"time"
)

type Maker interface {
	// CreateToken creates a new token for a specific username and duration.
	// The duration parameter is represented in hours.
	CreateToken(userId uuid.UUID, permissions []security.Role, duration time.Duration) (string, *Payload, error)
	// VerifyToken checks if the token is valid or not.
	// If the token is valid, it returns the Payload and nil.
	VerifyToken(token string) (*Payload, error)
}

package token

import "time"

type Maker interface {
	// CreateToken creates a new token for a specific username and duration.
	// The duration parameter is represented in hours.
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	// VerifyToken checks if the token is valid or not.
	// If the token is valid, it returns the Payload and nil.
	VerifyToken(token string) (*Payload, error)
}

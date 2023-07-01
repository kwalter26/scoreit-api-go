// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID    `json:"id"`
	UserID       uuid.UUID    `json:"user_id"`
	RefreshToken string       `json:"refresh_token"`
	UserAgent    string       `json:"user_agent"`
	ClientIp     string       `json:"client_ip"`
	IsBlocked    bool         `json:"is_blocked"`
	ExpiresAt    sql.NullTime `json:"expires_at"`
	CreatedAt    time.Time    `json:"created_at"`
}

type Team struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID                uuid.UUID `json:"id"`
	Username          string    `json:"username"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	IsEmailVerified   bool      `json:"is_email_verified"`
	HashedPassword    string    `json:"hashed_password"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type UserTeam struct {
	ID              uuid.UUID `json:"id"`
	Number          int64     `json:"number"`
	PrimaryPosition string    `json:"primary_position"`
	UserID          uuid.UUID `json:"user_id"`
	TeamID          uuid.UUID `json:"team_id"`
	CreatedAt       time.Time `json:"created_at"`
}

type VerifyEmail struct {
	ID         uuid.UUID `json:"id"`
	UserID     string    `json:"user_id"`
	Email      string    `json:"email"`
	SecretCode string    `json:"secret_code"`
	IsUsed     bool      `json:"is_used"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}

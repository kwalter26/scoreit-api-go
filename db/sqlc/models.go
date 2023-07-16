// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Atbat struct {
	ID         uuid.UUID     `json:"id"`
	InningID   uuid.NullUUID `json:"inning_id"`
	BatterID   uuid.NullUUID `json:"batter_id"`
	PitcherID  uuid.NullUUID `json:"pitcher_id"`
	Balls      int64         `json:"balls"`
	Strikes    int64         `json:"strikes"`
	InitBases  int64         `json:"init_bases"`
	TotalBases int64         `json:"total_bases"`
	Out        bool          `json:"out"`
}

type Game struct {
	ID         uuid.UUID `json:"id"`
	HomeTeamID uuid.UUID `json:"home_team_id"`
	AwayTeamID uuid.UUID `json:"away_team_id"`
	HomeScore  int64     `json:"home_score"`
	AwayScore  int64     `json:"away_score"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type GameParticipant struct {
	ID          uuid.UUID     `json:"id"`
	GameID      uuid.NullUUID `json:"game_id"`
	PlayerID    uuid.NullUUID `json:"player_id"`
	HomeTeam    bool          `json:"home_team"`
	BatPosition int64         `json:"bat_position"`
}

type GameStat struct {
	ID      uuid.UUID      `json:"id"`
	AtbatID uuid.NullUUID  `json:"atbat_id"`
	Type    sql.NullString `json:"type"`
}

type Inning struct {
	ID          uuid.UUID     `json:"id"`
	GameID      uuid.NullUUID `json:"game_id"`
	Number      int64         `json:"number"`
	HomeRuns    int64         `json:"home_runs"`
	HomeHits    int64         `json:"home_hits"`
	HomeErrors  int64         `json:"home_errors"`
	HomeLastBat uuid.UUID     `json:"home_last_bat"`
	AwayRuns    int64         `json:"away_runs"`
	AwayHits    int64         `json:"away_hits"`
	AwayErrors  int64         `json:"away_errors"`
	AwayLastBat uuid.UUID     `json:"away_last_bat"`
}

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

type TeamMember struct {
	ID              uuid.UUID `json:"id"`
	Number          int64     `json:"number"`
	PrimaryPosition string    `json:"primary_position"`
	UserID          uuid.UUID `json:"user_id"`
	TeamID          uuid.UUID `json:"team_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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

type VerifyEmail struct {
	ID         uuid.UUID `json:"id"`
	UserID     string    `json:"user_id"`
	Email      string    `json:"email"`
	SecretCode string    `json:"secret_code"`
	IsUsed     bool      `json:"is_used"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}

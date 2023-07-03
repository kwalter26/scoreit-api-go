package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/api/helpers"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"github.com/kwalter26/scoreit-api-go/security"
	"time"
)

// CreateUserRequest represents a request to create a new user.
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required,alphanum,min=3,max=40"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required,min=6,max=40"`
	Email     string `json:"email" binding:"required,email"`
}

// CreateUserResponse represents a response from a create user request.
type CreateUserResponse struct {
	Username          string `json:"username"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	PasswordChangedAt string `json:"password_changed_at"`
	CreatedAt         string `json:"created_at"`
}

// NewUserResponse creates a new CreateUserResponse from a db.User.
func NewUserResponse(user db.User) CreateUserResponse {
	return CreateUserResponse{
		Username:          user.Username,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt.String(),
		CreatedAt:         user.CreatedAt.String(),
	}
}

// CreateNewUser creates a new user account.
func (s *Server) CreateNewUser(context *gin.Context) {
	var req CreateUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	hashedPassword, err := security.HashPassword(req.Password)
	arg := db.CreateUserParams{
		Username:       req.Username,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}

	user, err := s.store.CreateUser(context, arg)
	if err != nil {
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	rsp := NewUserResponse(user)
	context.JSON(200, rsp)
}

// LoginUserRequest represents a request to login a user.
type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=40"`
	Password string `json:"password" binding:"required,min=6,max=40"`
}

// LoginUserResponse represents a response from a login user request.
type LoginUserResponse struct {
	SessionID             uuid.UUID          `json:"session_id"`
	AccessToken           string             `json:"access_token"`
	AccessTokenExpiresAt  time.Time          `json:"access_token_expires_at"`
	RefreshToken          string             `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time          `json:"refresh_token_expires_at"`
	User                  CreateUserResponse `json:"user"`
}

// LoginUser logs in a user.
func (s *Server) LoginUser(context *gin.Context) {
	var req LoginUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	user, err := s.store.GetUserByUsername(context, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(400, helpers.ErrorResponse(err))
			return
		}
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	if err := security.CheckPassword(req.Password, user.HashedPassword); err != nil {
		context.JSON(401, helpers.ErrorResponse(err))
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.RefreshTokenDuration)
	if err != nil {
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	session, err := s.store.CreateSession(context, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    context.GetHeader("User-Agent"),
		ClientIp:     context.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    sql.NullTime{Time: refreshPayload.ExpireAt, Valid: true},
	})

	rsp := LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpireAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpireAt,
		User:                  NewUserResponse(user),
	}

	context.JSON(200, rsp)
}

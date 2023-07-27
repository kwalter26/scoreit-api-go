package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/api/helpers"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"github.com/kwalter26/scoreit-api-go/security"
	"net/http"
	"time"
)

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
			context.JSON(http.StatusNotFound, helpers.ErrorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err))
		return
	}

	if err := security.CheckPassword(req.Password, user.HashedPassword); err != nil {
		context.JSON(http.StatusUnauthorized, helpers.ErrorResponse(err))
		return
	}

	roles, err := s.store.GetRoles(context, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, helpers.ErrorResponse(err))
			return
		}
	}

	var permissions []security.Role
	for _, role := range roles {
		permissions = append(permissions, security.Role(role.Name))
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(user.ID, permissions, s.config.AccessTokenDuration)
	if err != nil {
		context.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.ID, []security.Role{}, s.config.RefreshTokenDuration)
	if err != nil {
		context.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err))
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

	context.JSON(http.StatusOK, rsp)
}

// RefreshTokenRequest represents a request to refresh a token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse represents a response from a refresh token request.
type RefreshTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// RefreshToken refreshes a token.
func (s *Server) RefreshToken(context *gin.Context) {
	var req RefreshTokenRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, helpers.ErrorResponse(err))
		return
	}

	// Verify refresh token
	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		context.JSON(http.StatusUnauthorized, helpers.ErrorResponse(err))
		return
	}

	// Get the session from the database
	session, err := s.store.GetSession(context, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, helpers.ErrorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err))
		return
	}

	// check if session is blocked
	if session.IsBlocked {
		err := fmt.Errorf("session is blocked")
		context.JSON(http.StatusUnauthorized, helpers.ErrorResponse(err))
		return
	}

	// check if session username matches
	if session.UserID != refreshPayload.UserID {
		err := fmt.Errorf("incorrect session user")
		context.JSON(http.StatusUnauthorized, helpers.ErrorResponse(err))
		return
	}

	// check if session token matches request token
	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatch session token")
		context.JSON(http.StatusUnauthorized, helpers.ErrorResponse(err))
		return
	}

	// check if session token is expired
	if time.Now().After(session.ExpiresAt.Time) {
		err := fmt.Errorf("refresh token expired")
		context.JSON(http.StatusUnauthorized, helpers.ErrorResponse(err))
		return
	}

	// Get the user from the database
	permissions := new([]security.Role)

	// Create a new access token
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(refreshPayload.UserID, *permissions, s.config.AccessTokenDuration)
	if err != nil {
		context.JSON(http.StatusUnauthorized, helpers.ErrorResponse(err))
		return
	}

	// Return the access token
	rsp := RefreshTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpireAt,
	}

	context.JSON(http.StatusOK, rsp)
}

// LogoutUserRequest represents a request to logout a user.
type LogoutUserRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutUser represents a request to logout a user.
func (s *Server) LogoutUser(context *gin.Context) {
	var req LogoutUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, helpers.ErrorResponse(err))
		return
	}

	// Verify refresh token
	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		context.JSON(http.StatusUnauthorized, helpers.ErrorResponse(err))
		return
	}

	_, err = s.store.UpdateSession(context, db.UpdateSessionParams{
		ID:        refreshPayload.ID,
		IsBlocked: true,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, helpers.ErrorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err))
		return
	}

	context.JSON(http.StatusOK, nil)
}

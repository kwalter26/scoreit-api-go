package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"time"
)

// ListTeamsRequest represents a request to list teams.
type ListTeamsRequest struct {
	Limit  int32 `form:"limit,default=5"`
	Offset int32 `form:"offset,default=0"`
}

// TeamResponse represents a response from a team request.
type TeamResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (s *Server) ListTeams(context *gin.Context) {
	var req ListTeamsRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(400, errorResponse(err))
		return
	}

	arg := db.ListTeamsParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	teams, err := s.store.ListTeams(context, arg)
	if err != nil {
		context.JSON(500, errorResponse(err))
		return
	}

	var list []TeamResponse
	for _, team := range teams {
		list = append(list, NewTeamResponse(team))
	}

	context.JSON(200, list)
}

// NewTeamResponse creates a new TeamResponse from a db.Team.
func NewTeamResponse(team db.Team) TeamResponse {
	return TeamResponse{
		ID:   team.ID,
		Name: team.Name,
	}
}

// AddUserToTeamRequest represents a request to add a user to a team.
type AddUserToTeamRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	TeamID uuid.UUID `json:"team_id" binding:"required"`
}

// AddUserToTeamResponse represents a response from a add user to team request.
type AddUserToTeamResponse struct {
	UserID    uuid.UUID     `json:"user_id"`
	TeamID    uuid.UUID     `json:"team_id"`
	CreatedAt time.Duration `json:"created_at"`
	UpdateAt  time.Duration `json:"updated_at"`
}

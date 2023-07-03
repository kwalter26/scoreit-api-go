package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/api/helpers"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
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
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	arg := db.ListTeamsParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	teams, err := s.store.ListTeams(context, arg)
	if err != nil {
		context.JSON(500, helpers.ErrorResponse(err))
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

// AddTeamMemberRequest represents a request to add a user to a team.
type AddTeamMemberRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	TeamID uuid.UUID `json:"team_id" binding:"required"`
}

// AddTeamMember adds a user to a team.
func (s *Server) AddTeamMember(context *gin.Context) {
	var req AddTeamMemberRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	arg := db.AddTeamMemberParams{
		UserID: req.UserID,
		TeamID: req.TeamID,
	}

	userTeam, err := s.store.AddTeamMember(context, arg)
	if err != nil {
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	context.JSON(200, userTeam)
}

// GetTeamRequest represents a request to get a team.
type GetTeamRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// GetTeam gets a team.
func (s *Server) GetTeam(context *gin.Context) {
	var req GetTeamRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	team, err := s.store.GetTeam(context, id)
	if err != nil {
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	context.JSON(200, team)
}

// CreateTeamRequest represents a request to create a team.
type CreateTeamRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateTeam creates a team.
func (s *Server) CreateTeam(context *gin.Context) {
	var req CreateTeamRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	team, err := s.store.CreateTeam(context, req.Name)
	if err != nil {
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	context.JSON(200, team)
}

// ListTeamMembersRequest represents a request to list team members.
type ListTeamMembersRequest struct {
	ID     string `uri:"id" binding:"required,uuid"`
	Limit  int32  `form:"limit,default=5"`
	Offset int32  `form:"offset,default=0"`
}

// ListTeamMembers lists team members.
func (s *Server) ListTeamMembers(context *gin.Context) {
	var req ListTeamMembersRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	arg := db.ListTeamMembersParams{
		TeamID: id,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	members, err := s.store.ListTeamMembers(context, arg)
	if err != nil {
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	context.JSON(200, members)
}

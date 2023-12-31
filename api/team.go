package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/api/helpers"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
)

// ListTeamsRequest represents a request to list teams.
type ListTeamsRequest struct {
	PageSize int32 `form:"page_size,default=5" binding:"max=100,min=1"`
	PageID   int32 `form:"page_id,default=1" binding:"min=1"`
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
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	teams, err := s.store.ListTeams(context, arg)
	if err != nil {
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	var list = make([]TeamResponse, 0, len(teams))
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
	TeamID string `uri:"id" binding:"required,uuid"`
	UserID string `uri:"user_id" binding:"required,uuid"`
}

// AddTeamRequestBody represents the body of a request to add a team.
type AddTeamRequestBody struct {
	Number          int64  `json:"number" binding:"required,min=0"`
	PrimaryPosition string `json:"primary_position" binding:"required"`
}

// AddTeamMember adds a user to a team.
func (s *Server) AddTeamMember(context *gin.Context) {
	var req AddTeamMemberRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	var body AddTeamRequestBody
	if err := context.ShouldBindJSON(&body); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	userId := uuid.MustParse(req.UserID)

	teamId := uuid.MustParse(req.TeamID)

	arg := db.AddTeamMemberParams{
		UserID:          userId,
		TeamID:          teamId,
		Number:          body.Number,
		PrimaryPosition: body.PrimaryPosition,
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

	id := uuid.MustParse(req.ID)

	team, err := s.store.GetTeam(context, id)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(404, helpers.ErrorResponse(err))
			return
		}
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
	ID string `uri:"id" binding:"required,uuid"`
}

// ListTeamMembersRequestQuery represents a request to list team members.
type ListTeamMembersRequestQuery struct {
	PageSize int32 `form:"page_size,default=5" binding:"max=100,min=1"`
	PageId   int32 `form:"page_id,default=1" binding:"min=1"`
}

// ListTeamMembers lists team members.
func (s *Server) ListTeamMembers(context *gin.Context) {
	var req ListTeamMembersRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	var query ListTeamMembersRequestQuery
	if err := context.ShouldBindQuery(&query); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	id := uuid.MustParse(req.ID)

	arg := db.ListTeamMembersParams{
		TeamID: id,
		Limit:  query.PageSize * (query.PageId),
		Offset: query.PageId - 1,
	}

	members, err := s.store.ListTeamMembers(context, arg)
	if err != nil {
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	context.JSON(200, members)
}

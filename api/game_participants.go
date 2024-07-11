package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/api/helpers"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"net/http"
)

type CreateGameUriRequest struct {
	GameId string `uri:"id" binding:"required,uuid"`
}

type CreateGameParticipantRequest struct {
	PlayerID    string `json:"player_id" binding:"required,uuid"`
	TeamID      string `json:"team_id" binding:"required,uuid"`
	BatPosition int64  `json:"bat_position" binding:"required"`
}

type CreateGameParticipantResponse struct {
	ID          string `json:"id"`
	GameID      string `json:"game_id"`
	PlayerID    string `json:"player_id"`
	TeamID      string `json:"team_id"`
	BatPosition int64  `json:"bat_position"`
}

func (s *Server) CreateGameParticipant(context *gin.Context) {
	var uriReq CreateGameUriRequest
	if err := context.ShouldBindUri(&uriReq); err != nil {
		context.JSON(http.StatusBadRequest, helpers.ErrorResponse(err))
		return
	}

	var req CreateGameParticipantRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, helpers.ErrorResponse(err))
		return
	}

	gameID := uuid.MustParse(uriReq.GameId)
	playerID := uuid.MustParse(req.PlayerID)
	teamID := uuid.MustParse(req.TeamID)
	batPosition := req.BatPosition

	arg := db.CreateGameParticipantParams{
		GameID:      gameID,
		PlayerID:    playerID,
		TeamID:      teamID,
		BatPosition: batPosition,
	}

	participant, err := s.store.CreateGameParticipant(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err))
		return
	}

	context.JSON(http.StatusOK, CreateGameParticipantResponse{
		ID:          participant.ID.String(),
		GameID:      participant.GameID.String(),
		PlayerID:    participant.PlayerID.String(),
		TeamID:      participant.TeamID.String(),
		BatPosition: participant.BatPosition,
	})
}

type GetGameParticipantRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type GetGameParticipantResponse struct {
	ID          string `json:"id"`
	GameID      string `json:"game_id"`
	PlayerID    string `json:"player_id"`
	TeamID      string `json:"team_id"`
	BatPosition int64  `json:"bat_position"`
}

func (s *Server) GetGameParticipant(context *gin.Context) {
	var req GetGameParticipantRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, helpers.ErrorResponse(err))
		return
	}

	id := uuid.MustParse(req.ID)

	participant, err := s.store.GetGameParticipant(context, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.JSON(http.StatusNotFound, helpers.ErrorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err))
		return
	}

	context.JSON(http.StatusOK, GetGameParticipantResponse{
		ID:          participant.ID.String(),
		GameID:      participant.GameID.String(),
		PlayerID:    participant.PlayerID.String(),
		TeamID:      participant.TeamID.String(),
		BatPosition: participant.BatPosition,
	})
}

type ListGameParticipantsUriRequest struct {
	GameID string `uri:"id" binding:"required,uuid"`
}

type ListGameParticipantsRequest struct {
	PageSize int32 `form:"page_size,default=25" binding:"number,min=1,max=100"`
	PageID   int32 `form:"page_id,default=1" binding:"number,min=1"`
}

type ListGameParticipantsResponse struct {
	Participants []GetGameParticipantResponse `json:"participants"`
}

func (s *Server) ListGameParticipants(context *gin.Context) {
	var uriReq ListGameParticipantsUriRequest
	if err := context.ShouldBindUri(&uriReq); err != nil {
		context.JSON(http.StatusBadRequest, helpers.ErrorResponse(err))
		return
	}

	var req ListGameParticipantsRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, helpers.ErrorResponse(err))
		return
	}

	gameID := uuid.MustParse(uriReq.GameID)

	participants, err := s.store.ListGameParticipants(context, db.ListGameParticipantsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
		GameID: gameID,
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err))
		return
	}

	var rsp ListGameParticipantsResponse
	for _, participant := range participants {
		rsp.Participants = append(rsp.Participants, GetGameParticipantResponse{
			ID:          participant.ID.String(),
			GameID:      participant.GameID.String(),
			PlayerID:    participant.PlayerID.String(),
			TeamID:      participant.TeamID.String(),
			BatPosition: participant.BatPosition,
		})
	}

	context.JSON(http.StatusOK, rsp)
}

type ListGameParticipantsForPlayerRequest struct {
	PlayerID string `form:"player_id" binding:"required,uuid"`
	PageSize int32  `form:"page_size" binding:"omitempty,number,min=1,max=10"`
	PageID   int32  `form:"page_id" binding:"omitempty,number,min=1"`
}

type ListGameParticipantsForPlayerResponse struct {
	Participants []GetGameParticipantResponse `json:"participants"`
}

func (s *Server) ListGameParticipantsForPlayer(context *gin.Context) {
	var req ListGameParticipantsForPlayerRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, helpers.ErrorResponse(err))
		return
	}

	playerID := uuid.MustParse(req.PlayerID)

	participants, err := s.store.ListGameParticipantsForPlayer(context, db.ListGameParticipantsForPlayerParams{
		Limit:    req.PageSize,
		Offset:   (req.PageID - 1) * req.PageSize,
		PlayerID: playerID,
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err))
		return
	}

	var rsp ListGameParticipantsForPlayerResponse
	for _, participant := range participants {
		rsp.Participants = append(rsp.Participants, GetGameParticipantResponse{
			ID:          participant.ID.String(),
			GameID:      participant.GameID.String(),
			PlayerID:    participant.PlayerID.String(),
			TeamID:      participant.TeamID.String(),
			BatPosition: participant.BatPosition,
		})
	}

	context.JSON(http.StatusOK, rsp)
}

type UpdateGameParticipantRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type UpdateGameParticipantResponse struct {
	ID          string `json:"id"`
	GameID      string `json:"game_id"`
	PlayerID    string `json:"player_id"`
	TeamID      string `json:"team_id"`
	BatPosition int64  `json:"bat_position"`
}

func (s *Server) UpdateGameParticipant(context *gin.Context) {
	var req UpdateGameParticipantRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, helpers.ErrorResponse(err))
		return
	}

	id := uuid.MustParse(req.ID)

	participant, err := s.store.GetGameParticipant(context, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.JSON(http.StatusNotFound, helpers.ErrorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err))
		return
	}

	context.JSON(http.StatusOK, UpdateGameParticipantResponse{
		ID:          participant.ID.String(),
		GameID:      participant.GameID.String(),
		PlayerID:    participant.PlayerID.String(),
		TeamID:      participant.TeamID.String(),
		BatPosition: participant.BatPosition,
	})
}

type DeleteGameParticipantRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (s *Server) DeleteGameParticipant(context *gin.Context) {
	var req DeleteGameParticipantRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, helpers.ErrorResponse(err))
		return
	}

	id := uuid.MustParse(req.ID)

	err := s.store.DeleteGameParticipant(context, id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, helpers.ErrorResponse(err))
		return
	}

	context.Status(http.StatusOK)
}

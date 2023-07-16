package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/api/helpers"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"github.com/lib/pq"
)

// CreateGameRequest defines the request body for NewGameHandler.
type CreateGameRequest struct {
	HomeTeamID string `json:"home_team_id" binding:"required,uuid"`
	AwayTeamID string `json:"away_team_id" binding:"required,uuid"`
}

// CreateGameResponse defines the response body for NewGameHandler.
type CreateGameResponse struct {
	ID         string `json:"id"`
	HomeTeamID string `json:"home_team_id"`
	AwayTeamID string `json:"away_team_id"`
	HomeScore  int64  `json:"home_score"`
	AwayScore  int64  `json:"away_score"`
}

// CreateGame creates a new game.
func (s *Server) CreateGame(context *gin.Context) {
	var req CreateGameRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	homeID := uuid.MustParse(req.HomeTeamID)
	awayID := uuid.MustParse(req.AwayTeamID)

	game, err := s.store.CreateGame(context, db.CreateGameParams{
		HomeTeamID: homeID,
		AwayTeamID: awayID,
		HomeScore:  0,
		AwayScore:  0,
	})
	if err != nil {
		if pgErr, err := err.(*pq.Error); err {
			switch pgErr.Code.Name() {
			case "unique_violation":
				context.JSON(400, helpers.ErrorResponse(pgErr))
				return
			}
		}
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	context.JSON(200, CreateGameResponse{
		ID:         game.ID.String(),
		HomeTeamID: game.HomeTeamID.String(),
		AwayTeamID: game.AwayTeamID.String(),
		HomeScore:  game.HomeScore,
		AwayScore:  game.AwayScore,
	})
}

// ListGamesRequest defines the request body for ListGamesHandler.
type ListGamesRequest struct {
	HomeTeamID string `form:"home_team_id" binding:"omitempty,uuid"`
	AwayTeamID string `form:"away_team_id" binding:"omitempty,uuid"`
	Limit      int32  `form:"limit"`
	Offset     int32  `form:"offset"`
}

// ListGames lists all games.
func (s *Server) ListGames(context *gin.Context) {
	var req ListGamesRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	homeID := uuid.Nil
	awayID := uuid.Nil
	var err error
	if req.HomeTeamID != "" {
		homeID, err = uuid.Parse(req.HomeTeamID)
		if err != nil {
			context.JSON(400, helpers.ErrorResponse(err))
			return
		}
	}
	if req.AwayTeamID != "" {
		awayID, err = uuid.Parse(req.AwayTeamID)
		if err != nil {
			context.JSON(400, helpers.ErrorResponse(err))
			return
		}
	}

	games, err := s.store.ListGames(context, db.ListGamesParams{
		Limit:      req.Limit,
		Offset:     req.Offset,
		HomeTeamID: uuid.NullUUID{UUID: homeID, Valid: homeID != uuid.Nil},
		AwayTeamID: uuid.NullUUID{UUID: awayID, Valid: awayID != uuid.Nil},
	})
	if err != nil {
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	context.JSON(200, games)
}

// GetGameRequest defines the request body for GetGameHandler.
type GetGameRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// GetGameResponse defines the response body for GetGameHandler.
type GetGameResponse struct {
	ID         string `json:"id"`
	HomeTeamID string `json:"home_team_id"`
	AwayTeamID string `json:"away_team_id"`
	HomeScore  int64  `json:"home_score"`
	AwayScore  int64  `json:"away_score"`
}

// GetGame gets a game by ID.
func (s *Server) GetGame(context *gin.Context) {
	var req GetGameRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	id := uuid.MustParse(req.ID)

	game, err := s.store.GetGame(context, id)
	if err != nil {
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	context.JSON(200, GetGameResponse{
		ID:         game.ID.String(),
		HomeTeamID: game.HomeTeamID.String(),
		AwayTeamID: game.AwayTeamID.String(),
		HomeScore:  game.HomeScore,
		AwayScore:  game.AwayScore,
	})
}

// UpdateGameRequest defines the request body for UpdateGameHandler.
type UpdateGameRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// UpdateGameBody defines the request body for UpdateGameHandler.
type UpdateGameBody struct {
	HomeScore int64 `json:"home_score"`
	AwayScore int64 `json:"away_score"`
}

// UpdateGameResponse defines the response body for UpdateGameHandler.
type UpdateGameResponse struct {
	ID         string `json:"id"`
	HomeTeamID string `json:"home_team_id"`
	AwayTeamID string `json:"away_team_id"`
	HomeScore  int64  `json:"home_score"`
	AwayScore  int64  `json:"away_score"`
}

// UpdateGame updates a game by ID.
func (s *Server) UpdateGame(context *gin.Context) {
	var req UpdateGameRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	var body UpdateGameBody
	if err := context.ShouldBindJSON(&body); err != nil {
		context.JSON(400, helpers.ErrorResponse(err))
		return
	}

	id := uuid.MustParse(req.ID)

	game, err := s.store.UpdateGame(context, db.UpdateGameParams{
		ID:        id,
		HomeScore: body.HomeScore,
		AwayScore: body.AwayScore,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(404, helpers.ErrorResponse(err))
			return
		}
		context.JSON(500, helpers.ErrorResponse(err))
		return
	}

	context.JSON(200, UpdateGameResponse{
		ID:         game.ID.String(),
		HomeTeamID: game.HomeTeamID.String(),
		AwayTeamID: game.AwayTeamID.String(),
		HomeScore:  game.HomeScore,
		AwayScore:  game.AwayScore,
	})
}

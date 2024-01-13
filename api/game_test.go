package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/api/middleware"
	mockdb "github.com/kwalter26/scoreit-api-go/db/mock"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"github.com/kwalter26/scoreit-api-go/security"
	"github.com/kwalter26/scoreit-api-go/security/token"
	"github.com/kwalter26/scoreit-api-go/util"
	pg "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServerCreateGame(t *testing.T) {

	homeTeam := createRandomTeam()
	awayTeam := createRandomTeam()
	user, _ := createRandomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMake token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"home_team_id": homeTeam.ID,
				"away_team_id": awayTeam.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateGameParams{
					HomeTeamID: homeTeam.ID,
					AwayTeamID: awayTeam.ID,
					HomeScore:  0,
					AwayScore:  0,
				}
				store.EXPECT().
					CreateGame(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Game{}, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BadRequest (InvalidHomeTeamID)",
			body: gin.H{
				"home_team_id": "fd",
				"away_team_id": awayTeam.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// No expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest (InvalidAwayTeamID)",
			body: gin.H{
				"home_team_id": homeTeam.ID,
				"away_team_id": "fd",
			},
			buildStubs: func(store *mockdb.MockStore) {
				// No expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"home_team_id": homeTeam.ID,
				"away_team_id": awayTeam.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// No expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// No token
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"home_team_id": homeTeam.ID,
				"away_team_id": awayTeam.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateGame(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Game{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "AlreadyExists",
			body: gin.H{
				"home_team_id": homeTeam.ID,
				"away_team_id": awayTeam.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateGame(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Game{}, &pg.Error{Code: "23505"})
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := buildJsonRequest(t, tc.body)

			request, err := http.NewRequest(http.MethodPost, "/api/v1/games", &data)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestServerListGames(t *testing.T) {

	user, _ := createRandomUser(t)

	n := 10
	games := make([]db.Game, n)
	for i := 0; i < n; i++ {
		games[i], _, _ = createRandomGame()
	}

	type Query struct {
		pageID     int
		pageSize   int
		homeTeamID string
		awayTeamID string
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListGamesParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().
					ListGames(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(games, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGames(t, recorder.Body, games)
			},
		},
		{
			name: "BadRequest (InvalidPageSize)",
			query: Query{
				pageID:   1,
				pageSize: -1,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// No expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// No expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// No token
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListGamesParams{
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().
					ListGames(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "OK (WithNextPage)",
			query: Query{
				pageID:   2,
				pageSize: n - 5,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListGamesParams{
					Limit:  int32(n - 5),
					Offset: 5,
				}
				store.EXPECT().
					ListGames(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(games[5:], nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGames(t, recorder.Body, games[5:])
			},
		},
		{
			name: "BadRequest (InvalidHomeTeamID)",
			query: Query{
				pageID:     1,
				pageSize:   n,
				homeTeamID: "asdf",
			},
			buildStubs: func(store *mockdb.MockStore) {
				// No expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest (InvalidAwayTeamID)",
			query: Query{
				pageID:     1,
				pageSize:   n,
				awayTeamID: "asdf",
			},
			buildStubs: func(store *mockdb.MockStore) {
				// No expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "OK (WithHomeTeamID)",
			query: Query{
				pageID:     1,
				pageSize:   n,
				homeTeamID: games[0].HomeTeamID.String(),
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListGamesParams{
					Limit:      int32(n),
					Offset:     0,
					HomeTeamID: uuid.NullUUID{UUID: games[0].HomeTeamID, Valid: true},
				}
				store.EXPECT().
					ListGames(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Game{games[0]}, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGames(t, recorder.Body, games[0:1])
			},
		},
		{
			name: "OK (WithAwayTeamID)",
			query: Query{
				pageID:     1,
				pageSize:   n,
				awayTeamID: games[0].AwayTeamID.String(),
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListGamesParams{
					Limit:      int32(n),
					Offset:     0,
					AwayTeamID: uuid.NullUUID{UUID: games[0].AwayTeamID, Valid: true},
				}
				store.EXPECT().
					ListGames(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Game{games[0]}, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGames(t, recorder.Body, games[0:1])
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/games?page_id=%d&page_size=%d&home_team_id=%s&away_team_id=%s", tc.query.pageID, tc.query.pageSize, tc.query.homeTeamID, tc.query.awayTeamID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestServerGetGame(t *testing.T) {
	user, _ := createRandomUser(t)
	game, _, _ := createRandomGame()

	testCases := []struct {
		name          string
		gameID        string
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			gameID: game.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGame(gomock.Any(), gomock.Eq(game.ID)).
					Times(1).
					Return(game, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGame(t, recorder.Body, game)
			},
		},
		{
			name:   "NotFound",
			gameID: game.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGame(gomock.Any(), gomock.Eq(game.ID)).
					Times(1).
					Return(db.Game{}, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "BadRequest (InvalidID)",
			gameID: "asdf",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGame(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Unauthorized",
			gameID: game.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGame(gomock.Any(), gomock.Eq(game.ID)).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// Don't add authorization header to the request
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			gameID: game.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGame(gomock.Any(), gomock.Eq(game.ID)).
					Times(1).
					Return(db.Game{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/games/%s", tc.gameID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestServerUpdateGame(t *testing.T) {
	user, _ := createRandomUser(t)
	game, _, _ := createRandomGame()
	updateGame := game
	updateGame.HomeScore = util.RandomInt(0, 10)
	updateGame.AwayScore = util.RandomInt(0, 10)

	testCases := []struct {
		name          string
		gameID        string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			gameID: game.ID.String(),
			body: gin.H{
				"home_score": updateGame.HomeScore,
				"away_score": updateGame.AwayScore,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateGameParams{
					ID:        game.ID,
					HomeScore: updateGame.HomeScore,
					AwayScore: updateGame.AwayScore,
				}
				store.EXPECT().
					UpdateGame(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updateGame, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGame(t, recorder.Body, updateGame)
			},
		},
		{
			name:   "NotFound",
			gameID: game.ID.String(),
			body: gin.H{
				"home_score": updateGame.HomeScore,
				"away_score": updateGame.AwayScore,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateGameParams{
					ID:        game.ID,
					HomeScore: updateGame.HomeScore,
					AwayScore: updateGame.AwayScore,
				}
				store.EXPECT().
					UpdateGame(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Game{}, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "BadRequest (InvalidID)",
			gameID: "asdf",
			body:   gin.H{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateGame(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "BadRequest (InvalidHomeScore)",
			gameID: game.ID.String(),
			body: gin.H{
				"home_score": "asdf",
				"away_score": updateGame.AwayScore,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateGame(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "BadRequest (InvalidAwayScore)",
			gameID: game.ID.String(),
			body: gin.H{
				"home_score": updateGame.HomeScore,
				"away_score": "asdf",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateGame(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			gameID: game.ID.String(),
			body: gin.H{
				"home_score": updateGame.HomeScore,
				"away_score": updateGame.AwayScore,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateGame(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Game{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := buildJsonRequest(t, tc.body)

			url := fmt.Sprintf("/api/v1/games/%s", tc.gameID)
			request, err := http.NewRequest(http.MethodPut, url, &data)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func createRandomTeam() db.Team {
	return db.Team{
		ID:   uuid.New(),
		Name: util.RandomName(),
	}
}

func requireBodyMatchGames(t *testing.T, body *bytes.Buffer, games []db.Game) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotGames []db.Game
	err = json.Unmarshal(data, &gotGames)
	require.NoError(t, err)

	require.Equal(t, games, gotGames)
}

func requireBodyMatchGame(t *testing.T, body *bytes.Buffer, game db.Game) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotGame db.Game
	err = json.Unmarshal(data, &gotGame)
	require.NoError(t, err)

	require.Equal(t, game, gotGame)
}

func createRandomGame() (db.Game, db.Team, db.Team) {
	homeTeam := createRandomTeam()
	awayTeam := createRandomTeam()
	game := db.Game{
		HomeTeamID: homeTeam.ID,
		AwayTeamID: awayTeam.ID,
		HomeScore:  util.RandomInt(0, 5),
		AwayScore:  util.RandomInt(0, 5),
	}

	return game, homeTeam, awayTeam
}

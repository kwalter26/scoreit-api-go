package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/api/middleware"
	mockdb "github.com/kwalter26/scoreit-api-go/db/mock"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"github.com/kwalter26/scoreit-api-go/security/token"
	"github.com/kwalter26/udemy-simplebank/util"
	pg "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer_CreateGame(t *testing.T) {

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
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
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
			buildStubs: func(store *mockdb.MockStore) {},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
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
			buildStubs: func(store *mockdb.MockStore) {},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
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
			buildStubs: func(store *mockdb.MockStore) {},
			setupAuth:  func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
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
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
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
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
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

func createRandomTeam() db.Team {
	return db.Team{
		ID:   uuid.New(),
		Name: util.RandomOwner(),
	}
}

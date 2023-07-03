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
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer_ListTeams(t *testing.T) {

	user, _ := createRandomUser(t)

	testCases := []struct {
		name          string
		pageSize      int
		pageNumber    int
		setupAuth     func(t *testing.T, request *http.Request, tokenMake token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			pageSize:   10,
			pageNumber: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.ListTeamsParams{
					Limit:  10,
					Offset: 0,
				}
				store.EXPECT().
					ListTeams(gomock.Any(), args).
					Times(1).
					Return(randomTeams(10), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTeams(t, recorder.Body, 10)
			},
		},
		{
			name:       "Unauthorized",
			pageSize:   10,
			pageNumber: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// do not add authorization header
			},
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:       "BadRequestPageSizeTooSmall",
			pageSize:   -1,
			pageNumber: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "BadRequestPageSizeTooLarge",
			pageSize:   101,
			pageNumber: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)

			},
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "BadRequestPageNumberTooSmall",
			pageSize:   10,
			pageNumber: -1,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "InternalError",
			pageSize:   10,
			pageNumber: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.ListTeamsParams{
					Limit:  10,
					Offset: 0,
				}
				store.EXPECT().
					ListTeams(gomock.Any(), args).
					Times(1).
					Return(nil, sql.ErrConnDone)
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

			url := fmt.Sprintf("/api/v1/teams?page_size=%d&page_id=%d", tc.pageSize, tc.pageNumber)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestCreateTeamAPI(t *testing.T) {

	user, _ := createRandomUser(t)
	team := randomTeam()

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMake token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name": team.Name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateTeam(gomock.Any(), team.Name).
					Times(1).
					Return(team, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTeam(t, recorder.Body, team)
			},
		},
		{
			name: "NoToken",
			body: gin.H{
				"name": team.Name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// no token
			},
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "EmptyName",
			body: gin.H{
				"name": "",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"name": team.Name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateTeam(gomock.Any(), team.Name).
					Times(1).
					Return(db.Team{}, sql.ErrConnDone)
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

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.body)
			require.NoError(t, err)

			url := "/api/v1/teams"
			request, err := http.NewRequest(http.MethodPost, url, &buf)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestServer_ListTeamMembers(t *testing.T) {
	user, _ := createRandomUser(t)
	team := randomTeam()
	teamMember := randomTeamMember(user, team)
	teamMembers := []db.ListTeamMembersRow{teamMember}

	testCases := []struct {
		name          string
		teamID        string
		pageSize      int
		pageID        int
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			teamID:   team.ID.String(),
			pageSize: 5,
			pageID:   0,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListTeamMembersParams{
					TeamID: team.ID,
					Limit:  5,
					Offset: 0,
				}
				store.EXPECT().
					ListTeamMembers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(teamMembers, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUserTeams(t, recorder.Body, 1)
			},
		},
		{
			name:     "NoToken",
			teamID:   team.ID.String(),
			pageSize: 5,
			pageID:   0,
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// no token
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:     "NoTeamID",
			teamID:   "",
			pageSize: 5,
			pageID:   0,
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "InternalError",
			teamID:   team.ID.String(),
			pageSize: 5,
			pageID:   0,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListTeamMembersParams{
					TeamID: team.ID,
					Limit:  5,
					Offset: 0,
				}
				store.EXPECT().
					ListTeamMembers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "PageSizeTooSmall",
			teamID:   team.ID.String(),
			pageSize: 0,
			pageID:   0,
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "PageSizeTooLarge",
			teamID:   team.ID.String(),
			pageSize: 101,
			pageID:   0,
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "PageIDTooSmall",
			teamID:   team.ID.String(),
			pageSize: 101,
			pageID:   -1,
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
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

			url := fmt.Sprintf("/api/v1/teams/%s/members?page_size=%d&page_id=%d", tc.teamID, tc.pageSize, tc.pageID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestServer_AddTeamMember(t *testing.T) {
	user, _ := createRandomUser(t)
	team := randomTeam()
	position := string(util.RandomBaseballPosition())
	testCases := []struct {
		name          string
		teamID        string
		userID        string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			teamID: team.ID.String(),
			userID: user.ID.String(),
			body: gin.H{
				"number":           5,
				"primary_position": position,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.AddTeamMemberParams{
					UserID:          user.ID,
					TeamID:          team.ID,
					Number:          5,
					PrimaryPosition: position,
				}
				store.EXPECT().
					AddTeamMember(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.TeamMember{
						ID:              uuid.UUID{},
						Number:          5,
						PrimaryPosition: position,
						UserID:          user.ID,
						TeamID:          team.ID,
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					}, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "NoAuthorization",
			teamID: team.ID.String(),
			userID: user.ID.String(),
			body: gin.H{
				"number":           5,
				"primary_position": position,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// no authorization header
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoTeamId",
			// teamID: team.ID.String(),
			userID: user.ID.String(),
			body: gin.H{
				"number":           5,
				"primary_position": position,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InvalidTeamId",
			teamID: "invalid uuid",
			userID: user.ID.String(),
			body: gin.H{
				"number":           5,
				"primary_position": position,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InvalidUserId",
			teamID: team.ID.String(),
			userID: "invalid uuid",
			body: gin.H{
				"number":           5,
				"primary_position": position,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "NoNumber",
			teamID: team.ID.String(),
			userID: user.ID.String(),
			body: gin.H{
				"primary_position": position,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "BadNumber",
			teamID: team.ID.String(),
			userID: user.ID.String(),
			body: gin.H{
				"number":           "bad number",
				"primary_position": position,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "NoPrimaryPosition",
			teamID: team.ID.String(),
			userID: user.ID.String(),
			body: gin.H{
				"number": 5,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			teamID: team.ID.String(),
			userID: user.ID.String(),
			body: gin.H{
				"number":           5,
				"primary_position": position,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.AddTeamMemberParams{
					TeamID:          team.ID,
					UserID:          user.ID,
					Number:          5,
					PrimaryPosition: position,
				}
				store.EXPECT().
					AddTeamMember(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.TeamMember{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
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

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/teams/%s/members/%s", tc.teamID, tc.userID)
			request, err := http.NewRequest(http.MethodPut, url, &buf)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestServer_GetTeam(t *testing.T) {

	user, _ := createRandomUser(t)
	team := randomTeam()

	testCases := []struct {
		name          string
		teamID        string
		setupAuth     func(t *testing.T, request *http.Request, tokenMake token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			teamID: team.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTeam(gomock.Any(), gomock.Eq(team.ID)).
					Times(1).
					Return(team, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMake token.Maker) {
				addAuthorization(t, request, tokenMake, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTeam(t, recorder.Body, team)
			},
		},
		{
			name:   "NotFound",
			teamID: team.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTeam(gomock.Any(), gomock.Eq(team.ID)).
					Times(1).
					Return(db.Team{}, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMake token.Maker) {
				addAuthorization(t, request, tokenMake, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			teamID: "invalid_id",
			buildStubs: func(store *mockdb.MockStore) {
				// no expectations
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMake token.Maker) {
				addAuthorization(t, request, tokenMake, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			teamID: team.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTeam(gomock.Any(), gomock.Eq(team.ID)).
					Times(1).
					Return(db.Team{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMake token.Maker) {
				addAuthorization(t, request, tokenMake, middleware.AuthorizationTypeBearer, user.Username, time.Minute)
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

			url := fmt.Sprintf("/api/v1/teams/%s", tc.teamID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchUserTeams(t *testing.T, body *bytes.Buffer, i int) {
	teamMembers := make([]db.TeamMember, 5)
	err := json.NewDecoder(body).Decode(&teamMembers)
	require.NoError(t, err)
	require.NoError(t, err)
	for _, member := range teamMembers {
		require.NotEmpty(t, member)
	}
	require.Equal(t, i, len(teamMembers))
}

func randomTeamMember(user db.User, team db.Team) db.ListTeamMembersRow {
	return db.ListTeamMembersRow{
		ID:              uuid.New(),
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Number:          0,
		PrimaryPosition: string(util.RandomBaseballPosition()),
		TeamName:        team.Name,
	}
}

func requireBodyMatchTeam(t *testing.T, body *bytes.Buffer, team db.Team) {
	var createdUser db.Team
	err := json.NewDecoder(body).Decode(&createdUser)
	require.NoError(t, err)

	require.Equal(t, team.Name, createdUser.Name)
	require.NotEmpty(t, createdUser.ID)
	require.WithinDurationf(t, createdUser.CreatedAt.UTC(), time.Now().UTC(), time.Second, "now should be within 5 seconds of CreatedAt")
}

func requireBodyMatchTeams(t *testing.T, body *bytes.Buffer, i int) {
	teams := make([]db.Team, 5)
	err := json.NewDecoder(body).Decode(&teams)
	require.NoError(t, err)
	require.NoError(t, err)
	for _, account := range teams {
		require.NotEmpty(t, account)
	}
	require.Equal(t, i, len(teams))
}

func randomTeams(i int) []db.Team {
	teams := make([]db.Team, i)
	for i := range teams {
		teams[i] = randomTeam()
	}
	return teams

}

func randomTeam() db.Team {
	return db.Team{
		ID:        uuid.New(),
		Name:      util.RandomString(5),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func addAuthorization(t *testing.T, request *http.Request, tokenMaker token.Maker, authorizationType string, username string, duration time.Duration) {
	createToken, payload, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := authorizationType + " " + createToken
	request.Header.Set(middleware.AuthorizationHeaderKey, authorizationHeader)
}

func createRandomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(6)
	hashedPassword, err := security.HashPassword(password)
	require.NoError(t, err)

	user := db.User{
		ID:             uuid.New(),
		Username:       util.RandomName(),
		FirstName:      util.RandomName(),
		LastName:       util.RandomName(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}
	return user, password
}

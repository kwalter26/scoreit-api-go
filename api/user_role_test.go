package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/api/middleware"
	mockdb "github.com/kwalter26/scoreit-api-go/db/mock"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"github.com/kwalter26/scoreit-api-go/security"
	"github.com/kwalter26/scoreit-api-go/security/token"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestServer_CreateUserRole(t *testing.T) {
	user, _ := createRandomUser(t)
	testCases := []struct {
		name          string
		id            string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMake token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			id:   user.ID.String(),
			body: gin.H{
				"name": "admin",
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateRoleParams{
					Name:   "admin",
					UserID: user.ID,
				}
				store.EXPECT().
					CreateRole(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(db.UserRole{
						ID:     uuid.New(),
						Name:   "admin",
						UserID: user.ID,
					}, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.AdminRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUserRole(t, recorder.Body, "admin")
			},
		},
		{
			name: "BadRequest (Bad UUID)",
			id:   "invalid",
			body: gin.H{
				"name": "admin",
			},
			buildStubs: func(store *mockdb.MockStore) {
				// no stubs needed
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.AdminRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest (Empty Name)",
			id:   user.ID.String(),
			body: gin.H{
				"name": "",
			},
			buildStubs: func(store *mockdb.MockStore) {
				// no stubs needed
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.AdminRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			id:   user.ID.String(),
			body: gin.H{
				"name": "admin",
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateRoleParams{
					Name:   "admin",
					UserID: user.ID,
				}
				store.EXPECT().
					CreateRole(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(db.UserRole{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.AdminRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
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

			reqUrl, err := url.ParseRequestURI("/api/v1/players/" + tc.id + "/roles")
			require.NoError(t, err)

			request, _ := http.NewRequest(http.MethodPut, reqUrl.String(), &buf)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestServer_ListUserRoles(t *testing.T) {
	user, _ := createRandomUser(t)
	testCases := []struct {
		name          string
		pageSize      int32
		pageId        int32
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMake token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			pageSize: 10,
			pageId:   1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListRoles(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.UserRole{
						{
							ID:     uuid.New(),
							Name:   "admin",
							UserID: user.ID,
						},
						{
							ID:     uuid.New(),
							Name:   "user",
							UserID: user.ID,
						},
					}, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.AdminRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUserRoles(t, recorder.Body, []GetUserRolesResponse{{Name: "admin"}, {Name: "user"}})
			},
		},
		{
			name:     "NoRolesFound",
			pageSize: 10,
			pageId:   1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListRoles(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.UserRole{}, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.AdminRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUserRoles(t, recorder.Body, []GetUserRolesResponse{})
			},
		},
		{
			name:     "BadRequest",
			pageSize: 0,
			pageId:   1,
			buildStubs: func(store *mockdb.MockStore) {
				// no stubs needed
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.AdminRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "InternalError",
			pageSize: 10,
			pageId:   1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListRoles(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.UserRole{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.AdminRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			uri, err := url.ParseRequestURI("/api/v1/players/roles")
			require.NoError(t, err)

			query := uri.Query()
			if tc.pageSize >= 0 {
				query.Add("page_size", strconv.Itoa(int(tc.pageSize)))
			}

			if tc.pageId >= 0 {
				// tc.pageId to string

				query.Add("page_id", strconv.Itoa(int(tc.pageId)))
			}
			uri.RawQuery = query.Encode()

			request, _ := http.NewRequest(http.MethodGet, uri.String(), nil)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestServer_GetUserRoles(t *testing.T) {
	user, _ := createRandomUser(t)
	testCases := []struct {
		name          string
		id            string
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMake token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			id:   user.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRoles(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return([]db.UserRole{
						{
							ID:     uuid.New(),
							Name:   "admin",
							UserID: user.ID,
						},
						{
							ID:     uuid.New(),
							Name:   "user",
							UserID: user.ID,
						},
					}, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.AdminRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUserRoles(t, recorder.Body, []GetUserRolesResponse{{Name: "admin"}, {Name: "user"}})
			},
		},
		{
			name: "NoRolesFound",
			id:   user.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRoles(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return([]db.UserRole{}, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.AdminRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUserRoles(t, recorder.Body, []GetUserRolesResponse{})
			},
		},
		{
			name: "BadRequest (Bad UUID)",
			id:   "invalid",
			buildStubs: func(store *mockdb.MockStore) {
				// no stubs needed
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.AdminRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			id:   user.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRoles(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return([]db.UserRole{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, security.AdminRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			reqUrl, err := url.ParseRequestURI("/api/v1/players/" + tc.id + "/roles")
			require.NoError(t, err)

			request, _ := http.NewRequest(http.MethodGet, reqUrl.String(), nil)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchUserRoles(t *testing.T, body *bytes.Buffer, roles []GetUserRolesResponse) {
	var response []GetUserRolesResponse
	err := json.NewDecoder(body).Decode(&response)
	require.NoError(t, err)
	require.Equal(t, roles, response)
}

func requireBodyMatchUserRole(t *testing.T, body *bytes.Buffer, name string) {
	var response CreateUserRolesResponse
	err := json.NewDecoder(body).Decode(&response)
	require.NoError(t, err)
	require.NotEmpty(t, response.Name)
	require.Equal(t, name, response.Name)
}

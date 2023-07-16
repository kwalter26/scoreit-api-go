package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kwalter26/scoreit-api-go/api/middleware"
	mockdb "github.com/kwalter26/scoreit-api-go/db/mock"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"github.com/kwalter26/scoreit-api-go/security/token"
	"github.com/kwalter26/scoreit-api-go/util"
	pg "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer_CreateNewUser(t *testing.T) {

	user, password := createRandomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":   user.Username,
				"password":   password,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:  user.Username,
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Email:     user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParamsMatcher(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":   user.Username,
				"password":   password,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username":   user.Username,
				"password":   password,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pg.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":   "invalid-user#",
				"password":   password,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username":   user.Username,
				"password":   password,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      "invalid-email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username":   user.Username,
				"password":   util.RandomString(5),
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooLongPassword",
			body: gin.H{
				"username":   user.Username,
				"password":   util.RandomString(101),
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
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

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/users")
			request, err := http.NewRequest(http.MethodPost, url, &buf)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func TestServer_ListUsers(t *testing.T) {
	// create 3 users and store in users
	users := make([]db.User, 3)
	for i := range users {
		users[i], _ = createRandomUser(t)
	}

	listUserRows := make([]db.ListUsersRow, len(users))
	for i, user := range users {
		listUserRows[i] = db.ListUsersRow{
			ID:        user.ID,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
	}

	usersResponse := make([]UserResponse, len(users))
	for i, user := range users {
		usersResponse[i] = UserResponse{
			ID:        user.ID.String(),
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
	}

	testCases := []struct {
		name          string
		pageSize      int
		pageID        int
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMake token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			pageSize: 5,
			pageID:   1,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListUsersParams{
					Limit:  5,
					Offset: 0,
				}
				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(listUserRows, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, users[0].Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUsers(t, recorder.Body, usersResponse)
			},
		},
		{
			name:     "InternalError",
			pageSize: 5,
			pageID:   1,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListUsersParams{
					Limit:  5,
					Offset: 0,
				}
				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, users[0].Username, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "Unauthenticated",
			pageSize: 5,
			pageID:   1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// do nothing
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:     "InvalidPageSize",
			pageSize: 0,
			pageID:   1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middleware.AuthorizationTypeBearer, users[0].Username, time.Minute)
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

			url := fmt.Sprintf("/api/v1/players")
			if tc.pageSize > 0 || tc.pageID > 0 {
				url = fmt.Sprintf("/api/v1/players?page_size=%v&page_id=%v", tc.pageSize, tc.pageID)
			}
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestServer_LoginUser(t *testing.T) {
	user, password := createRandomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				loginResponseValid(t, recorder.Body, user)
			},
		},
		{
			name: "BadRequest (InvalidUsername)",
			body: gin.H{
				"username": "fd",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest (INVALIDPassword)",
			body: gin.H{
				"username": user.Username,
				"password": "as",
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadPassword",
			body: gin.H{
				"username": user.Username,
				"password": "badpassword",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			url := "/api/v1/login"
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, &buf)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	var createdUser db.User
	err := json.NewDecoder(body).Decode(&createdUser)
	require.NoError(t, err)

	require.NotEmpty(t, createdUser)
	require.Equal(t, user.Username, createdUser.Username)
	require.Equal(t, user.FirstName, createdUser.FirstName)
	require.Equal(t, user.LastName, createdUser.LastName)
	require.Equal(t, user.Email, createdUser.Email)
}

func requireBodyMatchUsers(t *testing.T, body *bytes.Buffer, users []UserResponse) {
	var createdUsers []db.User
	err := json.NewDecoder(body).Decode(&createdUsers)
	require.NoError(t, err)

	require.NotEmpty(t, createdUsers)
	require.Equal(t, len(users), len(createdUsers))
}

func loginResponseValid(t *testing.T, body *bytes.Buffer, user db.User) {
	var loginResponse LoginUserResponse
	err := json.NewDecoder(body).Decode(&loginResponse)
	require.NoError(t, err)

	require.NotEmpty(t, loginResponse)
	require.NotEmpty(t, loginResponse.AccessToken)
	require.Equal(t, user.Username, loginResponse.User.Username)
	require.Equal(t, user.FirstName, loginResponse.User.FirstName)
	require.Equal(t, user.LastName, loginResponse.User.LastName)
	require.Equal(t, user.Email, loginResponse.User.Email)
}

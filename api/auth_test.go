package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/db/mock"
	"github.com/kwalter26/scoreit-api-go/db/sqlc"
	"github.com/kwalter26/scoreit-api-go/security"
	"github.com/kwalter26/scoreit-api-go/security/token"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServerLoginUser(t *testing.T) {
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
					GetRoles(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return([]db.UserRole{{Name: "user"}}, nil)
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
			buildStubs: func(store *mockdb.MockStore) {
				// do nothing
			},
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
			buildStubs: func(store *mockdb.MockStore) {
				// do nothing
			},
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

			url := "/api/v1/auth/login"
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

func TestServerLogoutUser(t *testing.T) {
	user, _ := createRandomUser(t)

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore, tokenMaker token.Maker) gin.H
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore, tokenMaker token.Maker) gin.H {
				createToken, payload, err := tokenMaker.CreateToken(user.ID, security.UserRoles, time.Minute)
				require.NoError(t, err)

				arg := db.UpdateSessionParams{
					ID:        payload.ID,
					IsBlocked: true,
				}
				store.EXPECT().
					UpdateSession(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Session{}, nil)
				return gin.H{
					"refresh_token": createToken,
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			buildStubs: func(store *mockdb.MockStore, tokenMaker token.Maker) gin.H {
				createToken, payload, err := tokenMaker.CreateToken(user.ID, security.UserRoles, time.Minute)
				require.NoError(t, err)

				arg := db.UpdateSessionParams{
					ID:        payload.ID,
					IsBlocked: true,
				}
				store.EXPECT().
					UpdateSession(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Session{}, sql.ErrConnDone)
				return gin.H{
					"refresh_token": createToken,
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest (INVALIDToken)",
			buildStubs: func(store *mockdb.MockStore, tokenMaker token.Maker) gin.H {
				return gin.H{
					"refresh_token": "",
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Unauthorized (INVALIDToken)",
			buildStubs: func(store *mockdb.MockStore, tokenMaker token.Maker) gin.H {
				return gin.H{
					"refresh_token": "INVALIDToken",
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "SessionNotFound",
			buildStubs: func(store *mockdb.MockStore, tokenMaker token.Maker) gin.H {
				createToken, payload, err := tokenMaker.CreateToken(user.ID, security.UserRoles, time.Minute)
				require.NoError(t, err)

				arg := db.UpdateSessionParams{
					ID:        payload.ID,
					IsBlocked: true,
				}
				store.EXPECT().
					UpdateSession(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Session{}, sql.ErrNoRows)

				return gin.H{
					"refresh_token": createToken,
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			body := tc.buildStubs(store, server.tokenMaker)

			data, err := buildJsonRequest(t, body)

			url := "/api/v1/auth/logout"
			request, err := http.NewRequest(http.MethodPost, url, &data)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestServerRefreshToken(t *testing.T) {

	user, _ := createRandomUser(t)

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore, tokenMaker token.Maker) gin.H
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore, tokenMaker token.Maker) gin.H {
				createToken, payload, err := tokenMaker.CreateToken(user.ID, security.UserRoles, time.Minute)
				require.NoError(t, err)

				session := db.Session{
					ID:           payload.ID,
					UserID:       user.ID,
					RefreshToken: createToken,
					UserAgent:    "",
					ClientIp:     "",
					IsBlocked:    false,
					ExpiresAt: sql.NullTime{
						Time:  time.Now().Add(time.Minute),
						Valid: true,
					},
					CreatedAt: time.Now().Add(time.Minute),
				}

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(payload.ID)).
					Times(1).
					Return(session, nil)

				return gin.H{
					"refresh_token": createToken,
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				renewResponseValid(t, recorder.Body)
			},
		},
		{
			name: "InternalServerError (GetSession)",
			buildStubs: func(store *mockdb.MockStore, tokenMaker token.Maker) gin.H {
				createToken, payload, err := tokenMaker.CreateToken(user.ID, security.UserRoles, time.Minute)
				require.NoError(t, err)

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(payload.ID)).
					Times(1).
					Return(db.Session{}, sql.ErrConnDone)
				return gin.H{
					"refresh_token": createToken,
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest (INVALIDToken)",
			buildStubs: func(store *mockdb.MockStore, tokenMaker token.Maker) gin.H {
				return gin.H{
					"refresh_token": "",
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound (SessionNotFound)",
			buildStubs: func(store *mockdb.MockStore, tokenMaker token.Maker) gin.H {
				createToken, payload, err := tokenMaker.CreateToken(user.ID, security.UserRoles, time.Minute)
				require.NoError(t, err)

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(payload.ID)).
					Times(1).
					Return(db.Session{}, sql.ErrNoRows)

				return gin.H{
					"refresh_token": createToken,
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Unauthorized (SessionIsBlocked)",
			buildStubs: func(store *mockdb.MockStore, tokenMaker token.Maker) gin.H {
				createToken, payload, err := tokenMaker.CreateToken(user.ID, security.UserRoles, time.Minute)
				require.NoError(t, err)

				session := db.Session{
					ID:           payload.ID,
					UserID:       user.ID,
					RefreshToken: createToken,
					UserAgent:    "",
					ClientIp:     "",
					IsBlocked:    true,
					ExpiresAt: sql.NullTime{
						Time:  time.Now().Add(time.Minute),
						Valid: true,
					},
					CreatedAt: time.Now().Add(time.Minute),
				}

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(payload.ID)).
					Times(1).
					Return(session, nil)

				return gin.H{
					"refresh_token": createToken,
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Unauthorized (SessionExpired)",
			buildStubs: func(store *mockdb.MockStore, tokenMaker token.Maker) gin.H {
				createToken, payload, err := tokenMaker.CreateToken(user.ID, security.UserRoles, time.Minute)
				require.NoError(t, err)

				session := db.Session{
					ID:           payload.ID,
					UserID:       user.ID,
					RefreshToken: createToken,
					UserAgent:    "",
					ClientIp:     "",
					IsBlocked:    false,
					ExpiresAt: sql.NullTime{
						Time:  time.Now().Add(-time.Minute),
						Valid: true,
					},
					CreatedAt: time.Now().Add(-time.Minute),
				}

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(payload.ID)).
					Times(1).
					Return(session, nil)

				return gin.H{
					"refresh_token": createToken,
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Unauthorized (RefreshTokenMismatch)",
			buildStubs: func(store *mockdb.MockStore, tokenMaker token.Maker) gin.H {
				createToken, payload, err := tokenMaker.CreateToken(user.ID, security.UserRoles, time.Minute)
				require.NoError(t, err)
				anotherToken, _, err := tokenMaker.CreateToken(user.ID, security.UserRoles, time.Minute)
				require.NoError(t, err)

				session := db.Session{
					ID:           payload.ID,
					UserID:       user.ID,
					RefreshToken: anotherToken,
					UserAgent:    "",
					ClientIp:     "",
					IsBlocked:    false,
					ExpiresAt: sql.NullTime{
						Time:  time.Now().Add(-time.Minute),
						Valid: true,
					},
					CreatedAt: time.Now().Add(-time.Minute),
				}

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(session, nil)

				return gin.H{
					"refresh_token": createToken,
				}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Unauthorized (UserIDMismatch)",
			buildStubs: func(store *mockdb.MockStore, tokenMaker token.Maker) gin.H {
				createToken, payload, err := tokenMaker.CreateToken(user.ID, security.UserRoles, time.Minute)
				require.NoError(t, err)

				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(payload.ID)).
					Times(1).
					Return(db.Session{
						ID:           payload.ID,
						UserID:       uuid.New(),
						RefreshToken: createToken,
						UserAgent:    "",
						ClientIp:     "",
						IsBlocked:    false,
						ExpiresAt: sql.NullTime{
							Time:  time.Now().Add(time.Minute),
							Valid: true,
						},
						CreatedAt: time.Now().Add(time.Minute),
					}, nil)

				return gin.H{
					"refresh_token": createToken,
				}
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
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			body := tc.buildStubs(store, server.tokenMaker)

			data, err := buildJsonRequest(t, body)

			url := "/api/v1/auth/renew"
			request, err := http.NewRequest(http.MethodPost, url, &data)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func renewResponseValid(t *testing.T, body *bytes.Buffer) {
	var refreshToken RefreshTokenResponse
	err := json.NewDecoder(body).Decode(&refreshToken)
	require.NoError(t, err)

	require.NotEmpty(t, refreshToken)
	require.NotEmpty(t, refreshToken.AccessToken)
}

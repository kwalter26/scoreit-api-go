package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/kwalter26/scoreit-api-go/db/mock"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"github.com/kwalter26/scoreit-api-go/util"
	pg "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
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

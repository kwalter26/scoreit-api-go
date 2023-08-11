package api

import (
	"bytes"
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
	"testing"
	"time"
)

func TestServer_CreateUserRole(t *testing.T) {
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
				addAuthorization(t, request, tokenMaker, security.UserRoles, middleware.AuthorizationTypeBearer, user.ID, time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUserRole(t, recorder.Body, "admin")
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

			url := "/api/v1/players/" + user.ID.String() + "/roles"
			request, _ := http.NewRequest(http.MethodPut, url, &buf)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchUserRole(t *testing.T, body *bytes.Buffer, name string) {
	var response CreateUserRolesResponse
	err := json.NewDecoder(body).Decode(&response)
	require.NoError(t, err)
	require.NotEmpty(t, response.Name)
	require.Equal(t, name, response.Name)
}

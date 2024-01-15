package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kwalter26/scoreit-api-go/api/middleware"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"github.com/kwalter26/scoreit-api-go/security"
	"github.com/kwalter26/scoreit-api-go/security/token"
	"github.com/kwalter26/scoreit-api-go/util"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	TokenSecret   = "secret"
	TokenIssuer   = "https://localhost/"
	TokenAudience = "me"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
		CasbinPolicyPath:    "../security/authz_policy.csv",
		CasbinModelPath:     "../security/authz_model.conf",
	}

	jwtValidator, _ := validator.New(
		func(ctx context.Context) (interface{}, error) {
			return []byte(TokenSecret), nil
		},
		validator.HS256,
		TokenIssuer,
		[]string{TokenAudience},
		validator.WithCustomClaims(func() validator.CustomClaims {
			return &middleware.CustomClaims{}
		}),
		validator.WithAllowedClockSkew(time.Minute),
	)

	server, err := NewServer(config, store, jwtValidator)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := security.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParamsMatcher(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg: arg, password: password}
}

type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
}

func (e eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	err := security.CheckPassword(e.password, arg.CreateUserParams.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.CreateUserParams.HashedPassword = arg.CreateUserParams.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserTxParamsMatcher(arg db.CreateUserTxParams, password string) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg: arg, password: password}
}

func addAuthorization(t *testing.T, request *http.Request, tokenMaker token.Maker, roles []security.Role, authorizationType string, u uuid.UUID, duration time.Duration) {

	type CustomClaims struct {
		middleware.CustomClaims
		jwt.StandardClaims
	}

	claims := CustomClaims{
		CustomClaims: middleware.CustomClaims{Roles: roles, Scope: "openid"},
		StandardClaims: jwt.StandardClaims{
			Issuer:    TokenIssuer,
			Audience:  TokenAudience,
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	jwtTokenGenerator := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := jwtTokenGenerator.SignedString([]byte(TokenSecret))
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)

	authorizationHeader := authorizationType + " " + accessToken
	request.Header.Set(middleware.AuthorizationHeaderKey, authorizationHeader)
}

func buildJsonRequest(t *testing.T, body gin.H) (bytes.Buffer, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	require.NoError(t, err)
	return buf, err
}

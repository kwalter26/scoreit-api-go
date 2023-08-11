package api

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/kwalter26/scoreit-api-go/api/middleware"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"github.com/kwalter26/scoreit-api-go/security/token"
	"github.com/kwalter26/scoreit-api-go/util"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
	"time"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	//app    *newrelic.Application
}

func (s *Server) setupRouter() {
	_, err := casbin.NewEnforcer(s.config.CasbinModelPath, s.config.CasbinPolicyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create casbin enforcer")
	}
	log.Info().Msg("created casbin enforcer")

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())

	router.POST("/api/v1/users", s.CreateNewUser)
	router.POST("/api/v1/auth/login", s.LoginUser)
	router.POST("/api/v1/auth/renew", s.RefreshToken)
	router.POST("/api/v1/auth/logout", s.LogoutUser)

	authRoutes := router.Group("/api/")

	issuerURL, err := url.Parse("https://" + s.config.Auth0Domain + "/")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse the issuer url")
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, _ := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{s.config.Auth0Audience},
	)
	jwtMiddleware := jwtmiddleware.New(jwtValidator.ValidateToken)
	authRoutes.Use(adapter.Wrap(jwtMiddleware.CheckJWT))
	//authRoutes.Use(middleware.AuthMiddleware(s.tokenMaker))
	//authRoutes.Use(middleware.NewAuthorizeMiddleware(e))

	authRoutes.GET("/v1/teams", s.ListTeams)
	authRoutes.POST("/v1/teams", s.CreateTeam)
	authRoutes.PUT("/v1/teams/:id/members/:user_id", s.AddTeamMember)
	authRoutes.GET("/v1/teams/:id/members", s.ListTeamMembers)
	authRoutes.GET("/v1/teams/:id", s.GetTeam)

	authRoutes.GET("/v1/players", s.ListUsers)
	authRoutes.GET("/v1/players/roles", s.ListUserRoles)
	authRoutes.GET("/v1/players/:id", s.GetUser)
	authRoutes.GET("/v1/players/:id/roles", s.GetUserRoles)
	authRoutes.PUT("/v1/players/:id/roles", s.CreateUserRole)

	authRoutes.POST("/v1/games", s.CreateGame)
	authRoutes.GET("/v1/games", s.ListGames)
	authRoutes.GET("/v1/games/:id", s.GetGame)
	authRoutes.PUT("/v1/games/:id", s.UpdateGame)

	s.router = router
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	// print current working directory
	dir, err := os.Getwd()
	log.Info().Str("dir", dir).Msg("current working directory")
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{store: store, config: config, tokenMaker: tokenMaker}

	server.setupRouter()

	return server, nil
}

func (s *Server) Start(address string) error {
	log.Info().Str("address", address).Msg("starting server")
	return s.router.Run(address)
}

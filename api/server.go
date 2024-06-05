package api

import (
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/kwalter26/scoreit-api-go/api/middleware"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"github.com/kwalter26/scoreit-api-go/security"
	"github.com/kwalter26/scoreit-api-go/security/token"
	"github.com/kwalter26/scoreit-api-go/util"
	"github.com/rs/zerolog/log"
	"os"
)

type Server struct {
	config         util.Config
	store          db.Store
	tokenMaker     token.Maker
	tokenValidator *validator.Validator
	router         *gin.Engine
	//app    *newrelic.Application
}

func (s *Server) setupRouter() {
	enforcer, err := security.NewEnforcer(s.config, security.SecurityResources())
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
	authRoutes.Use(middleware.CheckJWT(s.tokenValidator))
	authRoutes.Use(middleware.NewAuthorizeMiddleware(enforcer))

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

func NewServer(config util.Config, store db.Store, tokenValidator *validator.Validator) (*Server, error) {
	// print current working directory
	dir, err := os.Getwd()
	log.Info().Str("dir", dir).Msg("current working directory")
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{store: store, config: config, tokenMaker: tokenMaker, tokenValidator: tokenValidator}

	server.setupRouter()

	return server, nil
}

func (s *Server) Start(address string) error {
	log.Info().Str("address", address).Msg("starting server")
	return s.router.Run(address)
}

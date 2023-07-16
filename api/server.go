package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kwalter26/scoreit-api-go/api/middleware"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"github.com/kwalter26/scoreit-api-go/security/token"
	"github.com/kwalter26/scoreit-api-go/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	//app    *newrelic.Application
}

func (s *Server) setupRouter() {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())

	router.POST("/api/v1/users", s.CreateNewUser)
	router.POST("/api/v1/login", s.LoginUser)

	authRoutes := router.Group("/api/").Use(middleware.AuthMiddleware(s.tokenMaker))

	authRoutes.GET("/v1/teams", s.ListTeams)
	authRoutes.POST("/v1/teams", s.CreateTeam)
	authRoutes.PUT("/v1/teams/:id/members/:user_id", s.AddTeamMember)
	authRoutes.GET("/v1/teams/:id/members", s.ListTeamMembers)
	authRoutes.GET("/v1/teams/:id", s.GetTeam)

	authRoutes.GET("/v1/players", s.ListUsers)

	s.router = router
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{store: store, config: config, tokenMaker: tokenMaker}

	server.setupRouter()

	return server, nil
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

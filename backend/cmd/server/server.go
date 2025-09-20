package server

import (
	"github.com/PingTower/internal/handlers"
	"github.com/PingTower/internal/repository/postgres"
	"github.com/PingTower/internal/repository/redis"
	"github.com/PingTower/internal/routes"
	"github.com/PingTower/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Server struct {
	Addr string
	DB   *sqlx.DB
}

func NewServer(addr string, db *sqlx.DB) *Server {
	return &Server{
		Addr: addr,
		DB:   db,
	}
}

func (s *Server) Run(e *echo.Echo, sugar *zap.SugaredLogger, db *sqlx.DB) error {
	api := e.Group("/api/v1")

	serviceRepoPostgres := postgres.NewServiceRepoPostgres(db)
	serviceRepoRedis := redis.NewServiceRepoRedis()
	serviceManager := service.NewServiceManager(serviceRepoPostgres, serviceRepoRedis)
	serviceHandler := handlers.NewServiceHandler(serviceManager)

	routes.RegisterRoutes(api, serviceHandler)

	// starting server
	e.Logger.Fatal(e.Start(s.Addr))

	return nil
}

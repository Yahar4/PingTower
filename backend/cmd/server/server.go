package server

import (
	"context"
	"github.com/PingTower/internal/handlers"
	"github.com/PingTower/internal/repository/postgres"
	"github.com/PingTower/internal/repository/redis"
	"github.com/PingTower/internal/routes"
	"github.com/PingTower/internal/service"
	"github.com/PingTower/pkg/cache"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	red "github.com/redis/go-redis/v9"
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

func (s *Server) Run(e *echo.Echo, sugar *zap.SugaredLogger, db *sqlx.DB, rdb *red.Client) error {
	api := e.Group("/api/v1")

	serviceRepoPostgres := postgres.NewServiceRepoPostgres(db)
	serviceRepoRedis := redis.NewServiceRepoRedis()

	cacheInitializer := cache.NewCacheInitializer(serviceRepoPostgres, rdb)
	if err := cacheInitializer.WarmUpCache(context.Background()); err != nil {
		sugar.Fatalw("failed to warm up cache", "error", err)
	}
	sugar.Info("cache successfully warmed up")

	serviceManager := service.NewServiceManager(serviceRepoPostgres, serviceRepoRedis)
	serviceHandler := handlers.NewServiceHandler(serviceManager)

	routes.RegisterRoutes(api, serviceHandler)

	// starting server
	e.Logger.Fatal(e.Start(s.Addr))

	return nil
}

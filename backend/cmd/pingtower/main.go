package main

import (
	"github.com/PingTower/cmd/server"
	"github.com/PingTower/pkg/database"
	"github.com/PingTower/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	l := logger.CreateLogger()

	sugar := l.Sugar()
	defer func(sugar *zap.SugaredLogger) {
		err := sugar.Sync()
		if err != nil {
			panic(err)
		}
	}(sugar)

	db, err := database.ConnectPostgresDB(sugar)
	if err != nil {
		sugar.Fatalw("Failed to initialize database",
			"error", err,
		)
	}
	defer db.Close()

	s := server.NewServer(":8080", db)
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderContentType},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.Run(e, sugar, db); err != nil {
			log.Fatal(err)
		}
	}()

	sugar.Infof("Server is succesfully running on port  %s", s.Addr)

	<-quit
	sugar.Info("Server is shutting down...")
}

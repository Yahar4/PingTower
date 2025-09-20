package database

import (
	"fmt"

	"github.com/PingTower/pkg/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// ConnectDB connects to a DataBase with specified
func ConnectPostgresDB(logger *zap.SugaredLogger, config *config.Config) (*sqlx.DB, error) {
	dsn := config.GetDSN()

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	logger.Info("Successfully connected to database")
	return db, nil
}

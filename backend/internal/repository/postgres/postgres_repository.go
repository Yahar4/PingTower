package postgres

import (
	"context"
	"fmt"
	"github.com/PingTower/internal/entities"
	"github.com/jmoiron/sqlx"
)

type ServiceRepositoryPostgres struct {
	db *sqlx.DB
}

func NewServiceRepoPostgres(database *sqlx.DB) *ServiceRepositoryPostgres {
	return &ServiceRepositoryPostgres{
		db: database,
	}
}

func (r *ServiceRepositoryPostgres) AddService(ctx context.Context, service entities.Service) error {
	query := `
		INSERT INTO services (id, service_name, url, interval_seconds, active)
		VALUES ($1, $2, $3, $4, $5)
    `

	_, err := r.db.ExecContext(
		ctx,
		query,
		service.ID,
		service.ServiceName,
		service.URL,
		int(service.Interval.Seconds()),
		service.Active,
	)

	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		return fmt.Errorf("postgres: %w", err)
	}

	return err
}

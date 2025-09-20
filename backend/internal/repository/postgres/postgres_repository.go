package postgres

import (
	"context"
	"fmt"
	"github.com/PingTower/internal/entities"
	"github.com/jmoiron/sqlx"
	"time"
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

func (r *ServiceRepositoryPostgres) GetAllServices(ctx context.Context) ([]entities.Service, error) {
	query := `
		SELECT id, service_name, url, interval_seconds, active
		FROM services
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("postgres: %w", err)
	}
	defer rows.Close()

	var services []entities.Service

	for rows.Next() {
		var svc entities.Service
		var intervalSeconds int64

		err := rows.Scan(
			&svc.ID,
			&svc.ServiceName,
			&svc.URL,
			&intervalSeconds,
			&svc.Active,
		)
		if err != nil {
			return nil, fmt.Errorf("postgres scan: %w", err)
		}

		svc.Interval = time.Duration(intervalSeconds) * time.Second
		services = append(services, svc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres rows: %w", err)
	}

	return services, nil
}

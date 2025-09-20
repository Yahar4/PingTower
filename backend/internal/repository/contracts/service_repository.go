package contracts

import (
	"context"
	"github.com/google/uuid"

	"github.com/PingTower/internal/entities"
)

type PostgresRepository interface {
	AddService(ctx context.Context, service entities.Service) error
	GetAllServices(ctx context.Context) ([]entities.Service, error)
	UpdateService(ctx context.Context, service entities.Service) error
	DeleteService(ctx context.Context, serviceID uuid.UUID) error
}

type RedisRepository interface {
	SetService(ctx context.Context, service entities.Service) error
	GetAllServices(ctx context.Context) ([]entities.Service, error)
	UpdateService(ctx context.Context, service entities.Service) error
	DeleteService(ctx context.Context, id uuid.UUID) error
}

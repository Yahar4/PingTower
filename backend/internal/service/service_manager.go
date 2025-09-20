package service

import (
	"context"
	"github.com/PingTower/internal/entities"
	"github.com/PingTower/internal/repository/contracts"
	"github.com/google/uuid"
	"time"
)

type ServiceManager struct {
	pg    contracts.PostgresRepository
	redis contracts.RedisRepository
}

func NewServiceManager(pg contracts.PostgresRepository, redis contracts.RedisRepository) *ServiceManager {
	return &ServiceManager{pg: pg, redis: redis}
}

func (s *ServiceManager) CreateService(ctx context.Context, req *entities.CreateServiceRequest) error {
	serv := entities.Service{
		ID:          uuid.New(),
		ServiceName: req.ServiceName,
		URL:         req.URL,
		Interval:    time.Duration(req.Interval) * time.Second,
		Active:      req.Active,
	}

	if err := s.pg.AddService(ctx, serv); err != nil {
		return err
	}

	if err := s.redis.SetService(ctx, serv); err != nil {
		return err
	}

	return nil
}

func (s *ServiceManager) GetAllServices(ctx context.Context) ([]entities.Service, error) {
	services, err := s.redis.GetAllServices(ctx)
	if err == nil && len(services) > 0 {
		return services, nil
	}

	services, err = s.pg.GetAllServices(ctx)
	if err != nil {
		return nil, err
	}

	for _, serv := range services {
		_ = s.redis.SetService(ctx, serv)
	}

	return services, nil
}

func (s *ServiceManager) UpdateService(ctx context.Context, req *entities.UpdateServiceRequest) error {
	serv := entities.Service{
		ID:          req.ID,
		ServiceName: req.ServiceName,
		URL:         req.URL,
		Interval:    time.Duration(req.Interval) * time.Second,
		Active:      req.Active,
	}

	if err := s.pg.UpdateService(ctx, serv); err != nil {
		return err
	}
	if err := s.redis.SetService(ctx, serv); err != nil {
		return err
	}

	return nil
}

func (s *ServiceManager) DeleteService(ctx context.Context, id uuid.UUID) error {
	if err := s.pg.DeleteService(ctx, id); err != nil {
		return err
	}
	if err := s.redis.DeleteService(ctx, id); err != nil {
		return err
	}
	return nil
}

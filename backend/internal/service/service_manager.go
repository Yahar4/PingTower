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
		ID:       uuid.New(),
		Name:     req.Name,
		URL:      req.URL,
		Interval: time.Duration(req.Interval) * time.Second,
		Active:   req.Active,
	}

	if err := s.pg.AddService(ctx, serv); err != nil {
		return err
	}

	if err := s.redis.SetService(ctx, serv); err != nil {
		return err
	}

	return nil
}

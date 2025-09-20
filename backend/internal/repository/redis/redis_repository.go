package redis

import (
	"context"
	"github.com/PingTower/internal/entities"
	"github.com/PingTower/internal/repository/contracts"
)

type RedisRepo struct{}

func NewServiceRepoRedis() contracts.RedisRepository {
	return &RedisRepo{}
}

func (r *RedisRepo) SetService(ctx context.Context, service entities.Service) error {
	return nil
}

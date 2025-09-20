package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PingTower/internal/repository/postgres"
	"github.com/redis/go-redis/v9"
)

type CacheInitializer struct {
	postgresRepo *postgres.ServiceRepositoryPostgres
	redisClient  *redis.Client
}

func NewCacheInitializer(pg *postgres.ServiceRepositoryPostgres, rdb *redis.Client) *CacheInitializer {
	return &CacheInitializer{
		postgresRepo: pg,
		redisClient:  rdb,
	}
}

func (c *CacheInitializer) WarmUpCache(ctx context.Context) error {
	services, err := c.postgresRepo.GetAllServices(ctx)
	if err != nil {
		return fmt.Errorf("failed to load services from postgres: %w", err)
	}

	for _, s := range services {
		data, err := json.Marshal(s)
		if err != nil {
			return fmt.Errorf("failed to marshal service %s: %w", s.ID, err)
		}

		key := fmt.Sprintf("service:%s", s.ID.String())

		if err := c.redisClient.Set(ctx, key, data, 0).Err(); err != nil {
			return fmt.Errorf("failed to set redis key %s: %w", key, err)
		}
	}

	return nil
}

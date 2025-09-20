package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PingTower/internal/entities"
	"github.com/PingTower/internal/repository/contracts"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewServiceRepoRedis(client *redis.Client) contracts.RedisRepository {
	return &RedisRepo{client: client}
}

func (r *RedisRepo) SetService(ctx context.Context, service entities.Service) error {
	data, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("redis: failed to marshal service %s: %w", service.ID, err)
	}

	key := fmt.Sprintf("service:%s", service.ID.String())

	if err := r.client.Set(ctx, key, data, 0).Err(); err != nil {
		return fmt.Errorf("redis: failed to set key %s: %w", key, err)
	}

	return nil
}

func (r *RedisRepo) GetAllServices(ctx context.Context) ([]entities.Service, error) {
	keys, err := r.client.Keys(ctx, "service:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch keys: %w", err)
	}

	services := make([]entities.Service, 0, len(keys))
	for _, key := range keys {
		val, err := r.client.Get(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get value for key %s: %w", key, err)
		}

		var s entities.Service
		if err := json.Unmarshal([]byte(val), &s); err != nil {
			return nil, fmt.Errorf("failed to unmarshal value for key %s: %w", key, err)
		}

		services = append(services, s)
	}

	return services, nil
}

func (r *RedisRepo) UpdateService(ctx context.Context, service entities.Service) error {
	data, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service: %w", err)
	}

	key := r.key(service.ID)
	if err := r.client.Set(ctx, key, data, 0).Err(); err != nil {
		return fmt.Errorf("failed to update service in redis: %w", err)
	}

	return nil
}

func (r *RedisRepo) DeleteService(ctx context.Context, id uuid.UUID) error {
	key := r.key(id)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete service from redis: %w", err)
	}
	return nil
}

func (r *RedisRepo) key(id uuid.UUID) string {
	return fmt.Sprintf("service:%s", id.String())
}

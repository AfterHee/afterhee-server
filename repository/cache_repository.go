package repository

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	GetValue(ctx context.Context, key string) (*string, error)
	SetValue(ctx context.Context, key string, value string, expiration time.Duration) error
}

type cacheRepository struct {
	rdb *redis.Client
}

func NewCacheRepository(rdb *redis.Client) CacheRepository {
	return &cacheRepository{rdb: rdb}
}

func (s *cacheRepository) GetValue(ctx context.Context, key string) (*string, error) {
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		log.Printf("falied to fetch data from redis. key=%s error=%v\n", key, err)
		return nil, err
	}

	return &val, nil
}

func (s *cacheRepository) SetValue(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := s.rdb.SetNX(ctx, key, value, expiration).Err()
	if err != nil {
		log.Printf("falied to put data to redis. key=%s error=%v\n", key, err)
		return err
	}
	return nil
}

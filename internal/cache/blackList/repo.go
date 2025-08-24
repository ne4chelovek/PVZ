package blackList

import (
	"PVZ/internal/cache"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type clientRedis struct {
	redisClient *redis.Client
}

func NewBlackList(redisClient *redis.Client) cache.BlackList {
	return &clientRedis{redisClient: redisClient}
}

func (b *clientRedis) AddToBlacklist(token string, expiration time.Duration) error {
	return b.redisClient.Set(context.Background(), token, "true", expiration).Err()
}

func (b *clientRedis) IsTokenBlacklisted(token string) (bool, error) {
	val, err := b.redisClient.Get(context.Background(), token).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return val == "true", nil
}

func (b *clientRedis) Get(ctx context.Context, token string) (string, error) {
	return b.redisClient.Get(ctx, "blacklist:"+token).Result()
}

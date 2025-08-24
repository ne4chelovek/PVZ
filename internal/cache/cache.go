package cache

import (
	"context"
	"time"
)

type BlackList interface {
	Get(ctx context.Context, token string) (string, error)
	AddToBlacklist(token string, expiration time.Duration) error
	IsTokenBlacklisted(token string) (bool, error)
}

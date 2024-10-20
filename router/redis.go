package router

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

// RedisClient defines the methods we will use from Redis
type RedisClient interface {
	Incr(key string) *redis.IntCmd
	Expire(key string, expiration time.Duration) *redis.BoolCmd
}

type RateLimiter interface {
	AllowRequest(userID string) (bool, error)
}

// Developing Rate Limiter
type RedisRateLimiter struct {
	client    RedisClient
	limit     int
	resetTime time.Duration
	ctx       context.Context
}

func NewRedisRateLimiter(ctx context.Context, rdb RedisClient, limit int, resetTime time.Duration) *RedisRateLimiter {

	return &RedisRateLimiter{
		client:    rdb,
		limit:     limit,
		resetTime: resetTime,
		ctx:       ctx,
	}
}

func (rl *RedisRateLimiter) AllowRequest(userID string) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s", userID)

	// Increment the count and set expiration
	count, err := rl.client.Incr(key).Result()
	if err != nil {
		return false, err
	}
	log.Println("Count", count)

	if count == 1 {
		rl.client.Expire(key, rl.resetTime)
	}

	if count > int64(rl.limit) {
		return false, nil // Request denied
	}

	return true, nil // Request allowed
}

package mocks

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/mock" // Optional if you want to track calls
)

// MockRedisClient is the mock struct
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Incr(key string) *redis.IntCmd {
	args := m.Called(key)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(key, expiration)
	return args.Get(0).(*redis.BoolCmd)
}

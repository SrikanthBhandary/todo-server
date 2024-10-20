package router

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/srikanthdoc/todo-server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestYourFunction(t *testing.T) {
	// Create a new instance of MockRedisClient
	mockRedis := &mocks.MockRedisClient{}

	// Create mock return values for Redis commands
	intCmd := redis.NewIntCmd(nil, 1)      // Simulate the response of Incr
	boolCmd := redis.NewBoolCmd(nil, true) // Simulate a successful response of Expire

	// Case 1: Test normal behavior
	mockRedis.On("Incr", "rate_limit:test-key").Return(intCmd)
	mockRedis.On("Expire", "rate_limit:test-key", 10*time.Second).Return(boolCmd)

	// Call the function under test and pass the mock
	limiter := NewRedisRateLimiter(context.TODO(), mockRedis, 2, 1*time.Second)
	result, err := limiter.AllowRequest("test-key")

	// Assertions
	assert.Equal(t, result, true)
	assert.Equal(t, err, nil)

	// Case 2: Test Expire failure
	boolCmdFail := redis.NewBoolCmd(nil, false) // Simulate an Expire failure (nil error but false value)

	mockRedis.On("Expire", "rate_limit:test-key", 10*time.Second).Return(boolCmdFail)

	// Call the function again and pass the mock with failing Expire
	result, err = limiter.AllowRequest("test-key")

	// Assertions: The request should still be allowed, but there might be an error in setting the expiration
	assert.Equal(t, result, true)
	assert.Nil(t, err, "There should be no error, even if Expire fails")
}

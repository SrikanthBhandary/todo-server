package router

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/srikanthdoc/todo-server/entity"
	"github.com/srikanthdoc/todo-server/mocks"

	"github.com/srikanthdoc/todo-server/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test for CreateToDo
func TestCreateUser_Success(t *testing.T) {
	mockToDoSvc := new(mocks.MockToDoService)
	mockUserSvc := new(mocks.MockUserService)
	jwtSvc := new(mocks.MockJWTValidator)
	emailSender := &mocks.MockEmailSender{}

	mockRedis := &mocks.MockRedisClient{}

	// Create mock return values for Redis commands
	intCmd := redis.NewIntCmd(nil, 1)      // Simulate the response of Incr
	boolCmd := redis.NewBoolCmd(nil, true) // Simulate a successful response of Expire

	// Case 1: Test normal behavior
	mockRedis.On("Incr", "rate_limit:1").Return(intCmd)
	mockRedis.On("Expire", "rate_limit:1", 10*time.Second).Return(boolCmd)
	ratelimiter := NewRedisRateLimiter(context.TODO(), mockRedis, 2, 1*time.Second)

	t.Run("TestCreateUser_SUCCESS", func(t *testing.T) {
		jobChannel := make(chan worker.Job, 10)
		pool := worker.NewWorkerPool(3, jobChannel)
		ctx, cancel := context.WithCancel(context.Background())

		defer cancel()
		pool.Init(ctx)

		// Create the router using the mock services
		r := NewRouter(mockToDoSvc, mockUserSvc, jwtSvc, ratelimiter, pool, emailSender)
		r.InitRoutes()

		user := &entity.User{
			UserName: "testuser",
			Password: "password",
		}

		mockUserSvc.On("CreateUser", mock.Anything, user).Return(nil)

		body, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockUserSvc.AssertExpectations(t)
	})

	t.Run("TestCreateUser_InvalidJSON", func(t *testing.T) {
		jobChannel := make(chan worker.Job, 10)
		pool := worker.NewWorkerPool(3, jobChannel)
		ctx, cancel := context.WithCancel(context.Background())

		defer cancel()
		pool.Init(ctx)

		// Create the router using the mock services
		r := NewRouter(mockToDoSvc, mockUserSvc, jwtSvc, ratelimiter, pool, emailSender)
		r.InitRoutes()

		user := &entity.User{
			UserName: "testuser",
			Password: "password",
		}

		mockUserSvc.On("CreateUser", mock.Anything, user).Return(nil)

		body := []byte("invalid json")
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("TestLoginUser_Success", func(t *testing.T) {
		jobChannel := make(chan worker.Job, 10)
		pool := worker.NewWorkerPool(3, jobChannel)
		ctx, cancel := context.WithCancel(context.Background())

		defer cancel()
		pool.Init(ctx)

		// Create the router using the mock services
		r := NewRouter(mockToDoSvc, mockUserSvc, jwtSvc, ratelimiter, pool, emailSender)
		r.InitRoutes()

		user := &entity.User{
			UserName: "testuser",
			Password: "hashedpassword",
		}

		loginRequest := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: "testuser",
			Password: "password",
		}

		mockUserSvc.On("GetUserByUserName", mock.Anything, loginRequest.Username).Return(user, nil)
		mockUserSvc.On("CheckPasswordHash", loginRequest.Password, user.Password).Return(true)
		jwtSvc.On("GenerateToken", user.UserID).Return("dummytoken", nil)

		body, _ := json.Marshal(loginRequest)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var result map[string]string
		err := json.NewDecoder(rr.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, "dummytoken", result["token"])
		mockUserSvc.AssertExpectations(t)
		jwtSvc.AssertExpectations(t)
	})
}

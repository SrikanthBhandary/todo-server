package router

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/srikanthbhandary/todo-server/entity"
	"github.com/srikanthbhandary/todo-server/mocks"

	"github.com/srikanthbhandary/todo-server/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test for CreateToDo
func TestCreateToDo(t *testing.T) {
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

	t.Run("TestCreateToDO_SUCCESS", func(t *testing.T) {
		jobChannel := make(chan worker.Job, 10)
		pool := worker.NewWorkerPool(3, jobChannel)
		ctx, cancel := context.WithCancel(context.Background())

		defer cancel()
		pool.Init(ctx)

		// Create the router using the mock services
		r := NewRouter(mockToDoSvc, mockUserSvc, jwtSvc, ratelimiter, pool, emailSender)
		r.InitRoutes()

		todo := &entity.ToDo{
			Title: "Test ToDo",
		}

		//Using mock.MatchedBy provides a targeted way to verify that your mock is being called correctly without enforcing overly strict requirements on the arguments. It’s a common practice in unit testing with mocks to ensure that tests remain clear and maintainable while still verifying the correct behavior of the code under test.
		mockToDoSvc.On("AddToDo", mock.Anything, mock.MatchedBy(func(t *entity.ToDo) bool {
			return t.Title == todo.Title // Check only the Title
		})).Return(nil)

		body, _ := json.Marshal(todo)
		req := httptest.NewRequest("POST", "/todos", bytes.NewBuffer(body))
		req = req.WithContext(context.WithValue(req.Context(), "userID", 1)) // Set userID in context
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer dummytoken") // Set Authorization Bearer token

		rr := httptest.NewRecorder()
		r.Router.ServeHTTP(rr, req)
		data, _ := io.ReadAll(rr.Result().Body)
		log.Println(string(data))
		defer rr.Result().Body.Close()
		assert.Equal(t, http.StatusCreated, rr.Code)

	})

	t.Run("TestGetAllToDos_SUCCESS", func(t *testing.T) {
		jobChannel := make(chan worker.Job, 10)
		pool := worker.NewWorkerPool(3, jobChannel)
		ctx, cancel := context.WithCancel(context.Background())

		defer cancel()
		pool.Init(ctx)

		// Create the router using the mock services
		r := NewRouter(mockToDoSvc, mockUserSvc, jwtSvc, ratelimiter, pool, emailSender)
		r.InitRoutes()

		todos := []entity.ToDo{
			{ToDoID: 1, Title: "ToDo 1", UserID: 1},
			{ToDoID: 2, Title: "ToDo 2", UserID: 1},
		}

		mockToDoSvc.On("GetAllTodos", mock.Anything, 1).Return(todos, nil)
		req := httptest.NewRequest("GET", "/todos", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userID", 1))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer dummytoken") // Set Authorization Bearer token

		rr := httptest.NewRecorder()
		r.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var result []entity.ToDo
		err := json.NewDecoder(rr.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Len(t, result, 2)

	})
}

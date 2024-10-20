package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/srikanthbhandary/todo-server/entity"
	"github.com/srikanthbhandary/todo-server/mocks"
	"github.com/srikanthbhandary/todo-server/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestToDo(t *testing.T) {
	mockRepo := new(mocks.MockToDoRepository)
	service := service.NewTodoService(mockRepo)
	t.Run("TestAddToDo_SUCCESS", func(t *testing.T) {
		todo := &entity.ToDo{
			Title: "Test ToDo",
		}
		mockRepo.On("AddToDo", mock.Anything, todo).Return(nil)
		err := service.AddToDo(context.Background(), todo)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("TestGetAllTodos_SUCCESS", func(t *testing.T) {
		todos := []entity.ToDo{
			{Title: "Todo 1", UserID: 1, DateTime: time.Now()},
			{Title: "Todo 2", UserID: 1, DateTime: time.Now()},
		}

		mockRepo.On("GetAllTodos", mock.Anything, 1).Return(todos, nil)

		result, err := service.GetAllTodos(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, todos, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestGetTodo_SUCCESS", func(t *testing.T) {
		todo := entity.ToDo{Title: "Todo 1", UserID: 1, DateTime: time.Now()}

		mockRepo.On("GetTodo", mock.Anything, 1, 1).Return(todo, nil)

		result, err := service.GetTodo(context.Background(), 1, 1)

		assert.NoError(t, err)
		assert.Equal(t, todo, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestDeleteToDo_SUCCESS", func(t *testing.T) {
		mockRepo.On("DeleteToDo", mock.Anything, 1, 1).Return(nil)

		err := service.DeleteToDo(context.Background(), 1, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestDeleteAllTodos_SUCCESS", func(t *testing.T) {
		mockRepo.On("DeleteAllTodos", mock.Anything, 1).Return(nil)

		err := service.DeleteAllTodos(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

}

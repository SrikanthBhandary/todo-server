package service

import (
	"context"
	"testing"

	"github.com/srikanthbhandary/todo-server/entity"
	"github.com/srikanthbhandary/todo-server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	t.Run("TestCreateUser_SUCCESS", func(t *testing.T) {
		user := &entity.User{
			UserName: "testuser",
			Password: "password", // Password should be hashed in real implementation
		}
		mockRepo.On("CreateUser", mock.Anything, user).Return(nil)

		err := service.CreateUser(context.Background(), user)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestGetUserByID_SUCCESS", func(t *testing.T) {
		user := &entity.User{UserID: 1, UserName: "testuser"}

		mockRepo.On("GetUserByID", mock.Anything, 1).Return(user, nil)

		result, err := service.GetUserByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestGetUserByUserName_SUCCESS", func(t *testing.T) {
		user := &entity.User{UserID: 1, UserName: "testuser"}

		mockRepo.On("GetUserByUserName", mock.Anything, "testuser").Return(user, nil)

		result, err := service.GetUserByUserName(context.Background(), "testuser")

		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestUpdateUser_SUCCESS", func(t *testing.T) {
		user := &entity.User{UserID: 1, UserName: "testuser"}

		mockRepo.On("UpdateUser", mock.Anything, user).Return(nil)

		err := service.UpdateUser(context.Background(), user)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("TestDeleteUser_SUCCESS", func(t *testing.T) {
		mockRepo.On("DeleteUser", mock.Anything, 1).Return(nil)

		err := service.DeleteUser(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

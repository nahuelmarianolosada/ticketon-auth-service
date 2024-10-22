package user

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"ticketon-auth-service/api/mocks"
	"ticketon-auth-service/api/model"
	"time"
)

func Test_gormDB_Create(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userToCreate := model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Time{},
		},
		FirstName: "testuser",
		Email:     "test@example.com",
	}

	t.Run("Success_Create", func(t *testing.T) {
		// Mock the Create method to return success
		mockRepo.On("Create", &userToCreate).Return(&gorm.DB{Error: nil})

		// Call the method
		result := mockRepo.Create(&userToCreate)

		// Assert that there are no errors
		assert.NoError(t, result.Error)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_Create", func(t *testing.T) {
		// Mock the Create method to return an error
		mockRepo := new(mocks.UserRepository)
		mockRepo.On("Create", &userToCreate).Return(&gorm.DB{Error: errors.New("duplicate entry")})

		// Call the method
		result := mockRepo.Create(&userToCreate)

		// Assert that an error is returned
		assert.Error(t, result.Error, "expected error but got none")
		assert.Equal(t, "duplicate entry", result.Error.Error())
		mockRepo.AssertExpectations(t)
	})
}

func Test_gormDB_First(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userID := "1"
	expectedUser := &model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Time{},
		},
		FirstName: "testuser",
		Email:     "test@example.com",
	}

	t.Run("Success_First", func(t *testing.T) {
		// Mock the First method to return success
		mockRepo.On("First", userID).Return(expectedUser, nil)

		// Call the method
		user, err := mockRepo.First(userID)

		// Assert that there are no errors and the user is returned
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_First", func(t *testing.T) {
		// Mock the First method to return an error
		mockRepo := new(mocks.UserRepository)
		mockRepo.On("First", userID).Return(nil, errors.New("user not found"))

		// Call the method
		user, err := mockRepo.First(userID)

		// Assert that an error is returned
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func Test_gormDB_Save(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userToSave := model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Time{},
		},
		FirstName: "testuser",
		Email:     "test@example.com",
	}

	t.Run("Success_Save", func(t *testing.T) {
		// Mock the Save method to return success
		mockRepo.On("Save", &userToSave).Return(&gorm.DB{Error: nil})

		// Call the method
		result := mockRepo.Save(&userToSave)

		// Assert that there are no errors
		assert.NoError(t, result.Error)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_Save", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		// Mock the Save method to return an error
		mockRepo.On("Save", &userToSave).Return(&gorm.DB{Error: errors.New("save failed")})

		// Call the method
		result := mockRepo.Save(&userToSave)

		// Assert that an error is returned
		assert.Error(t, result.Error)
		assert.Equal(t, "save failed", result.Error.Error())
		mockRepo.AssertExpectations(t)
	})
}

func Test_gormDB_Update(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userToUpdate := model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Time{},
		},
		FirstName: "testuser",
		Email:     "test@example.com",
	}

	t.Run("Success_Update", func(t *testing.T) {
		// Mock the Update method to return success
		mockRepo.On("Update", &userToUpdate).Return(&gorm.DB{Error: nil})

		// Call the method
		result := mockRepo.Update(&userToUpdate)

		// Assert that there are no errors
		assert.NoError(t, result.Error)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_Update", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		// Mock the Update method to return an error
		mockRepo.On("Update", &userToUpdate).Return(&gorm.DB{Error: errors.New("update failed")})

		// Call the method
		result := mockRepo.Update(&userToUpdate)

		// Assert that an error is returned
		assert.Error(t, result.Error)
		assert.Equal(t, "update failed", result.Error.Error())
		mockRepo.AssertExpectations(t)
	})
}

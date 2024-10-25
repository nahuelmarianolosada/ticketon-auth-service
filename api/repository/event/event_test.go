package event

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
	mockRepo := new(mocks.EventRepository)
	EventToCreate := model.EventBasic{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Time{},
		},
		Name: "testEvent",
	}

	t.Run("Success_Create", func(t *testing.T) {
		// Mock the Create method to return success
		mockRepo.On("Create", &EventToCreate).Return(&gorm.DB{Error: nil})

		// Call the method
		result := mockRepo.Create(&EventToCreate)

		// Assert that there are no errors
		assert.NoError(t, result.Error)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_Create", func(t *testing.T) {
		// Mock the Create method to return an error
		mockRepo := new(mocks.EventRepository)
		mockRepo.On("Create", &EventToCreate).Return(&gorm.DB{Error: errors.New("duplicate entry")})

		// Call the method
		result := mockRepo.Create(&EventToCreate)

		// Assert that an error is returned
		assert.Error(t, result.Error, "expected error but got none")
		assert.Equal(t, "duplicate entry", result.Error.Error())
		mockRepo.AssertExpectations(t)
	})
}

func Test_gormDB_First(t *testing.T) {
	mockRepo := new(mocks.EventRepository)
	EventID := "1"
	expectedEvent := &model.EventBasic{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Time{},
		},
		Name: "testEvent",
	}

	t.Run("Success_First", func(t *testing.T) {
		// Mock the First method to return success
		mockRepo.On("First", EventID).Return(expectedEvent, nil)

		// Call the method
		EventBasic, err := mockRepo.First(EventID)

		// Assert that there are no errors and the EventBasic is returned
		assert.NoError(t, err)
		assert.Equal(t, expectedEvent, EventBasic)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_First", func(t *testing.T) {
		// Mock the First method to return an error
		mockRepo := new(mocks.EventRepository)
		mockRepo.On("First", EventID).Return(nil, errors.New("EventBasic not found"))

		// Call the method
		EventBasic, err := mockRepo.First(EventID)

		// Assert that an error is returned
		assert.Error(t, err)
		assert.Nil(t, EventBasic)
		assert.Equal(t, "EventBasic not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func Test_gormDB_Save(t *testing.T) {
	mockRepo := new(mocks.EventRepository)
	EventToSave := model.EventBasic{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Time{},
		},
		Name: "testEvent",
	}

	t.Run("Success_Save", func(t *testing.T) {
		// Mock the Save method to return success
		mockRepo.On("Save", &EventToSave).Return(&gorm.DB{Error: nil})

		// Call the method
		result := mockRepo.Save(&EventToSave)

		// Assert that there are no errors
		assert.NoError(t, result.Error)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_Save", func(t *testing.T) {
		mockRepo := new(mocks.EventRepository)
		// Mock the Save method to return an error
		mockRepo.On("Save", &EventToSave).Return(&gorm.DB{Error: errors.New("save failed")})

		// Call the method
		result := mockRepo.Save(&EventToSave)

		// Assert that an error is returned
		assert.Error(t, result.Error)
		assert.Equal(t, "save failed", result.Error.Error())
		mockRepo.AssertExpectations(t)
	})
}

func Test_gormDB_Update(t *testing.T) {
	mockRepo := new(mocks.EventRepository)
	EventToUpdate := model.EventBasic{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Time{},
		},
		Name: "testEvent",
	}

	t.Run("Success_Update", func(t *testing.T) {
		// Mock the Update method to return success
		mockRepo.On("Update", &EventToUpdate).Return(&gorm.DB{Error: nil})

		// Call the method
		result := mockRepo.Update(&EventToUpdate)

		// Assert that there are no errors
		assert.NoError(t, result.Error)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_Update", func(t *testing.T) {
		mockRepo := new(mocks.EventRepository)
		// Mock the Update method to return an error
		mockRepo.On("Update", &EventToUpdate).Return(&gorm.DB{Error: errors.New("update failed")})

		// Call the method
		result := mockRepo.Update(&EventToUpdate)

		// Assert that an error is returned
		assert.Error(t, result.Error)
		assert.Equal(t, "update failed", result.Error.Error())
		mockRepo.AssertExpectations(t)
	})
}

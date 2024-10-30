package event

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"ticketon-auth-service/api/model"
	"ticketon-auth-service/api/repository"
)

// EventRepository defines the methods that the repository uses.
type EventRepository interface {
	Create(value model.EventBasic) *gorm.DB
	Save(value model.EventBasic) *gorm.DB
	Update(value model.EventBasic) *gorm.DB
	First(id string, user_id uint) (*model.EventBasic, error)
	Delete(value model.EventBasic) *gorm.DB
}

// Production DB that uses gorm
var DB EventRepository = &gormDB{}

type gormDB struct {
	*gorm.DB
}

func (db *gormDB) Create(value model.EventBasic) *gorm.DB {
	return repository.DB.Create(&value)
}

func (db *gormDB) Save(value model.EventBasic) *gorm.DB {
	return repository.DB.Save(&value)
}

func (db *gormDB) Update(value model.EventBasic) *gorm.DB {
	return repository.DB.Where("id = ?", value.ID).Updates(&value)
}

func (db *gormDB) First(id string, userID uint) (*model.EventBasic, error) {
	var existingEvent model.EventBasic

	// Convert the string to an integer (EventID)
	EventID, err := strconv.Atoi(id)
	if err != nil {
		// Return a descriptive error instead of nil
		return nil, fmt.Errorf("invalid Event ID: %s", id)
	}

	// Query the database with both EventID and user_id as conditions
	result := repository.DB.Where("id = ? AND user_id = ?", EventID, userID).First(&existingEvent)
	if result.Error != nil {
		return nil, result.Error
	}

	// Return the found Event and nil error
	return &existingEvent, nil
}

func (db *gormDB) Delete(value model.EventBasic) *gorm.DB {
	return repository.DB.Delete(&value, "id = ?", value.ID)
}

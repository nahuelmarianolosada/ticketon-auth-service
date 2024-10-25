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
	Create(value interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Update(value interface{}) *gorm.DB
	First(value string) (*model.EventBasic, error)
}

// Production DB that uses gorm
var DB EventRepository = &gormDB{}

type gormDB struct {
	*gorm.DB
}

func (db *gormDB) Create(value interface{}) *gorm.DB {
	return repository.DB.Create(value)
}

func (db *gormDB) Save(value interface{}) *gorm.DB {
	return repository.DB.Save(value)
}

func (db *gormDB) Update(value interface{}) *gorm.DB {
	if evt, ok := value.(*model.EventBasic); ok {
		return repository.DB.Where("id = ?", evt.ID).Updates(value)
	}
	return repository.DB.Where("id = ?").Updates(value)
}

func (db *gormDB) First(value string) (*model.EventBasic, error) {
	var existingEvent model.EventBasic

	// Convert the string to an integer (EventID)
	EventID, err := strconv.Atoi(value)
	if err != nil {
		// Return a descriptive error instead of nil
		return nil, fmt.Errorf("invalid Event ID: %s", value)
	}

	// Query the database for the Event with the given ID
	result := repository.DB.First(&existingEvent, EventID)
	if result.Error != nil {
		// Return the database error if the Event is not found or any other error occurs
		return nil, result.Error
	}

	// Return the found Event and nil error
	return &existingEvent, nil
}

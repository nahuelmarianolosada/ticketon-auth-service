package user

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"ticketon-auth-service/api/model"
	"ticketon-auth-service/api/repository"
)

// DBInterface defines the methods that the repository uses.
type DBInterface interface {
	Create(value interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Update(value interface{}) *gorm.DB
	First(value string) (*model.User, error)
}

// Production DB that uses gorm
var DB DBInterface = &gormDB{}

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
	if usr, ok := value.(*model.User); ok {
		return repository.DB.Where("id = ?", usr.ID).Updates(value)
	}
	return repository.DB.Where("id = ?").Updates(value)
}

func (db *gormDB) First(value string) (*model.User, error) {
	var existingUser model.User

	// Convert the string to an integer (userID)
	userID, err := strconv.Atoi(value)
	if err != nil {
		// Return a descriptive error instead of nil
		return nil, fmt.Errorf("invalid user ID: %s", value)
	}

	// Query the database for the user with the given ID
	result := repository.DB.First(&existingUser, userID)
	if result.Error != nil {
		// Return the database error if the user is not found or any other error occurs
		return nil, result.Error
	}

	// Return the found user and nil error
	return &existingUser, nil
}

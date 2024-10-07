package user

import (
	"gorm.io/gorm"
	"ticketon-auth-service/api/repository"
)

// DBInterface defines the methods that the repository uses.
type DBInterface interface {
	Create(value interface{}) *gorm.DB
}

// Production DB that uses gorm
var DB DBInterface = &gormDB{}

type gormDB struct {
	*gorm.DB
}

func (db *gormDB) Create(value interface{}) *gorm.DB {
	return repository.DB.Create(value)
}

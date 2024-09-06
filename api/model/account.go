package model

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	UserID          uint
	Cvu             *string `json:"cvu"`
	Alias           *string `json:"alias"`
	AvailableAmount string  `json:"available_amount"`
}

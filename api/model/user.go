package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUser interface {
	HashPassword(password string) error
	CheckPassword(providedPassword string) error
}

type User struct {
	gorm.Model
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname" binding:"required"`
	Dni       int    `json:"dni" binding:"required"`
	Email     string `json:"email" gorm:"unique" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
}

// DTO for binding JSON and validation
type CreateUserRequest struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname" binding:"required"`
	Dni       int    `json:"dni" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Phone     string `json:"phone" binding:"required"`

	HashPasswordFunc func(password string) error `json:"-"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (user *CreateUserRequest) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

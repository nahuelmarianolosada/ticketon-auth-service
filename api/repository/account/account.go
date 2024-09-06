package account

import (
	"ticketon-auth-service/api/model"
	"ticketon-auth-service/api/repository"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func GetByID(accountID int) (*model.Account, error) {
	var account model.Account
	result := repository.DB.First(&account, accountID)

	if result.Error != nil {
		fmt.Printf("ERROR %v", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, result.Error
	}
	return &account, nil
}

func GetByUserID(userId int) (*model.Account, error) {
	var account model.Account
	result := repository.DB.First(&account, "user_id = ?", userId)

	if result.Error != nil {
		fmt.Printf("ERROR %v", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, result.Error
	}
	return &account, nil
}

func GetByAliasCvu(alias, cvu string) (*model.Account, error) {
	var account model.Account
	result := repository.DB.First(&account, "alias = ? OR cvu = ?", alias, cvu)

	if result.Error != nil {
		fmt.Printf("ERROR %v", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, result.Error
	}
	return &account, nil
}

func Create(account model.Account) (*model.Account, error) {
	tx := repository.DB.Create(&account)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &account, nil
}

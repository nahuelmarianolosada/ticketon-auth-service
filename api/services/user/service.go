package user

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"ticketon-auth-service/api/model"
	accountRepo "ticketon-auth-service/api/repository/account"
	userRepo "ticketon-auth-service/api/repository/user"
	"time"
)

func CreateUser(ctx context.Context, bodyReq model.CreateUserRequest) (*model.CreateUserResponse, error) {
	newDefaultAccount := model.Account{AvailableAmount: "0"}

	userToCreate := model.User{
		Model: gorm.Model{
			CreatedAt: time.Now(),
		},
		FirstName: bodyReq.FirstName,
		LastName:  bodyReq.LastName,
		Dni:       bodyReq.Dni,
		Email:     bodyReq.Email,
		Password:  bodyReq.Password,
		Phone:     bodyReq.Phone,
	}

	record := userRepo.DB.Create(&userToCreate)
	if record.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(record.Error, &mysqlErr) {
			// Check if it's a duplicate entry error (error code 1062)
			if mysqlErr.Number == 1062 {
				return nil, model.ApiError{Message: "Email already exists", Err: record.Error}
			}
		}
		return nil, model.ApiError{Message: record.Error.Error(), Err: record.Error}

	}

	newDefaultAccount.UserID = userToCreate.ID

	accountCreated, err := accountRepo.DB.Create(newDefaultAccount)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			// 1062 is the error code for a duplicate entry
			return nil, model.ApiError{Message: "A user with this email already exists.", Err: err}
		}
		return nil, model.ApiError{Message: err.Error(), Err: err}
	}

	return &model.CreateUserResponse{UserID: userToCreate.ID, AccountID: accountCreated.ID, Email: userToCreate.Email}, nil
}

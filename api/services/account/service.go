package account

import (
	"context"
	"github.com/go-sql-driver/mysql"
	"ticketon-auth-service/api/model"
	accountRepo "ticketon-auth-service/api/repository/account"
)

func CreateAccount(ctx context.Context, userID uint) (*model.Account, error) {
	newDefaultAccount := model.Account{AvailableAmount: "0", UserID: userID}
	accountCreated, err := accountRepo.DB.Create(newDefaultAccount)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			// 1062 is the error code for a duplicate entry
			return nil, model.ApiError{Message: "A user with this email already exists.", Err: err}
		}
		return nil, model.ApiError{Message: err.Error(), Err: err}
	}
	return accountCreated, nil
}

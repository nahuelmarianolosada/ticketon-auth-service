package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"ticketon-auth-service/api/model"
	_ "ticketon-auth-service/api/model"
	accountRepo "ticketon-auth-service/api/repository/account"
)

func FindAccount(c *gin.Context) {
	userId := GetUserIDFromJWT(c)
	if userId == nil {
		return
	}
	userFound, err := accountRepo.GetByUserID(*userId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.AbortWithStatusJSON(http.StatusNotFound, model.ApiError{Message: err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, userFound)
}

func ValidateAccountWithToken(c *gin.Context, accountID int) bool {
	//Validamos que la cuenta solicitada coincida con la cuenta del usuario recibida en el token
	userId := GetUserIDFromJWT(c)
	if userId == nil {
		c.AbortWithStatusJSON(http.StatusForbidden, model.ApiError{Message: "UserID from token is required"})
		return false
	}
	acc, err := accountRepo.GetByUserID(*userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiError{Message: err.Error()})
		return false
	} else if int(acc.ID) != accountID {
		c.AbortWithStatusJSON(http.StatusForbidden, model.ApiError{Message: "Not allowed to access account"})
		return false
	}

	return true
}

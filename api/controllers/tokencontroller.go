package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"ticketon-auth-service/api/middlewares/auth"
	"ticketon-auth-service/api/model"
	"ticketon-auth-service/api/repository"
)

type TokenRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func GenerateToken(context *gin.Context) {
	var request TokenRequest
	var user model.User
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, model.ApiError{Message: err.Error()})
		context.Abort()
		return
	}
	// check if email exists and password is correct
	record := repository.DB.Where("email = ?", request.Email).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, model.ApiError{Message: record.Error.Error()})
		context.Abort()
		return
	}
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, model.ApiError{Message: "invalid credentials"})
		context.Abort()
		return
	}
	tokenString, err := auth.GenerateJWT(user.Email, strconv.Itoa(int(user.ID)))
	if err != nil {
		context.JSON(http.StatusInternalServerError, model.ApiError{Message: err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}

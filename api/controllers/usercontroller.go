package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"ticketon-auth-service/api/middlewares/auth"
	"ticketon-auth-service/api/model"
	accountRepo "ticketon-auth-service/api/repository/account"
	userRepo "ticketon-auth-service/api/repository/user"
	"time"
)

func GetUserIDFromJWT(c *gin.Context) *int {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ApiError{Message: "request does not contain an access token"})
		return nil
	}

	claims, err := auth.GetClaims(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ApiError{Message: err.Error()})
		return nil
	}

	userID := claims.Username
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ApiError{Message: "user_id is required"})
		return nil
	}

	userIdInt, err := strconv.Atoi(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ApiError{Message: "user_id is not a number"})
		return nil
	}

	return &userIdInt
}

func RegisterUser(c *gin.Context) {
	var user model.CreateUserRequest
	newDefaultAccount := model.Account{AvailableAmount: "0"}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ApiError{Message: err.Error()})
		return
	}

	// Ensure HashPasswordFunc is set to the default if not already set (useful for tests)
	if user.HashPasswordFunc == nil {
		user.HashPasswordFunc = user.HashPassword
	}

	if err := user.HashPassword(user.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiError{Message: err.Error()})
		return
	}

	userToCreate := model.User{
		Model: gorm.Model{
			CreatedAt: time.Now(),
		},
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Dni:       user.Dni,
		Email:     user.Email,
		Password:  user.Password,
		Phone:     user.Phone,
	}

	record := userRepo.DB.Create(userToCreate)
	if record.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiError{Message: record.Error.Error()})
		return
	}

	newDefaultAccount.UserID = user.ID

	accountCreated, err := accountRepo.DB.Create(newDefaultAccount)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": user.ID, "account_id": accountCreated.ID, "email": user.Email})
}

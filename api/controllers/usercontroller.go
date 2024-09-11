package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"ticketon-auth-service/api/middlewares/auth"
	"ticketon-auth-service/api/model"
	"ticketon-auth-service/api/repository"
	accountRepo "ticketon-auth-service/api/repository/account"
)

func RegisterUser(c *gin.Context) {
	var user model.CreateUserRequest
	newDefaultAccount := model.Account{AvailableAmount: "0"}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := user.HashPassword(user.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	record := repository.DB.Create(user)
	if record.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		return
	}

	newDefaultAccount.UserID = user.ID

	accountCreated, err := accountRepo.Create(newDefaultAccount)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": user.ID, "account_id": accountCreated.ID, "email": user.Email})
}

func GetUserIDFromJWT(c *gin.Context) *int {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "request does not contain an access token"})
		return nil
	}

	claims, err := auth.GetClaims(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return nil
	}

	userID := claims.Username
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id is required"})
		return nil
	}

	userIdInt, err := strconv.Atoi(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user_id is not a number"})
		return nil
	}

	return &userIdInt
}

package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
	"ticketon-auth-service/api/middlewares/auth"
	"ticketon-auth-service/api/model"
	accountRepo "ticketon-auth-service/api/repository/account"
	userRepo "ticketon-auth-service/api/repository/user"
	accountService "ticketon-auth-service/api/services/account"
	userService "ticketon-auth-service/api/services/user"
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

	userCreated, err := userService.CreateUser(c, user)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			// Check if it's a duplicate entry error (error code 1062)
			if mysqlErr.Number == 1062 {
				c.AbortWithStatusJSON(http.StatusConflict, model.ApiError{Message: "Email already exists"})
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiError{Message: err.Error()})
		return
	}

	accountService.CreateAccount(c, userCreated.UserID)
	accountCreated, err := accountRepo.DB.Create(newDefaultAccount)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			// Check if it's a duplicate entry error (error code 1062)
			if mysqlErr.Number == 1062 {
				c.AbortWithStatusJSON(http.StatusConflict, model.ApiError{Message: "Email already exists"})
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": userCreated.UserID, "account_id": accountCreated.ID, "email": user.Email})
}

func UpdateUser(c *gin.Context) {
	// Get the user ID from the URL path
	userID := c.Param("id")

	// Check if the user exists
	var existingUser model.User
	if _, err := userRepo.DB.First(userID); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, model.ApiError{Message: "User not found. " + err.Error()})
		return
	}

	// Bind the request body to the CreateUserRequest struct
	var updatedUserData model.CreateUserRequest
	if err := c.ShouldBindJSON(&updatedUserData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ApiError{Message: err.Error()})
		return
	}

	// If there's a password in the update request, hash it before saving
	if updatedUserData.Password != "" {
		if updatedUserData.HashPasswordFunc == nil {
			updatedUserData.HashPasswordFunc = updatedUserData.HashPassword
		}

		if err := updatedUserData.HashPasswordFunc(updatedUserData.Password); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiError{Message: "Password hashing failed"})
			return
		}
		existingUser.Password = updatedUserData.Password // Update password only if provided
	}

	// Update fields that are allowed to be updated
	existingUser.FirstName = updatedUserData.FirstName
	existingUser.LastName = updatedUserData.LastName
	existingUser.Dni = updatedUserData.Dni
	existingUser.Email = updatedUserData.Email
	existingUser.Phone = updatedUserData.Phone

	// Save the updated user to the database
	if err := userRepo.DB.Update(&existingUser).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiError{Message: err.Error()})
		return
	}

	// Return the updated user data in the response
	c.JSON(http.StatusOK, gin.H{
		"user_id":   existingUser.ID,
		"firstname": existingUser.FirstName,
		"lastname":  existingUser.LastName,
		"email":     existingUser.Email,
		"phone":     existingUser.Phone,
	})
}

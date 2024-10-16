package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"ticketon-auth-service/api/mocks"
	"ticketon-auth-service/api/model"
	accountRepo "ticketon-auth-service/api/repository/account"
	userRepo "ticketon-auth-service/api/repository/user"
	"time"
)

func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("BadRequest_ShouldReturn400", func(t *testing.T) {
		// Prepare an invalid request body (missing required fields, incorrect JSON format, etc.)
		invalidRequestBody := `{"invalidField": "invalidValue"}`

		// Use the mockery-generated mocks
		mockDB := new(mocks.UserRepository)
		mockAccountRepo := new(mocks.AccountRepository)

		// Inject the mock DB into the repository
		userRepo.DB = mockDB
		accountRepo.DB = mockAccountRepo

		// Create the router and register the handler
		router := gin.Default()
		router.POST("/api/users", RegisterUser)

		// Serve the invalid request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer([]byte(invalidRequestBody)))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "message") // Expect an error message in the response
	})

	t.Run("AccountCreateFailure_ShouldReturn500", func(t *testing.T) {
		mockUserRequest := model.CreateUserRequest{
			FirstName: "John",
			LastName:  "Doe",
			Dni:       1,
			Email:     "test@example.com",
			Password:  "password",
			Phone:     "+1234567890",
		}

		// Use mockery-generated mocks
		mockDB := new(mocks.UserRepository)
		mockAccountRepo := new(mocks.AccountRepository)

		// Mock the DB to successfully create the user
		mockDB.On("Create", mock.AnythingOfType("*model.User")).Return(&gorm.DB{})

		// Mock Account Repo to simulate a failure during account creation
		mockAccountRepo.On("Create", mock.AnythingOfType("model.Account")).Return(nil, errors.New("account create error"))

		// Inject the mock DB into the repository
		userRepo.DB = mockDB
		accountRepo.DB = mockAccountRepo

		// Create the router and register the handler
		router := gin.Default()
		router.POST("/api/users", RegisterUser)

		// Prepare the request body
		reqBody, _ := json.Marshal(mockUserRequest)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "message") // Expect an error message in the response
	})

	t.Run("SuccessfulRegistration_ShouldReturn201", func(t *testing.T) {
		mockUserRequest := model.CreateUserRequest{
			FirstName: "John",
			LastName:  "Doe",
			Dni:       1,
			Email:     "test@example.com",
			Password:  mock.Anything,
			Phone:     "+1234567890",
		}

		// Use mockery-generated mocks
		mockDB := new(mocks.UserRepository)
		mockAccountRepo := new(mocks.AccountRepository)

		// Mocking DB and Account creation
		mockDB.On("Create", mock.AnythingOfType("*model.User")).Return(&gorm.DB{})
		mockAccountRepo.On("Create", mock.AnythingOfType("model.Account")).Return(&model.Account{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserID:          1,
			AvailableAmount: "0",
		}, nil)

		// Inject the mock DB into repository
		userRepo.DB = mockDB
		accountRepo.DB = mockAccountRepo

		// Create the router and register the handler
		router := gin.Default()
		router.POST("/api/users", RegisterUser)

		// Prepare the request body
		reqBody, _ := json.Marshal(mockUserRequest)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json") // Ensure the correct content type is set

		// Serve the request through Gin's router
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "user_id")
		assert.Contains(t, w.Body.String(), "account_id")
	})
}

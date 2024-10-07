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
	"ticketon-auth-service/api/model"
	accountRepo "ticketon-auth-service/api/repository/account"
	userRepo "ticketon-auth-service/api/repository/user"
	"time"
)

// Mock dependencies
type MockDB struct {
	mock.Mock
}

type MockAccountRepo struct {
	mock.Mock
}

type MockUserHash struct {
	mock.Mock
}

// Mock for DB Create method
func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return &gorm.DB{
		Error: args.Error(0),
	}
}

// Mock for accountRepo.Create method
func (m *MockAccountRepo) Create(account model.Account) (*model.Account, error) {
	args := m.Called(account)

	// Safely handle nil value for the returned account
	if args.Get(0) != nil {
		return args.Get(0).(*model.Account), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserHash) HashPassword(pass string) error {
	args := m.Called(pass)

	// Safely handle nil value for the returned account
	if args.Get(0) != nil {
		return args.Error(1)
	}
	return args.Error(1)
}

func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("BadRequest_ShouldReturn400", func(t *testing.T) {
		// Prepare an invalid request body (missing required fields, incorrect JSON format, etc.)
		invalidRequestBody := `{"invalidField": "invalidValue"}`

		// Mock the DB and Account Repo, but they won't be called in this case
		mockDB := new(MockDB)
		mockAccountRepo := new(MockAccountRepo)

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

		// Mock the DB to successfully create the user
		mockDB := new(MockDB)
		mockDB.On("Create", mock.AnythingOfType("model.User")).Return(nil)

		// Mock Account Repo to simulate a failure during account creation
		mockAccountRepo := new(MockAccountRepo)
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

		mockDB := new(MockDB)
		mockAccountRepo := new(MockAccountRepo)

		// Mocking DB and Account creation
		mockDB.On("Create", mock.AnythingOfType("model.User")).Return(nil)
		mockAccountRepo.On("Create", mock.AnythingOfType("model.Account")).Return(&model.Account{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserID:          1,
			Cvu:             nil,
			Alias:           nil,
			AvailableAmount: "",
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

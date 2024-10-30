package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"ticketon-auth-service/api/middlewares/auth"
	"ticketon-auth-service/api/model"
	"ticketon-auth-service/api/repository"
)

// Setup function to initialize a mock database before running tests
func setupTestDB() {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&model.User{}) // Auto-migrate User model if necessary
	repository.DB = db
}

// Mocking the database response function to simulate GORM behaviors
func mockDatabaseResponse(user *model.User, err error) {
	repository.DB.Callback().Query().Replace("gorm:query", func(db *gorm.DB) {
		db.Statement.Dest = user
		db.Error = err
	})
}

func TestGenerateToken(t *testing.T) {
	// Initialize test database
	setupTestDB()
	gin.SetMode(gin.TestMode)

	// Save the original GenerateJWT function and reset after the test
	originalGenerateJWT := auth.GenerateJWT
	defer func() { auth.GenerateJWT = originalGenerateJWT }()

	// Table of test cases
	tests := []struct {
		name         string
		body         interface{}
		expectedCode int
		expectedMsg  string
		mockUser     *model.User
		mockDBError  error
		mockJWTError error
	}{
		{
			name:         "Missing fields",
			body:         map[string]string{},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Key: 'TokenRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag",
		},
		{
			name:         "Non-existent email",
			body:         map[string]string{"email": "notfound@example.com", "password": "somepassword"},
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "record not found",
			mockDBError:  gorm.ErrRecordNotFound,
		},
		{
			name:         "Incorrect password",
			body:         map[string]string{"email": "test@example.com", "password": "wrongpassword"},
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "invalid credentials",
			mockUser:     &model.User{Email: "test@example.com", Password: "hashedpassword"},
		},
		/*{
			name:         "Successful token generation",
			body:         map[string]string{"email": "test@example.com", "password": "correctpassword"},
			expectedCode: http.StatusOK,
			expectedMsg:  "mocked.token.string",
			mockUser:     &model.User{Email: "test@example.com", Password: "hashedpassword"},
		},*/
		{
			name:         "JWT generation failure",
			body:         map[string]string{"email": "fail@example.com", "password": "correctpassword"},
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "invalid credentials",
			mockUser:     &model.User{Email: "fail@example.com", Password: "hashedpassword"},
			mockJWTError: errors.New("token generation failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock for GenerateJWT
			auth.GenerateJWT = func(email, userID string) (string, error) {
				if tt.mockJWTError != nil {
					return "", tt.mockJWTError
				}
				return "mocked.token.string", nil
			}

			// Mock the database query for user retrieval
			mockDatabaseResponse(tt.mockUser, tt.mockDBError)

			// Set up request
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/token", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resp)
			ctx.Request = req

			// Execute function
			GenerateToken(ctx)

			// Validate response
			assert.Equal(t, tt.expectedCode, resp.Code)
			if tt.expectedCode == http.StatusOK {
				var response map[string]string
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMsg, response["token"])
			} else {
				var response model.ApiError
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.Message, tt.expectedMsg)
			}
		})
	}
}

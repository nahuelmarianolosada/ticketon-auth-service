package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth_NoToken(t *testing.T) {
	// Set up the router with the Auth middleware
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Auth(func(token string) error {
		return nil
	}))
	router.GET("/test", func(c *gin.Context) {
		c.String(200, "passed")
	})

	// Create a request without an Authorization header
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "request does not contain an access token")
}

func TestAuth_InvalidToken(t *testing.T) {
	// Set up the router with the Auth middleware
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Auth(func(token string) error {
		return assert.AnError // Simulate token validation error
	}))
	router.GET("/test", func(c *gin.Context) {
		c.String(200, "passed")
	})

	// Create a request with an invalid token in the Authorization header
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "invalid_token")
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), assert.AnError.Error())
}

func TestAuth_ValidToken(t *testing.T) {
	// Set up the router with the Auth middleware
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Auth(func(token string) error {
		return nil // Simulate valid token
	}))
	router.GET("/test", func(c *gin.Context) {
		c.String(200, "passed")
	})

	// Create a request with a valid token in the Authorization header
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "valid_token")
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "passed", w.Body.String())
}

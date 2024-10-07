package auth

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestGenerateJWT(t *testing.T) {
	// Set a dummy secret key for testing
	os.Setenv("JWT_SK", "testsecret")
	jwtKey = []byte(os.Getenv("JWT_SK"))

	email := "test@example.com"
	username := "testuser"

	tokenString, err := GenerateJWT(email, username)

	assert.NoError(t, err, "Expected no error while generating token")
	assert.NotEmpty(t, tokenString, "Expected token string to be non-empty")
}

func TestValidateToken(t *testing.T) {
	os.Setenv("JWT_SK", "testsecret")
	jwtKey = []byte(os.Getenv("JWT_SK"))

	email := "test@example.com"
	username := "testuser"
	tokenString, err := GenerateJWT(email, username)
	assert.NoError(t, err, "Expected no error while generating token")

	// Test valid token
	err = ValidateToken(tokenString)
	assert.NoError(t, err, "Expected no error while validating token")

	// Test expired token by manipulating the token expiration time
	claims := &JWTClaim{
		Email:    email,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(), // Set token to be expired
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredTokenString, _ := expiredToken.SignedString(jwtKey)

	err = ValidateToken(expiredTokenString)
	assert.Error(t, err, "Expected error for expired token")
	assert.Equal(t, "token is expired by 1h0m0s", err.Error(), "Expected 'token expired' error")
}

func TestGetClaims(t *testing.T) {
	os.Setenv("JWT_SK", "testsecret")
	jwtKey = []byte(os.Getenv("JWT_SK"))

	email := "test@example.com"
	username := "testuser"
	tokenString, err := GenerateJWT(email, username)
	assert.NoError(t, err, "Expected no error while generating token")

	claims, err := GetClaims(tokenString)
	assert.NoError(t, err, "Expected no error while getting claims")
	assert.Equal(t, email, claims.Email, "Expected email to match")
	assert.Equal(t, username, claims.Username, "Expected username to match")

	// Test invalid token
	invalidToken := "invalid.token.string"
	claims, err = GetClaims(invalidToken)
	assert.Error(t, err, "Expected error for invalid token")
	assert.Nil(t, claims, "Expected nil claims for invalid token")
}

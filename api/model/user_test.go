package model

import (
	"github.com/stretchr/testify/assert"
	bcrypt "golang.org/x/crypto/bcrypt"
	"testing"
)

func TestUser_HashPassword(t *testing.T) {
	t.Run("Success_HashPassword", func(t *testing.T) {
		// Create a user object
		user := User{}

		// Call HashPassword method
		err := user.HashPassword("password123")

		// Assert that there is no error
		assert.NoError(t, err, "Expected no error while hashing password")

		// Assert that password is not the plain text password
		assert.NotEqual(t, "password123", user.Password, "Hashed password should not equal plain text password")

		// Assert that password has been hashed correctly
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
		assert.NoError(t, err, "Expected bcrypt comparison to succeed")
	})
}

func TestUser_CheckPassword(t *testing.T) {
	t.Run("Success_CheckPassword", func(t *testing.T) {
		// Create a user object with a hashed password
		user := User{}
		err := user.HashPassword("password123")
		assert.NoError(t, err, "Expected no error while hashing password")

		// Call CheckPassword method with the correct password
		err = user.CheckPassword("password123")

		// Assert that there is no error
		assert.NoError(t, err, "Expected no error while checking correct password")
	})

	t.Run("Failure_CheckPassword", func(t *testing.T) {
		// Create a user object with a hashed password
		user := User{}
		err := user.HashPassword("password123")
		assert.NoError(t, err, "Expected no error while hashing password")

		// Call CheckPassword method with an incorrect password
		err = user.CheckPassword("wrongpassword")

		// Assert that there is an error
		assert.Error(t, err, "Expected an error while checking incorrect password")
		assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err, "Expected mismatch error")
	})
}

func TestCreateUserRequest_HashPassword(t *testing.T) {
	t.Run("Success_HashPassword", func(t *testing.T) {
		// Create a CreateUserRequest object
		userRequest := CreateUserRequest{}

		// Call HashPassword method
		err := userRequest.HashPassword("password123")

		// Assert that there is no error
		assert.NoError(t, err, "Expected no error while hashing password")

		// Assert that password is not the plain text password
		assert.NotEqual(t, "password123", userRequest.Password, "Hashed password should not equal plain text password")

		// Assert that password has been hashed correctly
		err = bcrypt.CompareHashAndPassword([]byte(userRequest.Password), []byte("password123"))
		assert.NoError(t, err, "Expected bcrypt comparison to succeed")
	})
}

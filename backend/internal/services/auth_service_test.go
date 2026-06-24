package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordComparison(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	err = bcrypt.CompareHashAndPassword(hash, []byte("testpassword"))
	assert.NoError(t, err, "should match correct password")

	err = bcrypt.CompareHashAndPassword(hash, []byte("wrongpassword"))
	assert.Error(t, err, "should reject wrong password")
}

func TestBcryptHashConstant(t *testing.T) {
	password := "user123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	assert.True(t, bcrypt.DefaultCost >= 10, "bcrypt cost should be at least 10")

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	assert.NoError(t, err)

	err = bcrypt.CompareHashAndPassword(hash, []byte("wrong"))
	assert.Error(t, err)
}

func TestHashUniqueness(t *testing.T) {
	password := "samepassword"
	hash1, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hash2, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	assert.NotEqual(t, string(hash1), string(hash2), "hashes should be unique due to salt")
}

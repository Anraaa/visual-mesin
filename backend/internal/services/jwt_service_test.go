package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTGenerateAndValidate(t *testing.T) {
	svc := NewJWTService("test-secret-key", time.Hour)

	token, err := svc.GenerateToken(1, "admin", "admin@test.com")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := svc.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "admin", claims.UserLevel)
	assert.Equal(t, "admin@test.com", claims.Email)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
}

func TestJWTExpiredToken(t *testing.T) {
	svc := NewJWTService("test-secret-key", -time.Hour)

	token, err := svc.GenerateToken(1, "admin", "admin@test.com")
	require.NoError(t, err)

	_, err = svc.ValidateToken(token)
	assert.Error(t, err)
}

func TestJWTInvalidSignature(t *testing.T) {
	svc1 := NewJWTService("secret-1", time.Hour)
	svc2 := NewJWTService("secret-2", time.Hour)

	token, err := svc1.GenerateToken(1, "admin", "admin@test.com")
	require.NoError(t, err)

	_, err = svc2.ValidateToken(token)
	assert.Error(t, err)
}

func TestJWTInvalidTokenString(t *testing.T) {
	svc := NewJWTService("test-secret-key", time.Hour)

	_, err := svc.ValidateToken("not-a-valid-jwt")
	assert.Error(t, err)

	_, err = svc.ValidateToken("")
	assert.Error(t, err)
}

func TestJWTTokenClaims(t *testing.T) {
	svc := NewJWTService("test-secret-key", time.Hour)

	token, err := svc.GenerateToken(42, "prod", "user@test.com")
	require.NoError(t, err)

	claims, err := svc.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, uint(42), claims.UserID)
	assert.Equal(t, "prod", claims.UserLevel)
	assert.Equal(t, "user@test.com", claims.Email)
}

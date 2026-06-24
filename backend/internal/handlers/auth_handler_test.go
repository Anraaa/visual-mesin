package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/repository"
	"github.com/anraaa/visual-mesin/internal/services"
)

func setupIntegrationTest(t *testing.T) (*gin.Engine, *services.JWTService) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	createTables(t, db)

	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := models.User{
		Email:     strPtr("admin@admin.com"),
		UserName:  "Admin",
		NIP:       strPtr("m26-134"),
		Password:  string(hash),
		UserLevel: "admin",
	}
	err = db.Create(&user).Error
	require.NoError(t, err)

	user2 := models.User{
		Email:     strPtr("user@visualmesin.com"),
		UserName:  "User Produksi",
		NIP:       strPtr("26-133"),
		Password:  string(hash),
		UserLevel: "prod",
	}
	err = db.Create(&user2).Error
	require.NoError(t, err)

	userRepo := repository.NewUserRepository(db)
	jwtSvc := services.NewJWTService("test-secret-key", 24*time.Hour)
	authSvc := services.NewAuthService(userRepo, jwtSvc)
	activityLogSvc := services.NewActivityLogService(db)
	authHandler := NewAuthHandler(authSvc, activityLogSvc)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	api := r.Group("/api/v1")
	api.POST("/auth/login", authHandler.Login)

	return r, jwtSvc
}

func TestAuthHandler_Login_Success(t *testing.T) {
	r, _ := setupIntegrationTest(t)

	body, _ := json.Marshal(models.LoginRequest{
		Email:    "admin@admin.com",
		Password: "password",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestAuthHandler_Login_InvalidPassword(t *testing.T) {
	r, _ := setupIntegrationTest(t)

	body, _ := json.Marshal(models.LoginRequest{
		Email:    "admin@admin.com",
		Password: "wrongpassword",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.False(t, resp.Success)
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	r, _ := setupIntegrationTest(t)

	body, _ := json.Marshal(models.LoginRequest{
		Email:    "unknown@email.com",
		Password: "password",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_Login_MissingFields(t *testing.T) {
	r, _ := setupIntegrationTest(t)

	tests := []struct {
		name string
		body map[string]string
	}{
		{"missing password", map[string]string{"email": "admin@admin.com"}},
		{"missing email", map[string]string{"password": "password"}},
		{"empty body", map[string]string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestAuthHandler_Login_WithNIP(t *testing.T) {
	r, _ := setupIntegrationTest(t)

	body, _ := json.Marshal(models.LoginRequest{
		Email:    "m26-134",
		Password: "password",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_Login_ReturnsToken(t *testing.T) {
	r, _ := setupIntegrationTest(t)

	body, _ := json.Marshal(models.LoginRequest{
		Email:    "admin@admin.com",
		Password: "password",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	respData, ok := resp.Data.(map[string]interface{})
	require.True(t, ok)

	token, ok := respData["token"].(string)
	require.True(t, ok)
	assert.NotEmpty(t, token)
}

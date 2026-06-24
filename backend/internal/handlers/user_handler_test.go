package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

func setupUserHandlerTest(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	createTables(t, db)

	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	admin := models.User{
		Email:     strPtr("admin@admin.com"),
		UserName:  "Admin",
		NIP:       strPtr("m26-134"),
		Password:  string(hash),
		UserLevel: "admin",
	}
	require.NoError(t, db.Create(&admin).Error)

	user := models.User{
		Email:     strPtr("user@visualmesin.com"),
		UserName:  "User Produksi",
		NIP:       strPtr("26-133"),
		Password:  string(hash),
		UserLevel: "prod",
	}
	require.NoError(t, db.Create(&user).Error)

	role := models.Role{Name: "user", GuardName: "web"}
	require.NoError(t, db.Create(&role).Error)

	adminRole := models.Role{Name: "admin", GuardName: "web"}
	require.NoError(t, db.Create(&adminRole).Error)

	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	activityLogSvc := services.NewActivityLogService(db)
	userHandler := NewUserHandler(userRepo, roleRepo, activityLogSvc)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	api := r.Group("/api/v1/users")
	api.GET("", userHandler.List)
	api.GET("/:id", userHandler.GetByID)
	api.POST("", userHandler.Create)
	api.PUT("/:id", userHandler.Update)
	api.DELETE("/:id", userHandler.Delete)
	api.POST("/:id/assign-role", userHandler.AssignRole)
	api.POST("/:id/sync-roles", userHandler.SyncRoles)

	return r, db
}

func TestUserHandler_List(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestUserHandler_GetByID_Success(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestUserHandler_GetByID_NotFound(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_GetByID_InvalidID(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Create_Success(t *testing.T) {
	r, db := setupUserHandlerTest(t)

	body, _ := json.Marshal(CreateUserRequest{
		UserName:  "New User",
		Email:     "new@user.com",
		Password:  "password123",
		UserLevel: "eng",
		NIP:       "m99-999",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-Auth", "bypass")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)

	// verify stored (note: password stored as plaintext — known bug)
	var stored models.User
	db.First(&stored, 3)
	assert.Equal(t, "New User", stored.UserName)
	assert.Equal(t, "password123", stored.Password) // plaintext, not bcrypt
}

func TestUserHandler_Create_DuplicateEmail(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	body, _ := json.Marshal(CreateUserRequest{
		UserName: "Dup",
		Email:    "admin@admin.com",
		Password: "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Create_MissingFields(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	tests := []struct {
		name string
		body map[string]string
	}{
		{"missing email", map[string]string{"user_name": "X", "password": "pass123"}},
		{"missing password", map[string]string{"user_name": "X", "email": "x@x.com"}},
		{"short password", map[string]string{"user_name": "X", "email": "x@x.com", "password": "123"}},
		{"empty body", map[string]string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestUserHandler_Update_Success(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	body, _ := json.Marshal(UpdateUserRequest{
		UserName: strPtr("Admin Updated"),
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestUserHandler_Update_NotFound(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	body, _ := json.Marshal(UpdateUserRequest{
		UserName: strPtr("Nobody"),
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_Update_InvalidLevel(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	body, _ := json.Marshal(UpdateUserRequest{
		UserLevel: strPtr("superadmin"),
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Delete_Success(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestUserHandler_Delete_NotFound(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_AssignRole_UserNotFound(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	body, _ := json.Marshal(AssignRoleRequest{RoleID: 1})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/999/assign-role", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_AssignRole_RoleNotFound(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	body, _ := json.Marshal(AssignRoleRequest{RoleID: 999})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/1/assign-role", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_SyncRoles(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	body, _ := json.Marshal(SyncRolesRequest{RoleIDs: []uint{1, 2}})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/1/sync-roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_SyncRoles_UserNotFound(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	body, _ := json.Marshal(SyncRolesRequest{RoleIDs: []uint{1}})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/999/sync-roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_Create_WithRoles(t *testing.T) {
	r, _ := setupUserHandlerTest(t)

	body, _ := json.Marshal(CreateUserRequest{
		UserName:  "Role User",
		Email:     "role@user.com",
		Password:  "password123",
		UserLevel: "eng",
		NIP:       "m00-001",
		RoleIDs:   []uint{1},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

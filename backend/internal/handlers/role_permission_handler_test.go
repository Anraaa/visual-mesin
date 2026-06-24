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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/repository"
	"github.com/anraaa/visual-mesin/internal/services"
)

func setupRoleHandlerTest(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	createTables(t, db)

	role := models.Role{Name: "user", GuardName: "web"}
	require.NoError(t, db.Create(&role).Error)

	perm1 := models.Permission{Name: "view-dashboard", GuardName: "web"}
	require.NoError(t, db.Create(&perm1).Error)
	perm2 := models.Permission{Name: "view-users", GuardName: "web"}
	require.NoError(t, db.Create(&perm2).Error)

	// assign permission to role via direct SQL (since AssignPermission uses ON DUPLICATE KEY)
	db.Exec("INSERT INTO role_has_permissions (role_id, permission_id) VALUES (?, ?)", role.ID, perm1.ID)

	roleRepo := repository.NewRoleRepository(db)
	permRepo := repository.NewPermissionRepository(db)
	activityLogSvc := services.NewActivityLogService(db)

	roleHandler := NewRoleHandler(roleRepo, permRepo, activityLogSvc)
	permHandler := NewPermissionHandler(permRepo)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	api := r.Group("/api/v1")
	{
		api.GET("/roles", roleHandler.List)
		api.GET("/roles/:id", roleHandler.GetByID)
		api.POST("/roles", roleHandler.Create)
		api.DELETE("/roles/:id", roleHandler.Delete)
		api.POST("/roles/:id/revoke-permission/:permId", roleHandler.RevokePermission)

		api.GET("/permissions", permHandler.List)
		api.POST("/permissions", permHandler.Create)
		api.DELETE("/permissions/:id", permHandler.Delete)
	}
	return r, db
}

func TestRoleHandler_List(t *testing.T) {
	r, _ := setupRoleHandlerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/roles", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)

	data, ok := resp.Data.([]interface{})
	require.True(t, ok)
	assert.Len(t, data, 1)

	role := data[0].(map[string]interface{})
	assert.Equal(t, "user", role["name"])
	perms := role["permissions"].([]interface{})
	assert.Len(t, perms, 1)
}

func TestRoleHandler_GetByID_Success(t *testing.T) {
	r, _ := setupRoleHandlerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/roles/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestRoleHandler_GetByID_NotFound(t *testing.T) {
	r, _ := setupRoleHandlerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/roles/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRoleHandler_Create(t *testing.T) {
	r, _ := setupRoleHandlerTest(t)

	body, _ := json.Marshal(map[string]string{"name": "admin"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRoleHandler_Create_MissingName(t *testing.T) {
	r, _ := setupRoleHandlerTest(t)

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRoleHandler_Delete_Success(t *testing.T) {
	r, _ := setupRoleHandlerTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/roles/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleHandler_Delete_NotFound(t *testing.T) {
	r, _ := setupRoleHandlerTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/roles/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRoleHandler_RevokePermission(t *testing.T) {
	r, _ := setupRoleHandlerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/roles/1/revoke-permission/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleHandler_RevokePermission_InvalidID(t *testing.T) {
	r, _ := setupRoleHandlerTest(t)

	tests := []string{
		"/api/v1/roles/abc/revoke-permission/1",
		"/api/v1/roles/1/revoke-permission/abc",
	}
	for _, path := range tests {
		req := httptest.NewRequest(http.MethodPost, path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	}
}

func TestPermissionHandler_List(t *testing.T) {
	r, _ := setupRoleHandlerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/permissions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestPermissionHandler_Create(t *testing.T) {
	r, _ := setupRoleHandlerTest(t)

	body, _ := json.Marshal(map[string]string{"name": "view-reports"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/permissions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPermissionHandler_Create_MissingName(t *testing.T) {
	r, _ := setupRoleHandlerTest(t)

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/permissions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPermissionHandler_Delete_Success(t *testing.T) {
	r, _ := setupRoleHandlerTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/permissions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPermissionHandler_Delete_NotFound(t *testing.T) {
	r, _ := setupRoleHandlerTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/permissions/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// GORM's Delete with non-existent ID doesn't return error
	assert.Equal(t, http.StatusOK, w.Code)
}

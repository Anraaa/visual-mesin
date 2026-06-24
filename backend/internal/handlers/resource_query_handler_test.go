package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/anraaa/visual-mesin/internal/db"
	"github.com/anraaa/visual-mesin/internal/repository"
	"github.com/anraaa/visual-mesin/internal/services"
)

func setupResourceQueryHandlerTest(t *testing.T) *gin.Engine {
	t.Helper()
	localDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	localDB.Exec(`CREATE TABLE resource_db_configs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		resource_name VARCHAR(255) NOT NULL,
		alias VARCHAR(255),
		driver VARCHAR(50) NOT NULL DEFAULT 'mysql',
		host VARCHAR(255) NOT NULL DEFAULT '127.0.0.1',
		port INTEGER NOT NULL DEFAULT 3306,
		username VARCHAR(255),
		password TEXT,
		database_name VARCHAR(255),
		is_active BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME,
		updated_at DATETIME
	)`)

	resourceRepo := repository.NewResourceDBConfigRepository(localDB)
	dbMgr := db.NewManager()
	svc := services.NewResourceQueryService(resourceRepo, dbMgr, localDB)
	activityLogSvc := services.NewActivityLogService(localDB)
	h := NewResourceQueryHandler(svc, activityLogSvc)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/resources-judgment/:resource", h.GetJudgmentSummary)
	return r
}

func TestGetJudgmentSummaryHandler_NotFound(t *testing.T) {
	r := setupResourceQueryHandlerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/resources-judgment/nonexistent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetJudgmentSummaryHandler_InactiveResource(t *testing.T) {
	localDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	localDB.Exec(`CREATE TABLE resource_db_configs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		resource_name VARCHAR(255) NOT NULL,
		alias VARCHAR(255),
		driver VARCHAR(50) NOT NULL DEFAULT 'mysql',
		host VARCHAR(255) NOT NULL DEFAULT '127.0.0.1',
		port INTEGER NOT NULL DEFAULT 3306,
		username VARCHAR(255),
		password TEXT,
		database_name VARCHAR(255),
		is_active BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME,
		updated_at DATETIME
	)`)
	localDB.Exec(`INSERT INTO resource_db_configs (resource_name, is_active) VALUES (?, ?)`, "inactive", 0)

	resourceRepo := repository.NewResourceDBConfigRepository(localDB)
	dbMgr := db.NewManager()
	svc := services.NewResourceQueryService(resourceRepo, dbMgr, localDB)
	activityLogSvc := services.NewActivityLogService(localDB)
	h := NewResourceQueryHandler(svc, activityLogSvc)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/resources-judgment/:resource", h.GetJudgmentSummary)

	req := httptest.NewRequest(http.MethodGet, "/resources-judgment/inactive", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

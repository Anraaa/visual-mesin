package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/anraaa/visual-mesin/internal/db"
	"github.com/anraaa/visual-mesin/internal/repository"
)

func setupResourceQueryTest(t *testing.T) (*ResourceQueryService, *gorm.DB) {
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
	svc := NewResourceQueryService(resourceRepo, dbMgr, localDB)
	return svc, localDB
}

func TestGetJudgmentSummary_ResourceNotFound(t *testing.T) {
	svc, _ := setupResourceQueryTest(t)

	result, err := svc.GetJudgmentSummary("nonexistent_table")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource not found")
	assert.Nil(t, result)
}

func TestGetJudgmentSummary_ResourceIsInactive(t *testing.T) {
	svc, _ := setupResourceQueryTest(t)

	svc.localDB.Exec(`INSERT INTO resource_db_configs (resource_name, is_active) VALUES (?, ?)`, "inactive_table", 0)

	result, err := svc.GetJudgmentSummary("inactive_table")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource is not active")
	assert.Nil(t, result)
}

func TestGetJudgmentSummary_EmptyName(t *testing.T) {
	svc, _ := setupResourceQueryTest(t)

	_, err := svc.GetJudgmentSummary("")
	assert.Error(t, err)
}

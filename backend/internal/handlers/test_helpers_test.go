package handlers

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func strPtr(s string) *string { return &s }

func createTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nip VARCHAR(50), user_id VARCHAR(100),
		user_name VARCHAR(100) NOT NULL,
		user_level VARCHAR(50) NOT NULL DEFAULT 'prod',
		email VARCHAR(255), email_verified_at DATETIME,
		password VARCHAR(255) NOT NULL,
		avatar_url VARCHAR(255), remember_token VARCHAR(100),
		department VARCHAR(100), jabatan VARCHAR(100),
		themes_settings TEXT, timestamp DATETIME,
		created_at DATETIME, updated_at DATETIME
	)`).Error)

	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS roles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		guard_name VARCHAR(255) NOT NULL DEFAULT 'web',
		created_at DATETIME, updated_at DATETIME
	)`).Error)

	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS permissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		guard_name VARCHAR(255) NOT NULL DEFAULT 'web',
		created_at DATETIME, updated_at DATETIME
	)`).Error)

	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS activity_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		log_name TEXT, description TEXT, subject_type TEXT,
		subject_id INTEGER, causer_type TEXT, causer_id INTEGER,
		properties TEXT, event TEXT, created_at DATETIME
	)`).Error)

	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS model_has_roles (
		role_id INTEGER, model_type TEXT, model_id INTEGER
	)`).Error)

	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS model_has_permissions (
		permission_id INTEGER, model_type TEXT, model_id INTEGER
	)`).Error)

	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS role_has_permissions (
		role_id INTEGER, permission_id INTEGER
	)`).Error)
}

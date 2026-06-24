package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/anraaa/visual-mesin/internal/models"
)

func setupUserRepoTest(t *testing.T) (*UserRepository, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nip VARCHAR(50),
		user_id VARCHAR(100),
		user_name VARCHAR(100) NOT NULL,
		user_level VARCHAR(50) NOT NULL DEFAULT 'prod',
		email VARCHAR(255),
		email_verified_at DATETIME,
		password VARCHAR(255) NOT NULL,
		avatar_url VARCHAR(255),
		remember_token VARCHAR(100),
		department VARCHAR(100),
		jabatan VARCHAR(100),
		themes_settings TEXT,
		timestamp DATETIME,
		created_at DATETIME,
		updated_at DATETIME
	)`)
	db.Exec(`CREATE TABLE roles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		guard_name VARCHAR(255) NOT NULL DEFAULT 'web',
		created_at DATETIME,
		updated_at DATETIME
	)`)
	db.Exec(`CREATE TABLE permissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		guard_name VARCHAR(255) NOT NULL DEFAULT 'web',
		created_at DATETIME,
		updated_at DATETIME
	)`)
	db.Exec(`CREATE TABLE model_has_roles (
		role_id INTEGER,
		model_type TEXT,
		model_id INTEGER
	)`)
	db.Exec(`CREATE TABLE model_has_permissions (
		permission_id INTEGER,
		model_type TEXT,
		model_id INTEGER
	)`)
	db.Exec(`CREATE TABLE role_has_permissions (
		role_id INTEGER,
		permission_id INTEGER
	)`)

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

	return NewUserRepository(db), db
}

func strPtr(s string) *string { return &s }

func TestFindByEmailOrNIP_ByEmail(t *testing.T) {
	repo, _ := setupUserRepoTest(t)

	u, err := repo.FindByEmailOrNIP("admin@admin.com")
	require.NoError(t, err)
	assert.Equal(t, "Admin", u.UserName)
}

func TestFindByEmailOrNIP_ByFullNIP(t *testing.T) {
	repo, _ := setupUserRepoTest(t)

	u, err := repo.FindByEmailOrNIP("m26-134")
	require.NoError(t, err)
	assert.Equal(t, "Admin", u.UserName)
}

func TestFindByEmailOrNIP_ByNIPWithoutMPrefix(t *testing.T) {
	repo, _ := setupUserRepoTest(t)

	u, err := repo.FindByEmailOrNIP("26-134")
	require.NoError(t, err)
	assert.Equal(t, "Admin", u.UserName)
}

func TestFindByEmailOrNIP_ByNIPWithCapitalM(t *testing.T) {
	repo, _ := setupUserRepoTest(t)

	u, err := repo.FindByEmailOrNIP("M26-134")
	require.NoError(t, err)
	assert.Equal(t, "Admin", u.UserName)
}

func TestFindByEmailOrNIP_ForUserWithoutMPrefix(t *testing.T) {
	repo, _ := setupUserRepoTest(t)

	u, err := repo.FindByEmailOrNIP("26-133")
	require.NoError(t, err)
	assert.Equal(t, "User Produksi", u.UserName)
}

func TestFindByEmailOrNIP_ForUserWithMPrefix(t *testing.T) {
	repo, _ := setupUserRepoTest(t)

	u, err := repo.FindByEmailOrNIP("m26-133")
	require.NoError(t, err)
	assert.Equal(t, "User Produksi", u.UserName)
}

func TestFindByEmailOrNIP_NotFound(t *testing.T) {
	repo, _ := setupUserRepoTest(t)

	_, err := repo.FindByEmailOrNIP("nonexistent@email.com")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestFindByID_Success(t *testing.T) {
	repo, _ := setupUserRepoTest(t)

	u, err := repo.FindByID(1)
	require.NoError(t, err)
	assert.Equal(t, "Admin", u.UserName)

	u, err = repo.FindByID(2)
	require.NoError(t, err)
	assert.Equal(t, "User Produksi", u.UserName)
}

func TestFindByID_NotFound(t *testing.T) {
	repo, _ := setupUserRepoTest(t)

	_, err := repo.FindByID(999)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCreate(t *testing.T) {
	repo, _ := setupUserRepoTest(t)

	hash, _ := bcrypt.GenerateFromPassword([]byte("newpass"), bcrypt.DefaultCost)
	newUser := models.User{
		Email:     strPtr("new@user.com"),
		UserName:  "New User",
		NIP:       strPtr("m99-999"),
		Password:  string(hash),
		UserLevel: "eng",
	}
	err := repo.Create(&newUser)
	require.NoError(t, err)
	assert.NotZero(t, newUser.ID)

	// verify FindByEmail works
	u, err := repo.FindByEmailOrNIP("new@user.com")
	require.NoError(t, err)
	assert.Equal(t, "New User", u.UserName)
}

func TestUpdate(t *testing.T) {
	repo, _ := setupUserRepoTest(t)

	u, err := repo.FindByID(1)
	require.NoError(t, err)

	u.UserName = "Admin Updated"
	err = repo.Update(u)
	require.NoError(t, err)

	u2, err := repo.FindByID(1)
	require.NoError(t, err)
	assert.Equal(t, "Admin Updated", u2.UserName)
}

func TestDelete(t *testing.T) {
	repo, _ := setupUserRepoTest(t)

	err := repo.Delete(1)
	require.NoError(t, err)

	_, err = repo.FindByID(1)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestList(t *testing.T) {
	repo, _ := setupUserRepoTest(t)

	users, total, err := repo.List(1, 10)
	require.NoError(t, err)
	assert.EqualValues(t, 2, total)
	assert.Len(t, users, 2)
}

func TestListPagination(t *testing.T) {
	repo, db := setupUserRepoTest(t)

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	for i := 0; i < 5; i++ {
		email := "extra@user.com"
		nip := "m00-001"
		u := models.User{
			Email:     &email,
			UserName:  "Extra",
			NIP:       &nip,
			Password:  string(hash),
			UserLevel: "eng",
		}
		db.Create(&u)
	}

	users, total, err := repo.List(1, 3)
	require.NoError(t, err)
	assert.EqualValues(t, 7, total)
	assert.Len(t, users, 3)
}

func TestGetRolesAndPermissions(t *testing.T) {
	repo, db := setupUserRepoTest(t)

	// seed role + permission
	role := models.Role{Name: "admin", GuardName: "web"}
	db.Create(&role)
	perm := models.Permission{Name: "view-dashboard", GuardName: "web"}
	db.Create(&perm)

	db.Exec("INSERT INTO role_has_permissions (role_id, permission_id) VALUES (?, ?)", role.ID, perm.ID)
	db.Exec("INSERT INTO model_has_roles (role_id, model_type, model_id) VALUES (?, ?, ?)", role.ID, "App\\Models\\User", 1)

	roles, err := repo.GetRoles(1)
	require.NoError(t, err)
	assert.Len(t, roles, 1)
	assert.Equal(t, "admin", roles[0].Name)

	permissions, err := repo.GetPermissions(1)
	require.NoError(t, err)
	assert.Len(t, permissions, 1)
	assert.Equal(t, "view-dashboard", permissions[0].Name)
}

func TestGetPermissionsViaDirectAssignment(t *testing.T) {
	repo, db := setupUserRepoTest(t)

	perm := models.Permission{Name: "view-report", GuardName: "web"}
	db.Create(&perm)
	db.Exec("INSERT INTO model_has_permissions (permission_id, model_type, model_id) VALUES (?, ?, ?)", perm.ID, "App\\Models\\User", 1)

	permissions, err := repo.GetPermissions(1)
	require.NoError(t, err)
	assert.Len(t, permissions, 1)
	assert.Equal(t, "view-report", permissions[0].Name)
}

func TestRevokeRoles(t *testing.T) {
	repo, db := setupUserRepoTest(t)

	role := models.Role{Name: "user", GuardName: "web"}
	db.Create(&role)

	db.Exec("INSERT INTO model_has_roles (role_id, model_type, model_id) VALUES (?, ?, ?)", role.ID, "App\\Models\\User", 1)

	roles, err := repo.GetRoles(1)
	require.NoError(t, err)
	assert.Len(t, roles, 1)

	err = repo.RevokeAllRoles(1)
	require.NoError(t, err)

	roles, err = repo.GetRoles(1)
	require.NoError(t, err)
	assert.Len(t, roles, 0)
}

func TestFindByEmailOrNIP_NIPCaseInsensitive(t *testing.T) {
	repo, _ := setupUserRepoTest(t)

	u, err := repo.FindByEmailOrNIP("M26-134")
	require.NoError(t, err)
	assert.Equal(t, "Admin", u.UserName)

	u, err = repo.FindByEmailOrNIP("26-133")
	require.NoError(t, err)
	assert.Equal(t, "User Produksi", u.UserName)
}

func TestFindByEmailOrNIP_EmptyIdentifier(t *testing.T) {
	repo, _ := setupUserRepoTest(t)

	_, err := repo.FindByEmailOrNIP("")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

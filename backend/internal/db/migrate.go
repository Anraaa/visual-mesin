package db

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var allPermissions = []string{
	"view-dashboard",
	"view-any-dashboard",
	"view-data-produksi",
	"view-any-data-produksi",
	"view-export",
	"view-any-export",
	"create-export",
	"delete-export",
	"view-ai-chat",
	"view-any-ai-chat",
	"view-user",
	"view-any-user",
	"create-user",
	"update-user",
	"delete-user",
	"delete-any-user",
	"view-role",
	"view-any-role",
	"create-role",
	"update-role",
	"delete-role",
	"delete-any-role",
	"view-permission",
	"view-any-permission",
	"create-permission",
	"update-permission",
	"delete-permission",
	"delete-any-permission",
	"view-resource-connection",
	"view-any-resource-connection",
	"create-resource-connection",
	"update-resource-connection",
	"delete-resource-connection",
	"delete-any-resource-connection",
	"view-data-produksi-config",
	"view-any-data-produksi-config",
	"view-activity-log",
	"view-any-activity-log",
}

func RunMigrations(dsn, migrationsPath string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("mysql://%s", dsn),
	)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}

func ensureRole(db *gorm.DB, name string) uint {
	var role struct{ ID uint }
	var count int64
	db.Table("roles").Where("name = ?", name).Count(&count)
	if count == 0 {
		db.Table("roles").Create(map[string]interface{}{
			"name":       name,
			"guard_name": "web",
		})
	}
	db.Table("roles").Where("name = ?", name).First(&role)
	return role.ID
}

func SeedRolesAndPermissions(db *gorm.DB) {
	for _, name := range allPermissions {
		var count int64
		db.Model(&struct{}{}).Table("permissions").Where("name = ?", name).Count(&count)
		if count == 0 {
			if err := db.Table("permissions").Create(map[string]interface{}{
				"name":       name,
				"guard_name": "web",
			}).Error; err != nil {
				log.Printf("Seed: failed to create permission %s: %v", name, err)
			}
		}
	}

	adminRoleID := ensureRole(db, "admin")
	userRoleID := ensureRole(db, "user")

	var permCount int64
	db.Model(&struct{}{}).Table("role_has_permissions").Where("role_id = ?", adminRoleID).Count(&permCount)
	if permCount == 0 {
		var perms []struct{ ID uint }
		db.Table("permissions").Find(&perms)
		for _, p := range perms {
			db.Exec("INSERT IGNORE INTO role_has_permissions (role_id, permission_id) VALUES (?, ?)", adminRoleID, p.ID)
		}
		log.Printf("Seed: assigned all permissions to admin role")
	} else {
		var unassigned []struct{ ID uint }
		db.Raw(`SELECT p.id FROM permissions p
			WHERE p.id NOT IN (SELECT rhp.permission_id FROM role_has_permissions rhp WHERE rhp.role_id = ?)`, adminRoleID).
			Scan(&unassigned)
		for _, p := range unassigned {
			db.Exec("INSERT IGNORE INTO role_has_permissions (role_id, permission_id) VALUES (?, ?)", adminRoleID, p.ID)
		}
		if len(unassigned) > 0 {
			log.Printf("Seed: assigned %d new permissions to admin role", len(unassigned))
		}
	}

	assignRoleToUserByEmail(db, adminRoleID, "admin@admin.com")
	assignRoleToUserByEmail(db, userRoleID, "user@visualmesin.com")

	fixRoleAssignments(db)
}

func assignRoleToUserByEmail(db *gorm.DB, roleID uint, email string) {
	var user struct{ ID uint }
	result := db.Table("users").Where("email = ?", email).First(&user)
	if result.Error != nil {
		return
	}
	var exists int64
	db.Model(&struct{}{}).Table("model_has_roles").
		Where("role_id = ? AND model_id = ? AND model_type = ?", roleID, user.ID, "App\\Models\\User").
		Count(&exists)
	if exists == 0 {
		db.Exec("INSERT IGNORE INTO model_has_roles (role_id, model_type, model_id) VALUES (?, ?, ?)",
			roleID, "App\\Models\\User", user.ID)
		log.Printf("Seed: assigned role_id=%d to user %s", roleID, email)
	}
}

func fixRoleAssignments(db *gorm.DB) {
	type userInfo struct {
		ID        uint
		Email     string
		UserLevel string
	}
	var adminRole struct{ ID uint }
	if err := db.Table("roles").Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return
	}

	var users []userInfo
	db.Table("users").Find(&users)
	for _, u := range users {
		if u.UserLevel == "admin" {
			db.Exec("DELETE FROM model_has_roles WHERE model_type = ? AND model_id = ? AND role_id != ?",
				"App\\Models\\User", u.ID, adminRole.ID)
			var exists int64
			db.Model(&struct{}{}).Table("model_has_roles").
				Where("role_id = ? AND model_id = ? AND model_type = ?", adminRole.ID, u.ID, "App\\Models\\User").
				Count(&exists)
			if exists == 0 {
				db.Exec("INSERT IGNORE INTO model_has_roles (role_id, model_type, model_id) VALUES (?, ?, ?)",
					adminRole.ID, "App\\Models\\User", u.ID)
				log.Printf("Fix: assigned admin role to admin user %s (id=%d)", u.Email, u.ID)
			}
		} else {
			db.Exec("DELETE FROM model_has_roles WHERE model_type = ? AND model_id = ? AND role_id = ?",
				"App\\Models\\User", u.ID, adminRole.ID)
		}
	}
}

func SeedDefaultUsers(db *gorm.DB) {
	type User struct {
		Email     string
		Password  string
		UserName  string
		UserLevel string
		NIP       string
		UserID    string
	}

	defaultUsers := []User{
		{"admin@admin.com", "password", "Admin", "admin", "m26-134", "admin"},
		{"user@visualmesin.com", "user123", "User Produksi", "prod", "PRD001", "user"},
	}

	for _, u := range defaultUsers {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Seed: failed to hash password for %s: %v", u.Email, err)
			continue
		}
		passwordStr := string(hashed)

		var count int64
		db.Model(&struct{}{}).Table("users").Where("email = ?", u.Email).Count(&count)
		if count > 0 {
			if err := db.Table("users").Where("email = ?", u.Email).Updates(map[string]interface{}{
				"password": passwordStr,
				"nip":      u.NIP,
			}).Error; err != nil {
				log.Printf("Seed: failed to update user %s: %v", u.Email, err)
			} else {
				log.Printf("Seed: updated user %s", u.Email)
			}
		} else {
			user := map[string]interface{}{
				"user_name":  u.UserName,
				"user_level": u.UserLevel,
				"email":      u.Email,
				"password":   passwordStr,
				"nip":        u.NIP,
				"user_id":    u.UserID,
			}
			if err := db.Table("users").Create(user).Error; err != nil {
				log.Printf("Seed: failed to create %s: %v", u.Email, err)
			} else {
				log.Printf("Seed: created user %s", u.Email)
			}
		}
	}
}

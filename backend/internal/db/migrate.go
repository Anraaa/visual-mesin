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

func SeedDefaultUsers(db *gorm.DB) {
	type User struct {
		Email    string
		Password string
		UserName string
		UserLevel string
		NIP      string
		UserID   string
	}

	defaultUsers := []User{
		{"admin@visualmesin.com", "admin12", "Admin", "admin", "ADM001", "admin"},
		{"user@visualmesin.com", "user123", "User Produksi", "prod", "PRD001", "user"},
	}

	for _, u := range defaultUsers {
		var count int64
		db.Model(&struct{}{}).Table("users").Where("email = ?", u.Email).Count(&count)
		if count > 0 {
			hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Seed: failed to hash password for %s: %v", u.Email, err)
				continue
			}
			if err := db.Table("users").Where("email = ?", u.Email).Update("password", string(hash)).Error; err != nil {
				log.Printf("Seed: failed to update password for %s: %v", u.Email, err)
			} else {
				log.Printf("Seed: updated password for %s", u.Email)
			}
		} else {
			hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Seed: failed to hash password for %s: %v", u.Email, err)
				continue
			}
			user := map[string]interface{}{
				"user_name":  u.UserName,
				"user_level": u.UserLevel,
				"email":      u.Email,
				"password":   string(hash),
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

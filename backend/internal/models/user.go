package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	NIP             *string    `gorm:"type:varchar(50);uniqueIndex" json:"nip"`
	UserID          *string    `gorm:"type:varchar(100);uniqueIndex" json:"user_id"`
	UserName        string     `gorm:"type:varchar(100);not null" json:"user_name"`
	UserLevel       string     `gorm:"type:enum('admin','eng','tech','prod');not null;default:'prod'" json:"user_level"`
	Email           *string    `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	Password        string     `gorm:"type:varchar(255);not null" json:"-"`
	AvatarURL       *string    `gorm:"type:varchar(255)" json:"avatar_url"`
	RememberToken   *string    `gorm:"type:varchar(100)" json:"-"`
	Department      *string    `gorm:"type:varchar(100)" json:"department"`
	Jabatan         *string    `gorm:"type:varchar(100)" json:"jabatan"`
	ThemesSettings  *string    `gorm:"type:json" json:"themes_settings"`
	Timestamp       *time.Time `gorm:"type:timestamp(3)" json:"timestamp"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Timestamp == nil {
		now := time.Now()
		u.Timestamp = &now
	}
	return nil
}

func (User) TableName() string {
	return "users"
}

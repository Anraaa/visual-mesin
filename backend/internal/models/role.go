package models

import "time"

type Role struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_role_name_guard" json:"name"`
	GuardName string    `gorm:"type:varchar(255);not null;default:'web';uniqueIndex:idx_role_name_guard" json:"guard_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Permissions []Permission `gorm:"many2many:role_has_permissions;" json:"permissions,omitempty"`
}

type Permission struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_perm_name_guard" json:"name"`
	GuardName string    `gorm:"type:varchar(255);not null;default:'web';uniqueIndex:idx_perm_name_guard" json:"guard_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Roles []Role `gorm:"many2many:role_has_permissions;" json:"roles,omitempty"`
}

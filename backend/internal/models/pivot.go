package models

type ModelHasRole struct {
	RoleID    uint   `gorm:"primaryKey;autoIncrement:false"`
	ModelType string `gorm:"type:varchar(255);primaryKey;autoIncrement:false"`
	ModelID   uint   `gorm:"primaryKey;autoIncrement:false"`
}

func (ModelHasRole) TableName() string {
	return "model_has_roles"
}

type ModelHasPermission struct {
	PermissionID uint   `gorm:"primaryKey;autoIncrement:false"`
	ModelType    string `gorm:"type:varchar(255);primaryKey;autoIncrement:false"`
	ModelID      uint   `gorm:"primaryKey;autoIncrement:false"`
}

func (ModelHasPermission) TableName() string {
	return "model_has_permissions"
}

type RoleHasPermission struct {
	PermissionID uint `gorm:"primaryKey;autoIncrement:false"`
	RoleID       uint `gorm:"primaryKey;autoIncrement:false"`
}

func (RoleHasPermission) TableName() string {
	return "role_has_permissions"
}

package models

import "time"

type ResourceGroup struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Color     string    `gorm:"type:varchar(7);not null;default:'#1677ff'" json:"color"`
	Icon      string    `gorm:"type:varchar(100);not null;default:'BuildOutlined'" json:"icon"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	Items     []ResourceGroupItem `gorm:"foreignKey:GroupID" json:"items,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ResourceGroup) TableName() string {
	return "resource_groups"
}

type ResourceGroupItem struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID      uint      `gorm:"not null;index" json:"group_id"`
	ResourceName string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"resource_name"`
	Label        string    `gorm:"type:varchar(255)" json:"label"`
	SortOrder    int       `gorm:"not null;default:0" json:"sort_order"`
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (ResourceGroupItem) TableName() string {
	return "resource_group_items"
}

type ResourceColumnDef struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ResourceName   string    `gorm:"type:varchar(255);not null;index" json:"resource_name"`
	ColumnName     string    `gorm:"type:varchar(255);not null" json:"column_name"`
	DataType       string    `gorm:"type:varchar(50);not null" json:"data_type"`
	Length         int       `gorm:"type:int" json:"length,omitempty"`
	DecimalPlaces  int       `gorm:"type:int" json:"decimal_places,omitempty"`
	EnumValues     string    `gorm:"type:text" json:"enum_values,omitempty"`
	IsNullable     bool      `gorm:"not null;default:true" json:"is_nullable"`
	DefaultValue   string    `gorm:"type:varchar(255)" json:"default_value,omitempty"`
	IsPrimary      bool      `gorm:"not null;default:false" json:"is_primary"`
	IsAutoIncrement bool    `gorm:"not null;default:false" json:"is_auto_increment"`
	SortOrder      int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (ResourceColumnDef) TableName() string {
	return "resource_column_defs"
}

type CreateResourceRequest struct {
	GroupID      uint                `json:"group_id" binding:"required"`
	ResourceName string              `json:"resource_name" binding:"required"`
	Label        string              `json:"label"`
	Columns      []CreateColumnInput `json:"columns" binding:"required"`
}

type CreateColumnInput struct {
	ColumnName     string `json:"column_name" binding:"required"`
	DataType       string `json:"data_type" binding:"required"`
	Length         *int   `json:"length"`
	DecimalPlaces  *int   `json:"decimal_places"`
	EnumValues     string `json:"enum_values"`
	IsNullable     bool   `json:"is_nullable"`
	DefaultValue   string `json:"default_value"`
	IsPrimary      bool   `json:"is_primary"`
	IsAutoIncrement bool  `json:"is_auto_increment"`
	SortOrder      int    `json:"sort_order"`
}

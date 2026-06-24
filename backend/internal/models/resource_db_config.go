package models

import (
	"time"
)

type ResourceDBConfig struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ResourceName      string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"resource_name"`
	Label             *string    `gorm:"type:varchar(255)" json:"label,omitempty"`
	Driver            string     `gorm:"type:varchar(20);not null;default:'mariadb'" json:"driver"`
	Host              string     `gorm:"type:varchar(255);not null" json:"host"`
	Port              int        `gorm:"type:int;not null;default:3306" json:"port"`
	DatabaseName      string     `gorm:"type:varchar(255);not null" json:"database_name"`
	Username          string     `gorm:"type:varchar(255);not null" json:"username"`
	Password          string     `gorm:"type:text;not null" json:"-"`
	IsActive          bool       `gorm:"not null;default:true;index:idx_is_active" json:"is_active"`
	IsLastTestSuccess *bool      `gorm:"type:boolean" json:"is_last_test_success,omitempty"`
	LastTestedAt      *time.Time `gorm:"type:timestamp" json:"last_tested_at,omitempty"`
	LastTestMessage   *string    `gorm:"type:text" json:"last_test_message,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func (ResourceDBConfig) TableName() string {
	return "resource_db_configs"
}

type ResourceDBConfigRequest struct {
	ResourceName string `json:"resource_name" binding:"required"`
	Label        string `json:"label"`
	Driver       string `json:"driver" binding:"required,oneof=mysql mariadb postgresql sqlite"`
	Host         string `json:"host" binding:"required"`
	Port         int    `json:"port" binding:"required"`
	DatabaseName string `json:"database_name" binding:"required"`
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	IsActive     *bool  `json:"is_active"`
}

type ResourceDBConfigUpdateRequest struct {
	ResourceName *string `json:"resource_name"`
	Label        *string `json:"label"`
	Driver       *string `json:"driver" binding:"omitempty,oneof=mysql mariadb postgresql sqlite"`
	Host         *string `json:"host"`
	Port         *int    `json:"port"`
	DatabaseName *string `json:"database_name"`
	Username     *string `json:"username"`
	Password     *string `json:"password"`
	IsActive     *bool   `json:"is_active"`
}

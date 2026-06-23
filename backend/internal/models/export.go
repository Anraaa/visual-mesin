package models

import (
	"time"

	"gorm.io/gorm"
)

type ExportJob struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ResourceName  string         `gorm:"type:varchar(255);not null" json:"resource_name"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	Status        string         `gorm:"type:enum('queued','processing','completed','failed');not null;default:'queued';index" json:"status"`
	TotalRows     int64          `gorm:"not null;default:0" json:"total_rows"`
	ProcessedRows int64          `gorm:"not null;default:0" json:"processed_rows"`
	FilePath      *string        `gorm:"type:varchar(255)" json:"file_path"`
	FileSize      *int64         `json:"file_size"`
	Format        string         `gorm:"type:enum('csv','xlsx');not null;default:'csv'" json:"format"`
	Columns       *string        `gorm:"type:json" json:"columns"`
	Filters       *string        `gorm:"type:json" json:"filters"`
	ErrorMessage  *string        `gorm:"type:text" json:"error_message"`
	StartedAt     *time.Time     `json:"started_at"`
	CompletedAt   *time.Time     `json:"completed_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ExportJob) TableName() string {
	return "export_jobs"
}

type ExportJobRequest struct {
	ResourceName string   `json:"resource_name" binding:"required"`
	Format       string   `json:"format" binding:"omitempty,oneof=csv xlsx"`
	Columns      []string `json:"columns"`
	Search       string   `json:"search"`
	SearchBy     string   `json:"search_by"`
	Filters      map[string]string `json:"filters"`
}

type ExportJobStatusResponse struct {
	ID            uint       `json:"id"`
	ResourceName  string     `json:"resource_name"`
	Status        string     `json:"status"`
	TotalRows     int64      `json:"total_rows"`
	ProcessedRows int64      `json:"processed_rows"`
	FileURL       *string    `json:"file_url"`
	FileSize      *int64     `json:"file_size"`
	Format        string     `json:"format"`
	ErrorMessage  *string    `json:"error_message"`
	CreatedAt     time.Time  `json:"created_at"`
	CompletedAt   *time.Time `json:"completed_at"`
}

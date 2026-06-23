package services

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type ActivityLogService struct {
	db *gorm.DB
}

func NewActivityLogService(db *gorm.DB) *ActivityLogService {
	return &ActivityLogService{db: db}
}

func (s *ActivityLogService) Log(userID uint, logName, description, event string, properties map[string]interface{}) {
	propsJSON, _ := json.Marshal(properties)

	s.db.Exec(`
		INSERT INTO activity_log (log_name, description, subject_type, subject_id, causer_type, causer_id, properties, event, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, logName, description, "App\\Models\\User", userID, "App\\Models\\User", userID, string(propsJSON), event, time.Now())
}

func (s *ActivityLogService) List(page, limit int) ([]map[string]interface{}, int64, error) {
	var total int64
	s.db.Table("activity_log").Count(&total)

	var logs []map[string]interface{}
	offset := (page - 1) * limit
	err := s.db.Table("activity_log").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, total, err
}

func (s *ActivityLogService) ListByUser(userID uint, page, limit int) ([]map[string]interface{}, int64, error) {
	var total int64
	s.db.Table("activity_log").Where("causer_id = ?", userID).Count(&total)

	var logs []map[string]interface{}
	offset := (page - 1) * limit
	err := s.db.Table("activity_log").
		Where("causer_id = ?", userID).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, total, err
}

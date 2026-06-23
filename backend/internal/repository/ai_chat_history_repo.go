package repository

import (
	"github.com/anraaa/visual-mesin/internal/models"
	"gorm.io/gorm"
)

type AiChatHistoryRepository struct {
	db *gorm.DB
}

func NewAiChatHistoryRepository(db *gorm.DB) *AiChatHistoryRepository {
	return &AiChatHistoryRepository{db: db}
}

func (r *AiChatHistoryRepository) FindByID(id uint) (*models.AiChatHistory, error) {
	var h models.AiChatHistory
	err := r.db.Preload("User").First(&h, id).Error
	return &h, err
}

func (r *AiChatHistoryRepository) Create(h *models.AiChatHistory) error {
	return r.db.Create(h).Error
}

func (r *AiChatHistoryRepository) Update(h *models.AiChatHistory) error {
	return r.db.Save(h).Error
}

func (r *AiChatHistoryRepository) ListBySession(sessionID string, page, limit int) ([]models.AiChatHistory, int64, error) {
	var items []models.AiChatHistory
	var total int64
	r.db.Model(&models.AiChatHistory{}).Where("session_id = ?", sessionID).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Where("session_id = ?", sessionID).Order("created_at ASC").Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *AiChatHistoryRepository) ListSessionsByUser(userID uint) ([]models.AiChatHistory, error) {
	var items []models.AiChatHistory
	err := r.db.Select("session_id, question, created_at").
		Where("user_id = ?", userID).
		Where("id IN (SELECT MIN(id) FROM ai_chat_history GROUP BY session_id)").
		Order("created_at DESC").
		Limit(50).
		Find(&items).Error
	return items, err
}

func (r *AiChatHistoryRepository) DeleteSession(sessionID string) error {
	return r.db.Where("session_id = ?", sessionID).Delete(&models.AiChatHistory{}).Error
}

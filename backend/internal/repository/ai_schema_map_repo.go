package repository

import (
	"github.com/anraaa/visual-mesin/internal/models"
	"gorm.io/gorm"
)

type AiSchemaMapRepository struct {
	db *gorm.DB
}

func NewAiSchemaMapRepository(db *gorm.DB) *AiSchemaMapRepository {
	return &AiSchemaMapRepository{db: db}
}

func (r *AiSchemaMapRepository) FindByID(id uint) (*models.AiSchemaMap, error) {
	var m models.AiSchemaMap
	err := r.db.First(&m, id).Error
	return &m, err
}

func (r *AiSchemaMapRepository) FindByIntentName(name string) (*models.AiSchemaMap, error) {
	var m models.AiSchemaMap
	err := r.db.Where("intent_name = ?", name).First(&m).Error
	return &m, err
}

func (r *AiSchemaMapRepository) Create(m *models.AiSchemaMap) error {
	return r.db.Create(m).Error
}

func (r *AiSchemaMapRepository) Update(m *models.AiSchemaMap) error {
	return r.db.Save(m).Error
}

func (r *AiSchemaMapRepository) Delete(id uint) error {
	return r.db.Delete(&models.AiSchemaMap{}, id).Error
}

func (r *AiSchemaMapRepository) List(page, limit int) ([]models.AiSchemaMap, int64, error) {
	var items []models.AiSchemaMap
	var total int64
	r.db.Model(&models.AiSchemaMap{}).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&items).Error
	return items, total, err
}

func (r *AiSchemaMapRepository) ListActive() ([]models.AiSchemaMap, error) {
	var items []models.AiSchemaMap
	err := r.db.Where("is_active = ?", true).Find(&items).Error
	return items, err
}

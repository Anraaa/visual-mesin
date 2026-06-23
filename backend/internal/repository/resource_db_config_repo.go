package repository

import (
	"github.com/anraaa/visual-mesin/internal/models"
	"gorm.io/gorm"
)

type ResourceDBConfigRepository struct {
	db *gorm.DB
}

func NewResourceDBConfigRepository(db *gorm.DB) *ResourceDBConfigRepository {
	return &ResourceDBConfigRepository{db: db}
}

func (r *ResourceDBConfigRepository) FindByID(id uint) (*models.ResourceDBConfig, error) {
	var cfg models.ResourceDBConfig
	err := r.db.First(&cfg, id).Error
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *ResourceDBConfigRepository) FindByResourceName(name string) (*models.ResourceDBConfig, error) {
	var cfg models.ResourceDBConfig
	err := r.db.Where("resource_name = ?", name).First(&cfg).Error
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *ResourceDBConfigRepository) Create(cfg *models.ResourceDBConfig) error {
	return r.db.Create(cfg).Error
}

func (r *ResourceDBConfigRepository) Update(cfg *models.ResourceDBConfig) error {
	return r.db.Save(cfg).Error
}

func (r *ResourceDBConfigRepository) Delete(id uint) error {
	return r.db.Delete(&models.ResourceDBConfig{}, id).Error
}

func (r *ResourceDBConfigRepository) List(page, limit int) ([]models.ResourceDBConfig, int64, error) {
	var configs []models.ResourceDBConfig
	var total int64

	r.db.Model(&models.ResourceDBConfig{}).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&configs).Error
	if err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}

func (r *ResourceDBConfigRepository) ListActive() ([]models.ResourceDBConfig, error) {
	var configs []models.ResourceDBConfig
	err := r.db.Where("is_active = ?", true).Find(&configs).Error
	return configs, err
}

package repository

import (
	"github.com/anraaa/visual-mesin/internal/models"
	"gorm.io/gorm"
)

type ExportJobRepository struct {
	db *gorm.DB
}

func NewExportJobRepository(db *gorm.DB) *ExportJobRepository {
	return &ExportJobRepository{db: db}
}

func (r *ExportJobRepository) Create(job *models.ExportJob) error {
	return r.db.Create(job).Error
}

func (r *ExportJobRepository) FindByID(id uint) (*models.ExportJob, error) {
	var job models.ExportJob
	err := r.db.First(&job, id).Error
	return &job, err
}

func (r *ExportJobRepository) Update(job *models.ExportJob) error {
	return r.db.Save(job).Error
}

func (r *ExportJobRepository) ListByUser(userID uint, page, limit int) ([]models.ExportJob, int64, error) {
	var jobs []models.ExportJob
	var total int64

	r.db.Model(&models.ExportJob{}).Where("user_id = ?", userID).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Where("user_id = ?", userID).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&jobs).Error

	return jobs, total, err
}

func (r *ExportJobRepository) ListAll(page, limit int) ([]models.ExportJob, int64, error) {
	var jobs []models.ExportJob
	var total int64

	r.db.Model(&models.ExportJob{}).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&jobs).Error

	return jobs, total, err
}

package repository

import (
	"github.com/anraaa/visual-mesin/internal/models"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) List() ([]models.Permission, error) {
	var perms []models.Permission
	err := r.db.Order("name ASC").Find(&perms).Error
	return perms, err
}

func (r *PermissionRepository) FindByID(id uint) (*models.Permission, error) {
	var perm models.Permission
	err := r.db.First(&perm, id).Error
	return &perm, err
}

func (r *PermissionRepository) FindByName(name string) (*models.Permission, error) {
	var perm models.Permission
	err := r.db.Where("name = ?", name).First(&perm).Error
	return &perm, err
}

func (r *PermissionRepository) Create(perm *models.Permission) error {
	return r.db.Create(perm).Error
}

func (r *PermissionRepository) Delete(perm *models.Permission) error {
	return r.db.Delete(perm).Error
}

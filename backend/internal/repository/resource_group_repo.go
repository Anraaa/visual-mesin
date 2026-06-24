package repository

import (
	"github.com/anraaa/visual-mesin/internal/models"
	"gorm.io/gorm"
)

type ResourceGroupRepo struct {
	db *gorm.DB
}

func NewResourceGroupRepo(db *gorm.DB) *ResourceGroupRepo {
	return &ResourceGroupRepo{db: db}
}

func (r *ResourceGroupRepo) List() ([]models.ResourceGroup, error) {
	var groups []models.ResourceGroup
	err := r.db.Order("sort_order ASC").Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Find(&groups).Error
	return groups, err
}

func (r *ResourceGroupRepo) FindByID(id uint) (*models.ResourceGroup, error) {
	var group models.ResourceGroup
	err := r.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).First(&group, id).Error
	return &group, err
}

func (r *ResourceGroupRepo) Create(group *models.ResourceGroup) error {
	return r.db.Create(group).Error
}

func (r *ResourceGroupRepo) Update(group *models.ResourceGroup) error {
	return r.db.Save(group).Error
}

func (r *ResourceGroupRepo) Delete(id uint) error {
	return r.db.Delete(&models.ResourceGroup{}, id).Error
}

// Items

func (r *ResourceGroupRepo) FindItemByResourceName(name string) (*models.ResourceGroupItem, error) {
	var item models.ResourceGroupItem
	err := r.db.Where("resource_name = ?", name).First(&item).Error
	return &item, err
}

func (r *ResourceGroupRepo) CreateItem(item *models.ResourceGroupItem) error {
	return r.db.Create(item).Error
}

func (r *ResourceGroupRepo) UpdateItem(item *models.ResourceGroupItem) error {
	return r.db.Save(item).Error
}

func (r *ResourceGroupRepo) DeleteItem(id uint) error {
	return r.db.Delete(&models.ResourceGroupItem{}, id).Error
}

// Column defs

func (r *ResourceGroupRepo) ListColumnDefs(resourceName string) ([]models.ResourceColumnDef, error) {
	var cols []models.ResourceColumnDef
	err := r.db.Where("resource_name = ?", resourceName).Order("sort_order ASC").Find(&cols).Error
	return cols, err
}

func (r *ResourceGroupRepo) CreateColumnDef(col *models.ResourceColumnDef) error {
	return r.db.Create(col).Error
}

func (r *ResourceGroupRepo) DeleteColumnDefsByResource(resourceName string) error {
	return r.db.Where("resource_name = ?", resourceName).Delete(&models.ResourceColumnDef{}).Error
}

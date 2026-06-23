package repository

import (
	"github.com/anraaa/visual-mesin/internal/models"
	"gorm.io/gorm"
)

type DBConnectionRepository struct {
	db *gorm.DB
}

func NewDBConnectionRepository(db *gorm.DB) *DBConnectionRepository {
	return &DBConnectionRepository{db: db}
}

func (r *DBConnectionRepository) FindByID(id uint) (*models.DBConnection, error) {
	var conn models.DBConnection
	err := r.db.First(&conn, id).Error
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

func (r *DBConnectionRepository) FindByName(name string) (*models.DBConnection, error) {
	var conn models.DBConnection
	err := r.db.Where("name = ?", name).First(&conn).Error
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

func (r *DBConnectionRepository) Create(conn *models.DBConnection) error {
	return r.db.Create(conn).Error
}

func (r *DBConnectionRepository) Update(conn *models.DBConnection) error {
	return r.db.Save(conn).Error
}

func (r *DBConnectionRepository) Delete(id uint) error {
	return r.db.Delete(&models.DBConnection{}, id).Error
}

func (r *DBConnectionRepository) List(page, limit int) ([]models.DBConnection, int64, error) {
	var conns []models.DBConnection
	var total int64

	r.db.Model(&models.DBConnection{}).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&conns).Error
	if err != nil {
		return nil, 0, err
	}

	return conns, total, nil
}

func (r *DBConnectionRepository) ListActive() ([]models.DBConnection, error) {
	var conns []models.DBConnection
	err := r.db.Where("is_active = ?", true).Find(&conns).Error
	return conns, err
}

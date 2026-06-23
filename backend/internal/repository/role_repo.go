package repository

import (
	"github.com/anraaa/visual-mesin/internal/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	return &role, err
}

func (r *RoleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	return &role, err
}

func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) List() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) GetPermissions(roleID uint) ([]models.Permission, error) {
	var perms []models.Permission
	err := r.db.Raw(`
		SELECT p.* FROM permissions p
		INNER JOIN role_has_permissions rhp ON rhp.permission_id = p.id
		WHERE rhp.role_id = ?
	`, roleID).Scan(&perms).Error
	return perms, err
}

func (r *RoleRepository) Delete(role *models.Role) error {
	return r.db.Delete(role).Error
}

func (r *RoleRepository) AssignPermission(roleID, permissionID uint) error {
	return r.db.Exec(
		"INSERT INTO role_has_permissions (role_id, permission_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE role_id=role_id",
		roleID, permissionID,
	).Error
}

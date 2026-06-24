package repository

import (
	"strings"

	"github.com/anraaa/visual-mesin/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByNIP(nip string) (*models.User, error) {
	var user models.User
	err := r.db.Where("nip = ?", nip).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmailOrNIP(identifier string) (*models.User, error) {
	var user models.User

	err := r.db.Where("email = ? OR nip = ?", identifier, identifier).First(&user).Error
	if err == nil {
		return &user, nil
	}

	if !strings.Contains(identifier, "@") && identifier != "" {
		normalized := strings.ToLower(identifier)

		if normalized != identifier {
			err = r.db.Where("nip = ?", normalized).First(&user).Error
			if err == nil {
				return &user, nil
			}
		}

		if strings.HasPrefix(normalized, "m") && len(normalized) > 1 {
			err = r.db.Where("nip = ?", normalized[1:]).First(&user).Error
			if err == nil {
				return &user, nil
			}
		}

		if !strings.HasPrefix(normalized, "m") {
			err = r.db.Where("nip = ?", "m"+normalized).First(&user).Error
			if err == nil {
				return &user, nil
			}
		}
	}

	return nil, gorm.ErrRecordNotFound
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) List(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	r.db.Model(&models.User{}).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) GetRoles(userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Raw(`
		SELECT r.* FROM roles r
		INNER JOIN model_has_roles mhr ON mhr.role_id = r.id
		WHERE mhr.model_id = ? AND mhr.model_type = ?
	`, userID, "App\\Models\\User").Scan(&roles).Error
	return roles, err
}

func (r *UserRepository) GetPermissions(userID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Raw(`
		SELECT DISTINCT p.* FROM permissions p
		INNER JOIN model_has_permissions mhp ON mhp.permission_id = p.id
		WHERE mhp.model_id = ? AND mhp.model_type = ?
		UNION
		SELECT DISTINCT p.* FROM permissions p
		INNER JOIN role_has_permissions rhp ON rhp.permission_id = p.id
		INNER JOIN model_has_roles mhr ON mhr.role_id = rhp.role_id
		WHERE mhr.model_id = ? AND mhr.model_type = ?
	`, userID, "App\\Models\\User", userID, "App\\Models\\User").Scan(&permissions).Error
	return permissions, err
}

func (r *UserRepository) AssignRole(userID, roleID uint) error {
	return r.db.Exec(
		"INSERT INTO model_has_roles (role_id, model_type, model_id) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE role_id=role_id",
		roleID, "App\\Models\\User", userID,
	).Error
}

func (r *UserRepository) RevokeAllRoles(userID uint) error {
	return r.db.Exec(
		"DELETE FROM model_has_roles WHERE model_type = ? AND model_id = ?",
		"App\\Models\\User", userID,
	).Error
}

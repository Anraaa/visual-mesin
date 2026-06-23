package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/anraaa/visual-mesin/internal/db"
	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/repository"
	"github.com/anraaa/visual-mesin/pkg/utils"
	"gorm.io/gorm"
)

type ResourceDBConfigService struct {
	repo    *repository.ResourceDBConfigRepository
	dbMgr   *db.Manager
	localDB *gorm.DB
}

func NewResourceDBConfigService(repo *repository.ResourceDBConfigRepository, dbMgr *db.Manager, localDB *gorm.DB) *ResourceDBConfigService {
	return &ResourceDBConfigService{repo: repo, dbMgr: dbMgr, localDB: localDB}
}

func (s *ResourceDBConfigService) List(page, limit int) ([]models.ResourceDBConfig, int64, error) {
	configs, total, err := s.repo.List(page, limit)
	if err != nil {
		return nil, 0, err
	}
	for i := range configs {
		configs[i].Password = ""
	}
	return configs, total, nil
}

func (s *ResourceDBConfigService) GetByID(id uint) (*models.ResourceDBConfig, error) {
	cfg, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	cfg.Password = ""
	return cfg, nil
}

func (s *ResourceDBConfigService) Create(req *models.ResourceDBConfigRequest) (*models.ResourceDBConfig, error) {
	encryptedPass, err := utils.Encrypt(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %w", err)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	var label *string
	if req.Label != "" {
		label = &req.Label
	}

	cfg := &models.ResourceDBConfig{
		ResourceName: req.ResourceName,
		Label:        label,
		Driver:       req.Driver,
		Host:         req.Host,
		Port:         req.Port,
		DatabaseName: req.DatabaseName,
		Username:     req.Username,
		Password:     encryptedPass,
		IsActive:     isActive,
	}

	if err := s.repo.Create(cfg); err != nil {
		return nil, err
	}

	cfg.Password = ""
	return cfg, nil
}

func (s *ResourceDBConfigService) Update(id uint, req *models.ResourceDBConfigUpdateRequest) (*models.ResourceDBConfig, error) {
	cfg, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.ResourceName != nil {
		cfg.ResourceName = *req.ResourceName
	}
	if req.Label != nil {
		cfg.Label = req.Label
	}
	if req.Driver != nil {
		cfg.Driver = *req.Driver
	}
	if req.Host != nil {
		cfg.Host = *req.Host
	}
	if req.Port != nil {
		cfg.Port = *req.Port
	}
	if req.DatabaseName != nil {
		cfg.DatabaseName = *req.DatabaseName
	}
	if req.Username != nil {
		cfg.Username = *req.Username
	}
	if req.Password != nil {
		encrypted, err := utils.Encrypt(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt password: %w", err)
		}
		cfg.Password = encrypted
	}
	if req.IsActive != nil {
		cfg.IsActive = *req.IsActive
	}

	password := cfg.Password
	if req.Password != nil {
		password = *req.Password
	}

	dsn := db.BuildDSN(cfg.Host, cfg.Port, cfg.Username, password, cfg.DatabaseName)
	s.dbMgr.RemoveConnection(dsn)

	if err := s.repo.Update(cfg); err != nil {
		return nil, err
	}

	cfg.Password = ""
	return cfg, nil
}

func (s *ResourceDBConfigService) Delete(id uint) error {
	cfg, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	decrypted, err := utils.Decrypt(cfg.Password)
	if err == nil {
		dsn := db.BuildDSN(cfg.Host, cfg.Port, cfg.Username, decrypted, cfg.DatabaseName)
		s.dbMgr.RemoveConnection(dsn)
	}

	return s.repo.Delete(id)
}

func (s *ResourceDBConfigService) TestByID(id uint) error {
	cfg, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	decrypted, err := utils.Decrypt(cfg.Password)
	if err != nil {
		s.updateTestResult(id, false, "Failed to decrypt password")
		return errors.New("failed to decrypt password")
	}

	dsn := db.BuildDSN(cfg.Host, cfg.Port, cfg.Username, decrypted, cfg.DatabaseName)

	testDB, err := sql.Open("mysql", dsn)
	if err != nil {
		s.updateTestResult(id, false, err.Error())
		return fmt.Errorf("failed to open connection: %w", err)
	}
	defer testDB.Close()

	testDB.SetMaxOpenConns(1)
	testDB.SetConnMaxLifetime(5 * time.Second)

	if err := testDB.Ping(); err != nil {
		s.updateTestResult(id, false, err.Error())
		return fmt.Errorf("connection test failed: %w", err)
	}

	s.updateTestResult(id, true, "Connection successful")
	return nil
}

func (s *ResourceDBConfigService) updateTestResult(id uint, success bool, message string) {
	cfg, err := s.repo.FindByID(id)
	if err != nil {
		return
	}

	now := time.Now()
	cfg.IsLastTestSuccess = &success
	cfg.LastTestedAt = &now
	cfg.LastTestMessage = &message
	s.localDB.Model(cfg).Updates(map[string]interface{}{
		"is_last_test_success": success,
		"last_tested_at":       now,
		"last_test_message":    message,
	})
}

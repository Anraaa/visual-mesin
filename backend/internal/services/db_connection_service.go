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

type DBConnectionService struct {
	repo    *repository.DBConnectionRepository
	dbMgr   *db.Manager
	localDB *gorm.DB
}

func NewDBConnectionService(repo *repository.DBConnectionRepository, dbMgr *db.Manager, localDB *gorm.DB) *DBConnectionService {
	return &DBConnectionService{repo: repo, dbMgr: dbMgr, localDB: localDB}
}

func (s *DBConnectionService) List(page, limit int) ([]models.DBConnection, int64, error) {
	conns, total, err := s.repo.List(page, limit)
	if err != nil {
		return nil, 0, err
	}
	for i := range conns {
		conns[i].Password = ""
	}
	return conns, total, nil
}

func (s *DBConnectionService) GetByID(id uint) (*models.DBConnection, error) {
	conn, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	conn.Password = ""
	return conn, nil
}

func (s *DBConnectionService) Create(req *models.DBConnectionRequest) (*models.DBConnection, error) {
	encryptedPass, err := utils.Encrypt(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %w", err)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	conn := &models.DBConnection{
		Name:         req.Name,
		Driver:       req.Driver,
		Host:         req.Host,
		Port:         req.Port,
		DatabaseName: req.DatabaseName,
		Username:     req.Username,
		Password:     encryptedPass,
		IsActive:     isActive,
	}

	if err := s.repo.Create(conn); err != nil {
		return nil, err
	}

	conn.Password = ""
	return conn, nil
}

func (s *DBConnectionService) Update(id uint, req *models.DBConnectionUpdateRequest) (*models.DBConnection, error) {
	conn, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		conn.Name = *req.Name
	}
	if req.Driver != nil {
		conn.Driver = *req.Driver
	}
	if req.Host != nil {
		conn.Host = *req.Host
	}
	if req.Port != nil {
		conn.Port = *req.Port
	}
	if req.DatabaseName != nil {
		conn.DatabaseName = *req.DatabaseName
	}
	if req.Username != nil {
		conn.Username = *req.Username
	}
	if req.Password != nil {
		encrypted, err := utils.Encrypt(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt password: %w", err)
		}
		conn.Password = encrypted
	}
	if req.IsActive != nil {
		conn.IsActive = *req.IsActive
	}

	password := conn.Password
	if req.Password != nil {
		password = *req.Password
	}

	dsn := db.BuildDSN(conn.Host, conn.Port, conn.Username, password, conn.DatabaseName)
	s.dbMgr.RemoveConnection(dsn)

	if err := s.repo.Update(conn); err != nil {
		return nil, err
	}

	conn.Password = ""
	return conn, nil
}

func (s *DBConnectionService) Delete(id uint) error {
	conn, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	decrypted, err := utils.Decrypt(conn.Password)
	if err == nil {
		dsn := db.BuildDSN(conn.Host, conn.Port, conn.Username, decrypted, conn.DatabaseName)
		s.dbMgr.RemoveConnection(dsn)
	}

	return s.repo.Delete(id)
}

func (s *DBConnectionService) TestConnection(req *models.TestConnectionRequest) error {
	dsn := db.BuildDSN(req.Host, req.Port, req.Username, req.Password, req.DatabaseName)

	testDB, err := sql.Open("mysql", dsn)
	if err != nil {
		s.updateTestResult(0, false, err.Error())
		return fmt.Errorf("failed to open connection: %w", err)
	}
	defer testDB.Close()

	testDB.SetMaxOpenConns(1)
	testDB.SetConnMaxLifetime(5 * time.Second)

	if err := testDB.Ping(); err != nil {
		s.updateTestResult(0, false, err.Error())
		return fmt.Errorf("connection test failed: %w", err)
	}

	s.updateTestResult(0, true, "Connection successful")
	return nil
}

func (s *DBConnectionService) TestByID(id uint) error {
	conn, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	decrypted, err := utils.Decrypt(conn.Password)
	if err != nil {
		s.updateTestResult(id, false, "Failed to decrypt password")
		return errors.New("failed to decrypt password")
	}

	dsn := db.BuildDSN(conn.Host, conn.Port, conn.Username, decrypted, conn.DatabaseName)

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

func (s *DBConnectionService) updateTestResult(id uint, success bool, message string) {
	if id == 0 {
		return
	}

	conn, err := s.repo.FindByID(id)
	if err != nil {
		return
	}

	now := time.Now()
	conn.IsLastTestSuccess = &success
	conn.LastTestedAt = &now
	conn.LastTestMessage = &message
	s.localDB.Model(conn).Updates(map[string]interface{}{
		"is_last_test_success": success,
		"last_tested_at":       now,
		"last_test_message":    message,
	})
}

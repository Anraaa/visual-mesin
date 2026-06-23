package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Manager struct {
	pools     map[string]*sql.DB
	gormPools map[string]*gorm.DB
	mu        sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		pools:     make(map[string]*sql.DB),
		gormPools: make(map[string]*gorm.DB),
	}
}

func BuildDSN(host string, port int, username, password, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		username, password, host, port, dbName)
}

func (m *Manager) GetConnection(dsn string) (*sql.DB, error) {
	m.mu.RLock()
	if db, ok := m.pools[dsn]; ok {
		m.mu.RUnlock()
		if err := db.Ping(); err == nil {
			return db, nil
		}
		m.mu.RUnlock()
		m.mu.Lock()
		delete(m.pools, dsn)
		delete(m.gormPools, dsn)
		m.mu.Unlock()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping: %w", err)
	}

	go m.healthCheck(dsn, db)
	m.pools[dsn] = db
	return db, nil
}

func (m *Manager) GetGORM(dsn string) (*gorm.DB, error) {
	m.mu.RLock()
	if gdb, ok := m.gormPools[dsn]; ok {
		m.mu.RUnlock()
		return gdb, nil
	}
	m.mu.RUnlock()

	sqlDB, err := m.GetConnection(dsn)
	if err != nil {
		return nil, err
	}

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm: %w", err)
	}

	m.mu.Lock()
	m.gormPools[dsn] = gdb
	m.mu.Unlock()

	return gdb, nil
}

func (m *Manager) RemoveConnection(dsn string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if db, ok := m.pools[dsn]; ok {
		db.Close()
		delete(m.pools, dsn)
	}
	delete(m.gormPools, dsn)
}

func (m *Manager) healthCheck(dsn string, db *sql.DB) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.RLock()
		current, ok := m.pools[dsn]
		m.mu.RUnlock()

		if !ok || current != db {
			return
		}

		if err := db.Ping(); err != nil {
			m.mu.Lock()
			delete(m.pools, dsn)
			delete(m.gormPools, dsn)
			m.mu.Unlock()
			db.Close()
			return
		}
	}
}

func (m *Manager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for dsn, db := range m.pools {
		db.Close()
		delete(m.pools, dsn)
		delete(m.gormPools, dsn)
	}
}

package services

import (
	"fmt"

	"github.com/anraaa/visual-mesin/internal/db"
	"github.com/anraaa/visual-mesin/internal/repository"
	"github.com/anraaa/visual-mesin/pkg/utils"
	"gorm.io/gorm"
)

type ResourceQueryService struct {
	resourceRepo *repository.ResourceDBConfigRepository
	dbMgr        *db.Manager
}

func NewResourceQueryService(resourceRepo *repository.ResourceDBConfigRepository, dbMgr *db.Manager) *ResourceQueryService {
	return &ResourceQueryService{resourceRepo: resourceRepo, dbMgr: dbMgr}
}

type QueryParams struct {
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
	Filters  map[string]string `json:"filters"`
	SortBy   string            `json:"sort_by"`
	SortDir  string            `json:"sort_dir"`
	Search   string            `json:"search"`
	SearchBy string            `json:"search_by"`
}

type QueryResult struct {
	Data     []map[string]interface{} `json:"data"`
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	Limit    int                      `json:"limit"`
	LastPage int                      `json:"last_page"`
	Columns  []ColumnInfo             `json:"columns"`
}

type ColumnInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (s *ResourceQueryService) getGormDB(resourceName string) (*gorm.DB, error) {
	cfg, err := s.resourceRepo.FindByResourceName(resourceName)
	if err != nil {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}
	if !cfg.IsActive {
		return nil, fmt.Errorf("resource is not active: %s", resourceName)
	}
	decrypted, err := utils.Decrypt(cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}
	dsn := db.BuildDSN(cfg.Host, cfg.Port, cfg.Username, decrypted, cfg.DatabaseName)
	gdb, err := s.dbMgr.GetGORM(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to resource db: %w", err)
	}
	return gdb, nil
}

func (s *ResourceQueryService) QueryResource(resourceName string, params QueryParams) (*QueryResult, error) {
	gdb, err := s.getGormDB(resourceName)
	if err != nil {
		return nil, err
	}

	columns, err := getTableColumns(gdb, resourceName)
	if err != nil {
		columns = []ColumnInfo{}
	}

	q := gdb.Table(resourceName)

	if params.Search != "" && params.SearchBy != "" {
		q = q.Where(fmt.Sprintf("%s LIKE ?", params.SearchBy), "%"+params.Search+"%")
	}
	if params.SortBy != "" {
		dir := "ASC"
		if params.SortDir == "desc" {
			dir = "DESC"
		}
		q = q.Order(fmt.Sprintf("%s %s", params.SortBy, dir))
	}

	var total int64
	q.Count(&total)

	page := params.Page
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	if limit < 1 || limit > 1000 {
		limit = 25
	}
	offset := (page - 1) * limit

	var results []map[string]interface{}
	if err := q.Offset(offset).Limit(limit).Find(&results).Error; err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	lastPage := int(total) / limit
	if int(total)%limit > 0 {
		lastPage++
	}

	return &QueryResult{
		Data:     results,
		Total:    total,
		Page:     page,
		Limit:    limit,
		LastPage: lastPage,
		Columns:  columns,
	}, nil
}

func (s *ResourceQueryService) QueryByID(resourceName string, idColumn string, idValue interface{}) (map[string]interface{}, error) {
	gdb, err := s.getGormDB(resourceName)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := gdb.Table(resourceName).Where(fmt.Sprintf("%s = ?", idColumn), idValue).First(&result).Error; err != nil {
		return nil, fmt.Errorf("record not found: %w", err)
	}
	return result, nil
}

func (s *ResourceQueryService) Create(resourceName string, data map[string]interface{}) (map[string]interface{}, error) {
	gdb, err := s.getGormDB(resourceName)
	if err != nil {
		return nil, err
	}

	if err := gdb.Table(resourceName).Create(data).Error; err != nil {
		return nil, fmt.Errorf("failed to create record: %w", err)
	}

	if id, ok := data["id"]; ok {
		var result map[string]interface{}
		gdb.Table(resourceName).Where("id = ?", id).First(&result)
		return result, nil
	}
	return data, nil
}

func (s *ResourceQueryService) Update(resourceName string, idColumn string, idValue interface{}, data map[string]interface{}) (map[string]interface{}, error) {
	gdb, err := s.getGormDB(resourceName)
	if err != nil {
		return nil, err
	}

	delete(data, idColumn)

	if err := gdb.Table(resourceName).Where(fmt.Sprintf("%s = ?", idColumn), idValue).Updates(data).Error; err != nil {
		return nil, fmt.Errorf("failed to update record: %w", err)
	}

	var result map[string]interface{}
	if err := gdb.Table(resourceName).Where(fmt.Sprintf("%s = ?", idColumn), idValue).First(&result).Error; err != nil {
		return nil, fmt.Errorf("record not found after update: %w", err)
	}
	return result, nil
}

func (s *ResourceQueryService) Delete(resourceName string, idColumn string, idValue interface{}) error {
	gdb, err := s.getGormDB(resourceName)
	if err != nil {
		return err
	}

	if err := gdb.Table(resourceName).Where(fmt.Sprintf("%s = ?", idColumn), idValue).Delete(nil).Error; err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}
	return nil
}

func getTableColumns(gdb *gorm.DB, tableName string) ([]ColumnInfo, error) {
	var columns []ColumnInfo

	query := fmt.Sprintf("SHOW COLUMNS FROM `%s`", tableName)
	rows, err := gdb.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var field, typeStr, null, key, extra string
		var defaultVal interface{}
		if err := rows.Scan(&field, &typeStr, &null, &key, &defaultVal, &extra); err != nil {
			continue
		}
		columns = append(columns, ColumnInfo{
			Name: field,
			Type: typeStr,
		})
	}

	return columns, nil
}

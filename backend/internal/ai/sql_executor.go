package ai

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type ExecutorResult struct {
	Columns []string                   `json:"columns"`
	Rows    []map[string]interface{}   `json:"rows"`
	Total   int                        `json:"total"`
	Latency string                     `json:"latency"`
}

type SQLExecutor struct {
	getDB func(driver, host string, port int, username, password, dbName string) (*sql.DB, error)
}

func NewSQLExecutor(getDB func(driver, host string, port int, username, password, dbName string) (*sql.DB, error)) *SQLExecutor {
	return &SQLExecutor{getDB: getDB}
}

func (e *SQLExecutor) Execute(sqlStr string, dbConfig *DBConfig) (*ExecutorResult, error) {
	start := time.Now()

	db, err := e.getDB(dbConfig.Driver, dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(10 * time.Second)

	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []map[string]interface{}
	rowCount := 0
	maxRows := 100

	for rows.Next() {
		if rowCount >= maxRows {
			break
		}

		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			continue
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			switch v := val.(type) {
			case []byte:
				row[col] = string(v)
			case time.Time:
				row[col] = v.Format("2006-01-02 15:04:05")
			default:
				row[col] = v
			}
		}
		results = append(results, row)
		rowCount++
	}

	latency := time.Since(start).Round(time.Millisecond).String()

	var totalCount int
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_subquery", sqlStr)
	if err := db.QueryRow(countSQL).Scan(&totalCount); err != nil {
		totalCount = rowCount
	}

	return &ExecutorResult{
		Columns: columns,
		Rows:    results,
		Total:   totalCount,
		Latency: latency,
	}, nil
}

type DBConfig struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func (e *SQLExecutor) BuildDBConfigFromJSON(jsonData string) (*DBConfig, error) {
	var cfg DBConfig
	if err := json.Unmarshal([]byte(jsonData), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

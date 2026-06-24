package services

import (
	"fmt"
	"strings"

	"github.com/anraaa/visual-mesin/internal/db"
	"github.com/anraaa/visual-mesin/internal/repository"
	"github.com/anraaa/visual-mesin/pkg/utils"
	"gorm.io/gorm"
)

type ResourceQueryService struct {
	resourceRepo *repository.ResourceDBConfigRepository
	dbMgr        *db.Manager
	localDB      *gorm.DB
}

func NewResourceQueryService(resourceRepo *repository.ResourceDBConfigRepository, dbMgr *db.Manager, localDB *gorm.DB) *ResourceQueryService {
	return &ResourceQueryService{resourceRepo: resourceRepo, dbMgr: dbMgr, localDB: localDB}
}

type QueryParams struct {
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	Filters    map[string]string `json:"filters"`
	SortBy     string            `json:"sort_by"`
	SortDir    string            `json:"sort_dir"`
	Search     string            `json:"search"`
	SearchBy   string            `json:"search_by"`
	StartDate  string            `json:"start_date"`
	EndDate    string            `json:"end_date"`
	DateColumn string            `json:"date_column"`
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
		if s.localDB != nil && s.localDB.Migrator().HasTable(resourceName) {
			return s.localDB, nil
		}
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
	if params.StartDate != "" && params.DateColumn != "" {
		q = q.Where(fmt.Sprintf("%s >= ?", params.DateColumn), params.StartDate)
	}
	if params.EndDate != "" && params.DateColumn != "" {
		q = q.Where(fmt.Sprintf("%s <= ?", params.DateColumn), params.EndDate)
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

type SPCResult struct {
	Labels   []string      `json:"labels"`
	Datasets []SPCDataset  `json:"datasets"`
}

type SPCDataset struct {
	Label        string     `json:"label"`
	Data         []*float64 `json:"data"`
	BorderColor  string     `json:"borderColor"`
	BorderDash   []int      `json:"borderDash,omitempty"`
	Fill         bool       `json:"fill"`
	Tension      float64    `json:"tension,omitempty"`
	PointRadius  int        `json:"pointRadius"`
}

type SPCConfig struct {
	TimeColumn string `json:"time_column"`
	Actual     string `json:"actual"`
	Target     string `json:"target"`
	TolPP      string `json:"tol_pp,omitempty"`
	TolMM      string `json:"tol_mm,omitempty"`
	TolP       string `json:"tol_p,omitempty"`
	TolM       string `json:"tol_m,omitempty"`
	Limit      int    `json:"limit"`
}

func (s *ResourceQueryService) GetSPC(resourceName string, cfg *SPCConfig) (*SPCResult, error) {
	gdb, err := s.getGormDB(resourceName)
	if err != nil {
		return nil, err
	}

	limit := cfg.Limit
	if limit < 1 || limit > 200 {
		limit = 50
	}

	tolCols := []string{}
	if cfg.TolPP != "" { tolCols = append(tolCols, cfg.TolPP) }
	if cfg.TolMM != "" { tolCols = append(tolCols, cfg.TolMM) }
	if cfg.TolP != ""  { tolCols = append(tolCols, cfg.TolP) }
	if cfg.TolM != ""  { tolCols = append(tolCols, cfg.TolM) }

	allCols := append([]string{cfg.TimeColumn, cfg.Actual, cfg.Target}, tolCols...)
	selectParts := ""
	for i, c := range allCols {
		if i > 0 { selectParts += ", " }
		selectParts += fmt.Sprintf("`%s`", c)
	}

	query := fmt.Sprintf("SELECT %s FROM `%s` WHERE `%s` IS NOT NULL AND `%s` IS NOT NULL ORDER BY `%s` DESC LIMIT %d",
		selectParts, resourceName, cfg.Actual, cfg.Target, cfg.TimeColumn, limit)

	rows, err := gdb.Raw(query).Rows()
	if err != nil {
		return nil, fmt.Errorf("spc query failed: %w", err)
	}
	defer rows.Close()

	var rawRows []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		if err := gdb.ScanRows(rows, &row); err != nil {
			continue
		}
		rawRows = append(rawRows, row)
	}

	datasets := []SPCDataset{
		{Label: "Actual Weight", BorderColor: "#10b981", Fill: false, Tension: 0.3, PointRadius: 3},
		{Label: "Set Weight", BorderColor: "#9ca3af", BorderDash: []int{5, 5}, Fill: false, PointRadius: 0},
		{Label: "USL", BorderColor: "#ef4444", BorderDash: []int{6, 3}, Fill: false, PointRadius: 0},
		{Label: "LSL", BorderColor: "#ef4444", BorderDash: []int{6, 3}, Fill: false, PointRadius: 0},
		{Label: "UCL (Warning)", BorderColor: "#f59e0b", BorderDash: []int{2, 2}, Fill: false, PointRadius: 0},
		{Label: "LCL (Warning)", BorderColor: "#f59e0b", BorderDash: []int{2, 2}, Fill: false, PointRadius: 0},
	}

	result := &SPCResult{Datasets: datasets}

	for i := len(rawRows) - 1; i >= 0; i-- {
		row := rawRows[i]
		label := fmt.Sprintf("%v", row[cfg.TimeColumn])
		if len(label) > 16 {
			label = label[:16]
		}
		result.Labels = append(result.Labels, label)

		actual := toFloatPtr(row[cfg.Actual])
		target := toFloatPtr(row[cfg.Target])

		result.Datasets[0].Data = append(result.Datasets[0].Data, actual)
		result.Datasets[1].Data = append(result.Datasets[1].Data, target)

		usl := calcLimit(target, row, cfg.TolPP, true)
		lsl := calcLimit(target, row, cfg.TolMM, false)
		ucl := calcLimit(target, row, cfg.TolP, true)
		lcl := calcLimit(target, row, cfg.TolM, false)

		result.Datasets[2].Data = append(result.Datasets[2].Data, usl)
		result.Datasets[3].Data = append(result.Datasets[3].Data, lsl)
		result.Datasets[4].Data = append(result.Datasets[4].Data, ucl)
		result.Datasets[5].Data = append(result.Datasets[5].Data, lcl)
	}

	return result, nil
}

func calcLimit(target *float64, row map[string]interface{}, tolCol string, add bool) *float64 {
	if target == nil || tolCol == "" {
		return nil
	}
	tol := toFloatPtr(row[tolCol])
	if tol == nil {
		return nil
	}
	v := *target + *tol
	if !add {
		v = *target - *tol
	}
	return &v
}

func toFloatPtr(v interface{}) *float64 {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case float64:
		return &val
	case float32:
		f := float64(val)
		return &f
	case int:
		f := float64(val)
		return &f
	case int64:
		f := float64(val)
		return &f
	case string:
		if val == "" {
			return nil
		}
		var f float64
		if _, err := fmt.Sscanf(val, "%f", &f); err == nil {
			return &f
		}
		return nil
	default:
		return nil
	}
}

type ResourceStats struct {
	TotalRecords int64                  `json:"total_records"`
	CTSummary    map[string]float64     `json:"ct_summary"`
	AvgCT        float64                `json:"avg_ct"`
	MinCT        float64                `json:"min_ct"`
	MaxCT        float64                `json:"max_ct"`
}

type TrendPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Count int64   `json:"count"`
}

type DistributionItem struct {
	Label string `json:"label"`
	Value int64  `json:"value"`
}

func (s *ResourceQueryService) GetStats(resourceName string, durationCol string) (*ResourceStats, error) {
	gdb, err := s.getGormDB(resourceName)
	if err != nil {
		return nil, err
	}

	columns, err := getTableColumns(gdb, resourceName)
	if err != nil {
		return nil, err
	}

	ctColumns := []string{}
	prodColumns := []string{}
	durCol := ""
	for _, col := range columns {
		n := col.Name
		if n == "PUD_CT" || n == "BTD_CT" || n == "BDD_CT" || n == "GT_CT" {
			ctColumns = append(ctColumns, n)
		}
		if n == "prod_OK" || n == "prod_NG" || n == "prod_Tot" ||
			n == "OK" || n == "NG" || n == "Total" ||
			n == "ok" || n == "ng" || n == "total" {
			prodColumns = append(prodColumns, n)
		}
		if durationCol != "" && n == durationCol {
			durCol = n
		}
	}

	stats := &ResourceStats{
		CTSummary: make(map[string]float64),
	}

	gdb.Table(resourceName).Select("COUNT(*)").Scan(&stats.TotalRecords)

	for _, col := range ctColumns {
		var avg float64
		gdb.Table(resourceName).Select(fmt.Sprintf("AVG(%s)", col)).Scan(&avg)
		stats.CTSummary[col] = avg
	}

	for _, col := range prodColumns {
		var sum float64
		gdb.Table(resourceName).Select(fmt.Sprintf("SUM(%s)", col)).Scan(&sum)
		stats.CTSummary[col] = sum
	}

	if durCol != "" {
		var avgDur, minDur, maxDur float64
		gdb.Table(resourceName).Select(fmt.Sprintf("AVG(TIME_TO_SEC(%s))", durCol)).Scan(&avgDur)
		gdb.Table(resourceName).Select(fmt.Sprintf("MIN(TIME_TO_SEC(%s))", durCol)).Scan(&minDur)
		gdb.Table(resourceName).Select(fmt.Sprintf("MAX(TIME_TO_SEC(%s))", durCol)).Scan(&maxDur)
		stats.CTSummary["duration_avg"] = avgDur
		stats.CTSummary["duration_min"] = minDur
		stats.CTSummary["duration_max"] = maxDur
	}

	if len(ctColumns) > 0 {
		var avgCT, minCT, maxCT float64
		var hasMin, hasMax bool
		for _, col := range ctColumns {
			var a float64
			gdb.Table(resourceName).Select(fmt.Sprintf("AVG(%s)", col)).Scan(&a)
			avgCT += a

			var mi, ma float64
			if err := gdb.Table(resourceName).Select(fmt.Sprintf("MIN(%s)", col)).Scan(&mi).Error; err == nil {
				if !hasMin || mi < minCT {
					minCT = mi
					hasMin = true
				}
			}
			if err := gdb.Table(resourceName).Select(fmt.Sprintf("MAX(%s)", col)).Scan(&ma).Error; err == nil {
				if !hasMax || ma > maxCT {
					maxCT = ma
					hasMax = true
				}
			}
		}
		if len(ctColumns) > 0 {
			stats.AvgCT = avgCT / float64(len(ctColumns))
		}
		if hasMin {
			stats.MinCT = minCT
		}
		if hasMax {
			stats.MaxCT = maxCT
		}
	}

	return stats, nil
}

func (s *ResourceQueryService) GetTrend(resourceName string, column string, timeColumn string) ([]TrendPoint, error) {
	gdb, err := s.getGormDB(resourceName)
	if err != nil {
		return nil, err
	}

	cols, err := getTableColumns(gdb, resourceName)
	if err != nil {
		return nil, err
	}

	colExists := false
	for _, c := range cols {
		if c.Name == column {
			colExists = true
			break
		}
	}
	if !colExists {
		return []TrendPoint{}, nil
	}

	if timeColumn == "" {
		timeColumn = "Timestamp"
		timeExists := false
		for _, c := range cols {
			if c.Name == timeColumn {
				timeExists = true
				break
			}
		}
		if !timeExists {
			for _, c := range cols {
				lower := strings.ToLower(c.Type)
				if strings.Contains(lower, "datetime") || strings.Contains(lower, "timestamp") {
					timeColumn = c.Name
					timeExists = true
					break
				}
			}
		}
		if !timeExists {
			return []TrendPoint{}, nil
		}
	}

	valueExpr := fmt.Sprintf("AVG(`%s`)", column)
	for _, c := range cols {
		if c.Name == column {
			lower := strings.ToLower(c.Type)
			if strings.Contains(lower, "time") {
				valueExpr = fmt.Sprintf("AVG(TIME_TO_SEC(`%s`))", column)
			}
			break
		}
	}

	var results []TrendPoint
	query := fmt.Sprintf("SELECT DATE_FORMAT(`%s`, '%%Y-%%m-%%d %%H:00') as label, %s as value, COUNT(*) as count FROM `%s` WHERE `%s` IS NOT NULL GROUP BY label ORDER BY label ASC",
		timeColumn, valueExpr, resourceName, column)
	if err := gdb.Raw(query).Scan(&results).Error; err != nil {
		return []TrendPoint{}, nil
	}

	return results, nil
}

func (s *ResourceQueryService) GetDistribution(resourceName string, column string) ([]DistributionItem, error) {
	gdb, err := s.getGormDB(resourceName)
	if err != nil {
		return nil, err
	}

	cols, err := getTableColumns(gdb, resourceName)
	if err != nil {
		return nil, err
	}

	colExists := false
	for _, c := range cols {
		if c.Name == column {
			colExists = true
			break
		}
	}
	if !colExists {
		return []DistributionItem{}, nil
	}

	var results []DistributionItem
	if err := gdb.Table(resourceName).
		Select(fmt.Sprintf("`%s` as label, COUNT(*) as value", column)).
		Group(column).
		Order("value DESC").
		Limit(20).
		Scan(&results).Error; err != nil {
		return []DistributionItem{}, nil
	}

	return results, nil
}

type QualityTrendPoint struct {
	Label string  `json:"label"`
	OK    int64   `json:"ok"`
	NG    int64   `json:"ng"`
	Total int64   `json:"total"`
	Rate  float64 `json:"rate"`
}

func (s *ResourceQueryService) GetQualityTrend(resourceName string, timeCol string, statusCol string, okValue string, ngValue string) ([]QualityTrendPoint, error) {
	gdb, err := s.getGormDB(resourceName)
	if err != nil {
		return nil, err
	}

	cols, err := getTableColumns(gdb, resourceName)
	if err != nil {
		return nil, err
	}

	timeExists, statusExists := false, false
	for _, c := range cols {
		if c.Name == timeCol {
			timeExists = true
		}
		if c.Name == statusCol {
			statusExists = true
		}
	}
	if !timeExists || !statusExists {
		return []QualityTrendPoint{}, nil
	}

	if okValue == "" {
		okValue = "OK"
	}
	if ngValue == "" {
		ngValue = "NG"
	}

	var results []QualityTrendPoint
	query := fmt.Sprintf(
		"SELECT DATE_FORMAT(`%s`, '%%Y-%%m-%%d %%H:00') as label, "+
			"SUM(CASE WHEN `%s` = '%s' THEN 1 ELSE 0 END) as ok, "+
			"SUM(CASE WHEN `%s` = '%s' THEN 1 ELSE 0 END) as ng, "+
			"COUNT(*) as total, "+
			"ROUND(SUM(CASE WHEN `%s` = '%s' THEN 1 ELSE 0 END) / COUNT(*) * 100, 1) as rate "+
			"FROM `%s` WHERE `%s` IS NOT NULL GROUP BY label ORDER BY label ASC",
		timeCol, statusCol, okValue, statusCol, ngValue, statusCol, okValue, resourceName, timeCol,
	)

	if err := gdb.Raw(query).Scan(&results).Error; err != nil {
		return []QualityTrendPoint{}, nil
	}

	return results, nil
}

type JudgmentSummaryItem struct {
	Parameter string `json:"parameter"`
	LeftNG    int64  `json:"left_ng"`
	RightNG   int64  `json:"right_ng"`
	TotalNG   int64  `json:"total_ng"`
}

func (s *ResourceQueryService) GetJudgmentSummary(resourceName string) ([]JudgmentSummaryItem, error) {
	gdb, err := s.getGormDB(resourceName)
	if err != nil {
		return nil, err
	}

	cols, err := getTableColumns(gdb, resourceName)
	if err != nil {
		return nil, err
	}
	colMap := make(map[string]bool)
	for _, c := range cols {
		colMap[c.Name] = true
	}

	params := []string{"extemp", "platen", "jacket", "intemp", "inpressN2", "inpressSt"}
	var result []JudgmentSummaryItem

	for _, p := range params {
		leftCol := p + "_jdgL"
		rightCol := p + "_jdgR"
		if !colMap[leftCol] || !colMap[rightCol] {
			continue
		}
		var item JudgmentSummaryItem
		query := fmt.Sprintf(
			"SELECT ? as parameter, "+
				"SUM(CASE WHEN `%s` = 'NG' THEN 1 ELSE 0 END) as left_ng, "+
				"SUM(CASE WHEN `%s` = 'NG' THEN 1 ELSE 0 END) as right_ng "+
				"FROM `%s`",
			leftCol, rightCol, resourceName,
		)
		gdb.Raw(query, p).Scan(&item)
		item.TotalNG = item.LeftNG + item.RightNG
		result = append(result, item)
	}

	return result, nil
}

func (s *ResourceQueryService) GetResourceCount(resourceName string) (int64, error) {
	gdb, err := s.getGormDB(resourceName)
	if err != nil {
		return 0, err
	}
	var count int64
	if err := gdb.Table(resourceName).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
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

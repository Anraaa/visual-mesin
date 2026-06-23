package services

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/anraaa/visual-mesin/internal/db"
	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/repository"
)

const exportQueueKey = "export:queue"
const exportResultPrefix = "export:result:"

type ExportService struct {
	repo      *repository.ExportJobRepository
	querySvc  *ResourceQueryService
	rdb       *redis.Client
	dbMgr     *db.Manager
	exportDir string
}

func NewExportService(
	repo *repository.ExportJobRepository,
	querySvc *ResourceQueryService,
	rdb *redis.Client,
	dbMgr *db.Manager,
	exportDir string,
) *ExportService {
	return &ExportService{
		repo:      repo,
		querySvc:  querySvc,
		rdb:       rdb,
		dbMgr:     dbMgr,
		exportDir: exportDir,
	}
}

func (s *ExportService) SubmitJob(userID uint, req *models.ExportJobRequest) (*models.ExportJob, error) {
	columnsJSON, _ := json.Marshal(req.Columns)
	filtersJSON, _ := json.Marshal(req.Filters)

	format := req.Format
	if format == "" {
		format = "csv"
	}

	job := &models.ExportJob{
		ResourceName: req.ResourceName,
		UserID:       userID,
		Status:       "queued",
		Format:       format,
		Columns:      strPtr(string(columnsJSON)),
		Filters:      strPtr(string(filtersJSON)),
	}

	if err := s.repo.Create(job); err != nil {
		return nil, fmt.Errorf("failed to create export job: %w", err)
	}

	jobData, _ := json.Marshal(map[string]interface{}{
		"job_id":        job.ID,
		"resource_name": job.ResourceName,
		"format":        job.Format,
		"columns":       req.Columns,
		"search":        req.Search,
		"search_by":     req.SearchBy,
		"filters":       req.Filters,
	})

	if err := s.rdb.LPush(context.Background(), exportQueueKey, jobData).Err(); err != nil {
		s.failJob(job, fmt.Sprintf("failed to queue: %v", err))
		return nil, fmt.Errorf("failed to queue export job: %w", err)
	}

	return job, nil
}

func (s *ExportService) StartWorker() {
	go func() {
		for {
			result, err := s.rdb.BRPop(context.Background(), 5*time.Second, exportQueueKey).Result()
			if err != nil {
				continue
			}

			if len(result) < 2 {
				continue
			}

			var payload map[string]interface{}
			if err := json.Unmarshal([]byte(result[1]), &payload); err != nil {
				continue
			}

			jobID := uint(payload["job_id"].(float64))
			s.processJob(jobID)
		}
	}()
}

func (s *ExportService) processJob(jobID uint) {
	job, err := s.repo.FindByID(jobID)
	if err != nil {
		return
	}

	now := time.Now()
	job.Status = "processing"
	job.StartedAt = &now
	s.repo.Update(job)

	var columns []string
	if job.Columns != nil && *job.Columns != "" {
		json.Unmarshal([]byte(*job.Columns), &columns)
	}

	params := QueryParams{
		Page:   1,
		Limit:  1000,
		Search: "",
	}

	queryResult, err := s.querySvc.QueryResource(job.ResourceName, params)
	if err != nil {
		s.failJob(job, fmt.Sprintf("query failed: %v", err))
		return
	}

	job.TotalRows = queryResult.Total

	if queryResult.Total == 0 {
		completed := time.Now()
		job.Status = "completed"
		job.CompletedAt = &completed
		job.ProcessedRows = 0
		s.repo.Update(job)
		return
	}

	fileName := fmt.Sprintf("%s_%d_%s.%s", job.ResourceName, job.ID, time.Now().Format("20060102150405"), job.Format)
	filePath := filepath.Join(s.exportDir, fileName)

	os.MkdirAll(s.exportDir, 0755)

	file, err := os.Create(filePath)
	if err != nil {
		s.failJob(job, fmt.Sprintf("failed to create file: %v", err))
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headerWritten := false
	totalWritten := int64(0)
	page := 1

	for {
		params.Page = page
		params.Limit = 1000

		result, err := s.querySvc.QueryResource(job.ResourceName, params)
		if err != nil {
			break
		}

		if len(result.Data) == 0 {
			break
		}

		for _, row := range result.Data {
			if !headerWritten {
				var headers []string
				if len(columns) > 0 {
					headers = columns
				} else {
					for k := range row {
						headers = append(headers, k)
					}
				}
				writer.Write(headers)
				headerWritten = true
			}

			var record []string
			if len(columns) > 0 {
				for _, col := range columns {
					record = append(record, fmt.Sprintf("%v", row[col]))
				}
			} else {
				for _, v := range row {
					record = append(record, fmt.Sprintf("%v", v))
				}
			}
			writer.Write(record)
			totalWritten++
		}

		job.ProcessedRows = totalWritten
		s.repo.Update(job)

		if len(result.Data) < 1000 {
			break
		}
		page++
	}

	writer.Flush()
	file.Close()

	fileInfo, _ := os.Stat(filePath)
	if fileInfo != nil {
		size := fileInfo.Size()
		job.FileSize = &size
	}

	relPath := "/exports/" + fileName
	job.FilePath = &relPath
	completed := time.Now()
	job.Status = "completed"
	job.CompletedAt = &completed
	job.ProcessedRows = totalWritten
	s.repo.Update(job)
}

func (s *ExportService) GetJobStatus(jobID uint) (*models.ExportJobStatusResponse, error) {
	job, err := s.repo.FindByID(jobID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(job), nil
}

func (s *ExportService) ListJobs(userID uint, page, limit int) ([]models.ExportJobStatusResponse, int64, error) {
	jobs, total, err := s.repo.ListByUser(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]models.ExportJobStatusResponse, len(jobs))
	for i, j := range jobs {
		resp[i] = *s.toResponse(&j)
	}

	return resp, total, nil
}

func (s *ExportService) DeleteJob(jobID uint) error {
	job, err := s.repo.FindByID(jobID)
	if err != nil {
		return err
	}

	if job.FilePath != nil && *job.FilePath != "" {
		os.Remove(filepath.Join(s.exportDir, filepath.Base(*job.FilePath)))
	}

	return nil
}

func (s *ExportService) toResponse(job *models.ExportJob) *models.ExportJobStatusResponse {
	resp := &models.ExportJobStatusResponse{
		ID:            job.ID,
		ResourceName:  job.ResourceName,
		Status:        job.Status,
		TotalRows:     job.TotalRows,
		ProcessedRows: job.ProcessedRows,
		FileSize:      job.FileSize,
		Format:        job.Format,
		ErrorMessage:  job.ErrorMessage,
		CreatedAt:     job.CreatedAt,
		CompletedAt:   job.CompletedAt,
	}

	if job.FilePath != nil {
		url := *job.FilePath
		resp.FileURL = &url
	}

	return resp
}

func (s *ExportService) failJob(job *models.ExportJob, msg string) {
	job.Status = "failed"
	errMsg := msg
	job.ErrorMessage = &errMsg
	s.repo.Update(job)
}

func strPtr(s string) *string {
	return &s
}

package handlers

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/services"
)

type ExportHandler struct {
	svc *services.ExportService
}

func NewExportHandler(svc *services.ExportService) *ExportHandler {
	return &ExportHandler{svc: svc}
}

func (h *ExportHandler) Submit(c *gin.Context) {
	var req models.ExportJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	userID, _ := c.Get("user_id")

	job, err := h.svc.SubmitJob(userID.(uint), &req)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Export job submitted",
		Data:    job,
	})
}

func (h *ExportHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))
	userID, _ := c.Get("user_id")

	jobs, total, err := h.svc.ListJobs(userID.(uint), page, limit)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch export jobs")
		return
	}

	middleware.SuccessPaginated(c, "Export jobs retrieved", jobs, total, page, limit)
}

func (h *ExportHandler) GetStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "Invalid ID")
		return
	}

	job, err := h.svc.GetJobStatus(uint(id))
	if err != nil {
		middleware.NotFoundResponse(c, "Export job not found")
		return
	}

	middleware.SuccessResponse(c, "Export job status", job)
}

func (h *ExportHandler) Download(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "Invalid ID")
		return
	}

	job, err := h.svc.GetJobStatus(uint(id))
	if err != nil || job.FileURL == nil {
		middleware.NotFoundResponse(c, "File not found")
		return
	}

	filename := filepath.Base(*job.FileURL)
	filePath := "exports/" + filename

	c.FileAttachment(filePath, filename)
}

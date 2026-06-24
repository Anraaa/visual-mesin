package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/services"
)

type ResourceDBConfigHandler struct {
	svc            *services.ResourceDBConfigService
	activityLogSvc *services.ActivityLogService
}

func NewResourceDBConfigHandler(svc *services.ResourceDBConfigService, activityLogSvc *services.ActivityLogService) *ResourceDBConfigHandler {
	return &ResourceDBConfigHandler{svc: svc, activityLogSvc: activityLogSvc}
}

func (h *ResourceDBConfigHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))

	configs, total, err := h.svc.List(page, limit)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch resource database configs")
		return
	}

	middleware.SuccessPaginated(c, "Resource database configs retrieved successfully", configs, total, page, limit)
}

func (h *ResourceDBConfigHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "Invalid ID")
		return
	}

	cfg, err := h.svc.GetByID(uint(id))
	if err != nil {
		middleware.NotFoundResponse(c, "Resource database config not found")
		return
	}

	middleware.SuccessResponse(c, "Resource database config retrieved successfully", cfg)
}

func (h *ResourceDBConfigHandler) Create(c *gin.Context) {
	var req models.ResourceDBConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	cfg, err := h.svc.Create(&req)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	logActivity(c, h.activityLogSvc, "resource_db_config", "Resource DB config created: "+req.ResourceName, "create", map[string]interface{}{
		"config_id":     cfg.ID,
		"resource_name": req.ResourceName,
	})
	middleware.CreatedResponse(c, "Resource database config created successfully", cfg)
}

func (h *ResourceDBConfigHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "Invalid ID")
		return
	}

	var req models.ResourceDBConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	cfg, err := h.svc.Update(uint(id), &req)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	logActivity(c, h.activityLogSvc, "resource_db_config", "Resource DB config updated: "+cfg.ResourceName, "update", map[string]interface{}{
		"config_id":     cfg.ID,
		"resource_name": cfg.ResourceName,
	})
	middleware.SuccessResponse(c, "Resource database config updated successfully", cfg)
}

func (h *ResourceDBConfigHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "Invalid ID")
		return
	}

	cfg, err := h.svc.GetByID(uint(id))
	if err == nil {
		logActivity(c, h.activityLogSvc, "resource_db_config", "Resource DB config deleted: "+cfg.ResourceName, "delete", map[string]interface{}{
			"config_id":     id,
			"resource_name": cfg.ResourceName,
		})
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessResponse(c, "Resource database config deleted successfully", nil)
}

func (h *ResourceDBConfigHandler) TestConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "Invalid ID")
		return
	}

	if err := h.svc.TestByID(uint(id)); err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	logActivity(c, h.activityLogSvc, "resource_db_config", "Resource DB connection tested", "test_connection", map[string]interface{}{
		"config_id": id,
	})
	middleware.SuccessResponse(c, "Connection test successful", nil)
}

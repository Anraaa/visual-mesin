package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/services"
)

type ResourceQueryHandler struct {
	svc            *services.ResourceQueryService
	activityLogSvc *services.ActivityLogService
}

func NewResourceQueryHandler(svc *services.ResourceQueryService, activityLogSvc *services.ActivityLogService) *ResourceQueryHandler {
	return &ResourceQueryHandler{svc: svc, activityLogSvc: activityLogSvc}
}

func (h *ResourceQueryHandler) List(c *gin.Context) {
	resourceName := c.Param("resource")
	if resourceName == "" {
		middleware.BadRequestResponse(c, "Resource name is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))

	params := services.QueryParams{
		Page:       page,
		Limit:      limit,
		SortBy:     c.Query("sort_by"),
		SortDir:    c.Query("sort_dir"),
		Search:     c.Query("search"),
		SearchBy:   c.Query("search_by"),
		StartDate:  c.Query("start_date"),
		EndDate:    c.Query("end_date"),
		DateColumn: c.Query("date_column"),
	}

	result, err := h.svc.QueryResource(resourceName, params)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessPaginated(c, "Data retrieved successfully", result.Data, result.Total, result.Page, result.Limit)
}

func (h *ResourceQueryHandler) GetByID(c *gin.Context) {
	resourceName := c.Param("resource")
	if resourceName == "" {
		middleware.BadRequestResponse(c, "Resource name is required")
		return
	}

	idColumn := c.DefaultQuery("id_column", "id")
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		middleware.BadRequestResponse(c, "Invalid ID")
		return
	}

	result, err := h.svc.QueryByID(resourceName, idColumn, id)
	if err != nil {
		middleware.NotFoundResponse(c, err.Error())
		return
	}

	middleware.SuccessResponse(c, "Data retrieved successfully", result)
}

func (h *ResourceQueryHandler) Create(c *gin.Context) {
	resourceName := c.Param("resource")
	if resourceName == "" {
		middleware.BadRequestResponse(c, "Resource name is required")
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	result, err := h.svc.Create(resourceName, data)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	logActivity(c, h.activityLogSvc, "resource_data", "Record created in "+resourceName, "create", map[string]interface{}{
		"resource": resourceName,
		"data":     data,
	})
	middleware.CreatedResponse(c, "Record created successfully", result)
}

func (h *ResourceQueryHandler) Update(c *gin.Context) {
	resourceName := c.Param("resource")
	if resourceName == "" {
		middleware.BadRequestResponse(c, "Resource name is required")
		return
	}

	idColumn := c.DefaultQuery("id_column", "id")
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		middleware.BadRequestResponse(c, "Invalid ID")
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	result, err := h.svc.Update(resourceName, idColumn, id, data)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	logActivity(c, h.activityLogSvc, "resource_data", "Record updated in "+resourceName, "update", map[string]interface{}{
		"resource":  resourceName,
		"id_column": idColumn,
		"id":        id,
	})
	middleware.SuccessResponse(c, "Record updated successfully", result)
}

func (h *ResourceQueryHandler) Delete(c *gin.Context) {
	resourceName := c.Param("resource")
	if resourceName == "" {
		middleware.BadRequestResponse(c, "Resource name is required")
		return
	}

	idColumn := c.DefaultQuery("id_column", "id")
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		middleware.BadRequestResponse(c, "Invalid ID")
		return
	}

	if err := h.svc.Delete(resourceName, idColumn, id); err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	logActivity(c, h.activityLogSvc, "resource_data", "Record deleted from "+resourceName, "delete", map[string]interface{}{
		"resource":  resourceName,
		"id_column": idColumn,
		"id":        id,
	})
	middleware.SuccessResponse(c, "Record deleted successfully", nil)
}

func (h *ResourceQueryHandler) GetStats(c *gin.Context) {
	resourceName := c.Param("resource")
	if resourceName == "" {
		middleware.BadRequestResponse(c, "Resource name is required")
		return
	}

	durationCol := c.Query("duration_col")

	result, err := h.svc.GetStats(resourceName, durationCol)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessResponse(c, "Stats retrieved successfully", result)
}

func (h *ResourceQueryHandler) GetTrend(c *gin.Context) {
	resourceName := c.Param("resource")
	if resourceName == "" {
		middleware.BadRequestResponse(c, "Resource name is required")
		return
	}

	column := c.DefaultQuery("column", "PUD_CT")
	timeColumn := c.DefaultQuery("time_column", "")

	result, err := h.svc.GetTrend(resourceName, column, timeColumn)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessResponse(c, "Trend retrieved successfully", result)
}

func (h *ResourceQueryHandler) GetSPC(c *gin.Context) {
	resourceName := c.Param("resource")
	if resourceName == "" {
		middleware.BadRequestResponse(c, "Resource name is required")
		return
	}

	cfg := &services.SPCConfig{
		TimeColumn: c.DefaultQuery("time_col", "timestamp_record"),
		Actual:     c.DefaultQuery("actual", "aberat_control"),
		Target:     c.DefaultQuery("target", "sberat_control"),
		TolPP:      c.Query("tol_pp"),
		TolMM:      c.Query("tol_mm"),
		TolP:       c.Query("tol_p"),
		TolM:       c.Query("tol_m"),
		Limit:      50,
	}

	result, err := h.svc.GetSPC(resourceName, cfg)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessResponse(c, "SPC data retrieved successfully", result)
}

func (h *ResourceQueryHandler) GetDistribution(c *gin.Context) {
	resourceName := c.Param("resource")
	if resourceName == "" {
		middleware.BadRequestResponse(c, "Resource name is required")
		return
	}

	column := c.DefaultQuery("column", "specification")

	result, err := h.svc.GetDistribution(resourceName, column)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessResponse(c, "Distribution retrieved successfully", result)
}

func (h *ResourceQueryHandler) GetQualityTrend(c *gin.Context) {
	resourceName := c.Param("resource")
	if resourceName == "" {
		middleware.BadRequestResponse(c, "Resource name is required")
		return
	}

	timeCol := c.DefaultQuery("time_col", "eventdate")
	statusCol := c.DefaultQuery("status_col", "finaljdg")
	okValue := c.DefaultQuery("ok_value", "OK")
	ngValue := c.DefaultQuery("ng_value", "NG")

	result, err := h.svc.GetQualityTrend(resourceName, timeCol, statusCol, okValue, ngValue)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessResponse(c, "Quality trend retrieved successfully", result)
}

func (h *ResourceQueryHandler) GetColumns(c *gin.Context) {
	resourceName := c.Param("resource")
	if resourceName == "" {
		middleware.BadRequestResponse(c, "Resource name is required")
		return
	}

	result, err := h.svc.QueryResource(resourceName, services.QueryParams{Page: 1, Limit: 1})
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessResponse(c, "Columns retrieved successfully", result.Columns)
}

func (h *ResourceQueryHandler) GetJudgmentSummary(c *gin.Context) {
	resourceName := c.Param("resource")
	if resourceName == "" {
		middleware.BadRequestResponse(c, "Resource name is required")
		return
	}

	result, err := h.svc.GetJudgmentSummary(resourceName)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessResponse(c, "Judgment summary retrieved successfully", result)
}

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/services"
)

type DBConnectionHandler struct {
	svc            *services.DBConnectionService
	activityLogSvc *services.ActivityLogService
}

func NewDBConnectionHandler(svc *services.DBConnectionService, activityLogSvc *services.ActivityLogService) *DBConnectionHandler {
	return &DBConnectionHandler{svc: svc, activityLogSvc: activityLogSvc}
}

func (h *DBConnectionHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))

	conns, total, err := h.svc.List(page, limit)
	if err != nil {
		middleware.InternalErrorResponse(c, "Failed to fetch database connections")
		return
	}

	middleware.SuccessPaginated(c, "Database connections retrieved successfully", conns, total, page, limit)
}

func (h *DBConnectionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "Invalid ID")
		return
	}

	conn, err := h.svc.GetByID(uint(id))
	if err != nil {
		middleware.NotFoundResponse(c, "Database connection not found")
		return
	}

	middleware.SuccessResponse(c, "Database connection retrieved successfully", conn)
}

func (h *DBConnectionHandler) Create(c *gin.Context) {
	var req models.DBConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	conn, err := h.svc.Create(&req)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	logActivity(c, h.activityLogSvc, "db_connection", "DB connection created: "+req.Name, "create", map[string]interface{}{
		"connection_id":   conn.ID,
		"connection_name": req.Name,
		"host":            req.Host,
		"database_name":   req.DatabaseName,
	})
	middleware.CreatedResponse(c, "Database connection created successfully", conn)
}

func (h *DBConnectionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "Invalid ID")
		return
	}

	var req models.DBConnectionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	conn, err := h.svc.Update(uint(id), &req)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	logActivity(c, h.activityLogSvc, "db_connection", "DB connection updated: "+conn.Name, "update", map[string]interface{}{
		"connection_id":   conn.ID,
		"connection_name": conn.Name,
	})
	middleware.SuccessResponse(c, "Database connection updated successfully", conn)
}

func (h *DBConnectionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "Invalid ID")
		return
	}

	conn, err := h.svc.GetByID(uint(id))
	if err == nil {
		logActivity(c, h.activityLogSvc, "db_connection", "DB connection deleted: "+conn.Name, "delete", map[string]interface{}{
			"connection_id":   id,
			"connection_name": conn.Name,
		})
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessResponse(c, "Database connection deleted successfully", nil)
}

func (h *DBConnectionHandler) TestConnection(c *gin.Context) {
	idStr := c.Param("id")
	if idStr != "" {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			middleware.BadRequestResponse(c, "Invalid ID")
			return
		}

		if err := h.svc.TestByID(uint(id)); err != nil {
			middleware.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		middleware.SuccessResponse(c, "Connection test successful", nil)
		return
	}

	var req models.TestConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	if err := h.svc.TestConnection(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	middleware.SuccessResponse(c, "Connection test successful", nil)
}

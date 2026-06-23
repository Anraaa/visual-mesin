package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/services"
)

type ModuleHandler struct {
	svc *services.ResourceQueryService
}

func NewModuleHandler(svc *services.ResourceQueryService) *ModuleHandler {
	return &ModuleHandler{svc: svc}
}

func (h *ModuleHandler) ResourceList(resourceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))

		params := services.QueryParams{
			Page:     page,
			Limit:    limit,
			SortBy:   c.Query("sort_by"),
			SortDir:  c.Query("sort_dir"),
			Search:   c.Query("search"),
			SearchBy: c.Query("search_by"),
		}

		result, err := h.svc.QueryResource(resourceName, params)
		if err != nil {
			middleware.InternalErrorResponse(c, err.Error())
			return
		}

		middleware.SuccessPaginated(c, "Data retrieved successfully", result.Data, result.Total, result.Page, result.Limit)
	}
}

func (h *ModuleHandler) ResourceGetByID(resourceName, idCol string) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			middleware.BadRequestResponse(c, "Invalid ID")
			return
		}

		idColumn := c.DefaultQuery("id_column", idCol)

		result, err := h.svc.QueryByID(resourceName, idColumn, id)
		if err != nil {
			middleware.NotFoundResponse(c, err.Error())
			return
		}

		middleware.SuccessResponse(c, "Data retrieved successfully", result)
	}
}

func (h *ModuleHandler) ResourceCreate(resourceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		middleware.CreatedResponse(c, "Record created successfully", result)
	}
}

func (h *ModuleHandler) ResourceUpdate(resourceName, idCol string) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		result, err := h.svc.Update(resourceName, idCol, id, data)
		if err != nil {
			middleware.InternalErrorResponse(c, err.Error())
			return
		}

		middleware.SuccessResponse(c, "Record updated successfully", result)
	}
}

func (h *ModuleHandler) ResourceDelete(resourceName, idCol string) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			middleware.BadRequestResponse(c, "Invalid ID")
			return
		}

		if err := h.svc.Delete(resourceName, idCol, id); err != nil {
			middleware.InternalErrorResponse(c, err.Error())
			return
		}

		middleware.SuccessResponse(c, "Record deleted successfully", nil)
	}
}

func (h *ModuleHandler) ResourceBarcodeLookup(resourceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		barcode := c.Param("barcode")
		if barcode == "" {
			middleware.BadRequestResponse(c, "Barcode is required")
			return
		}

		result, err := h.svc.QueryByID(resourceName, "barcode", barcode)
		if err != nil {
			middleware.NotFoundResponse(c, err.Error())
			return
		}

		middleware.SuccessResponse(c, "Data retrieved successfully", result)
	}
}

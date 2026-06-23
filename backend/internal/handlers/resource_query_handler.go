package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/services"
)

type ResourceQueryHandler struct {
	svc *services.ResourceQueryService
}

func NewResourceQueryHandler(svc *services.ResourceQueryService) *ResourceQueryHandler {
	return &ResourceQueryHandler{svc: svc}
}

func (h *ResourceQueryHandler) Query(c *gin.Context) {
	resourceName := c.Param("resource")
	if resourceName == "" {
		middleware.BadRequestResponse(c, "Resource name is required")
		return
	}

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

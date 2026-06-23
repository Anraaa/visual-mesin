package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/services"
)

type AiSchemaMapHandler struct {
	svc *services.AiSchemaMapService
}

func NewAiSchemaMapHandler(svc *services.AiSchemaMapService) *AiSchemaMapHandler {
	return &AiSchemaMapHandler{svc: svc}
}

func (h *AiSchemaMapHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))

	items, total, err := h.svc.List(page, limit)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessPaginated(c, "AI schema maps retrieved", items, total, page, limit)
}

func (h *AiSchemaMapHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		middleware.NotFoundResponse(c, "AI schema map not found")
		return
	}
	middleware.SuccessResponse(c, "AI schema map retrieved", item)
}

func (h *AiSchemaMapHandler) Create(c *gin.Context) {
	var req models.AiSchemaMapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	item, err := h.svc.Create(&req)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.CreatedResponse(c, "AI schema map created", item)
}

func (h *AiSchemaMapHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req models.AiSchemaMapUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	item, err := h.svc.Update(uint(id), &req)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessResponse(c, "AI schema map updated", item)
}

func (h *AiSchemaMapHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.Delete(uint(id)); err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}
	middleware.SuccessResponse(c, "AI schema map deleted", nil)
}

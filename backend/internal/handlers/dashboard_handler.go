package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/services"
)

type DashboardHandler struct {
	svc *services.DashboardService
}

func NewDashboardHandler(svc *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) GetSummary(c *gin.Context) {
	data, err := h.svc.GetSummary()
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessResponse(c, "Dashboard summary retrieved successfully", data)
}

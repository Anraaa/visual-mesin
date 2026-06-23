package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/services"
)

type ActivityLogHandler struct {
	svc *services.ActivityLogService
}

func NewActivityLogHandler(svc *services.ActivityLogService) *ActivityLogHandler {
	return &ActivityLogHandler{svc: svc}
}

func (h *ActivityLogHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))

	logs, total, err := h.svc.List(page, limit)
	if err != nil {
		middleware.InternalErrorResponse(c, "Gagal mengambil activity log")
		return
	}

	middleware.SuccessPaginated(c, "Activity log berhasil diambil", logs, total, page, limit)
}

func (h *ActivityLogHandler) ListMy(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))

	logs, total, err := h.svc.ListByUser(userID.(uint), page, limit)
	if err != nil {
		middleware.InternalErrorResponse(c, "Gagal mengambil activity log")
		return
	}

	middleware.SuccessPaginated(c, "Activity log berhasil diambil", logs, total, page, limit)
}

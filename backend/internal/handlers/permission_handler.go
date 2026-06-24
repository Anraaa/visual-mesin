package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/repository"
)

type PermissionHandler struct {
	permRepo *repository.PermissionRepository
}

func NewPermissionHandler(permRepo *repository.PermissionRepository) *PermissionHandler {
	return &PermissionHandler{permRepo: permRepo}
}

func (h *PermissionHandler) List(c *gin.Context) {
	perms, err := h.permRepo.List()
	if err != nil {
		middleware.InternalErrorResponse(c, "Gagal mengambil data permission")
		return
	}
	middleware.SuccessResponse(c, "Data permission berhasil diambil", perms)
}

func (h *PermissionHandler) Create(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	perm := &models.Permission{
		Name:      req.Name,
		GuardName: "web",
	}

	if err := h.permRepo.Create(perm); err != nil {
		middleware.InternalErrorResponse(c, "Gagal membuat permission")
		return
	}

	middleware.CreatedResponse(c, "Permission berhasil dibuat", perm)
}

func (h *PermissionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID tidak valid")
		return
	}

	perm := &models.Permission{ID: uint(id)}
	if err := h.permRepo.Delete(perm); err != nil {
		middleware.InternalErrorResponse(c, "Gagal menghapus permission")
		return
	}

	middleware.SuccessResponse(c, "Permission berhasil dihapus", nil)
}

package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/repository"
)

type RoleHandler struct {
	roleRepo *repository.RoleRepository
}

func NewRoleHandler(roleRepo *repository.RoleRepository) *RoleHandler {
	return &RoleHandler{roleRepo: roleRepo}
}

func (h *RoleHandler) List(c *gin.Context) {
	roles, err := h.roleRepo.List()
	if err != nil {
		middleware.InternalErrorResponse(c, "Gagal mengambil data role")
		return
	}

	type roleItem struct {
		ID          uint                `json:"id"`
		Name        string              `json:"name"`
		GuardName   string              `json:"guard_name"`
		Permissions []models.Permission `json:"permissions"`
	}

	items := make([]roleItem, 0)
	for _, r := range roles {
		perms, _ := h.roleRepo.GetPermissions(r.ID)
		items = append(items, roleItem{
			ID:          r.ID,
			Name:        r.Name,
			GuardName:   r.GuardName,
			Permissions: perms,
		})
	}

	middleware.SuccessResponse(c, "Data role berhasil diambil", items)
}

func (h *RoleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID tidak valid")
		return
	}

	role, err := h.roleRepo.FindByID(uint(id))
	if err != nil {
		middleware.NotFoundResponse(c, "Role tidak ditemukan")
		return
	}

	perms, _ := h.roleRepo.GetPermissions(role.ID)

	middleware.SuccessResponse(c, "Data role berhasil diambil", gin.H{
		"id":          role.ID,
		"name":        role.Name,
		"guard_name":  role.GuardName,
		"permissions": perms,
	})
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	role := &models.Role{
		Name:      req.Name,
		GuardName: "web",
	}

	if err := h.roleRepo.Create(role); err != nil {
		middleware.InternalErrorResponse(c, "Gagal membuat role")
		return
	}

	middleware.CreatedResponse(c, "Role berhasil dibuat", role)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID tidak valid")
		return
	}

	role := &models.Role{ID: uint(id)}
	if err := h.roleRepo.Delete(role); err != nil {
		middleware.InternalErrorResponse(c, "Gagal menghapus role")
		return
	}

	middleware.SuccessResponse(c, "Role berhasil dihapus", nil)
}

type AssignPermissionRequest struct {
	PermissionID uint `json:"permission_id" binding:"required"`
}

func (h *RoleHandler) AssignPermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID tidak valid")
		return
	}

	var req AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	if err := h.roleRepo.AssignPermission(uint(id), req.PermissionID); err != nil {
		middleware.InternalErrorResponse(c, "Gagal menetapkan permission")
		return
	}

	middleware.SuccessResponse(c, "Permission berhasil ditetapkan", nil)
}

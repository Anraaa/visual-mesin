package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/repository"
	"github.com/anraaa/visual-mesin/internal/services"
)

type RoleHandler struct {
	roleRepo       *repository.RoleRepository
	permRepo       *repository.PermissionRepository
	activityLogSvc *services.ActivityLogService
}

func NewRoleHandler(roleRepo *repository.RoleRepository, permRepo *repository.PermissionRepository, activityLogSvc *services.ActivityLogService) *RoleHandler {
	return &RoleHandler{roleRepo: roleRepo, permRepo: permRepo, activityLogSvc: activityLogSvc}
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

	logActivity(c, h.activityLogSvc, "role", "Role created: "+role.Name, "create", map[string]interface{}{
		"role_id":   role.ID,
		"role_name": role.Name,
	})
	middleware.CreatedResponse(c, "Role berhasil dibuat", role)
}

func (h *RoleHandler) Delete(c *gin.Context) {
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

	if err := h.roleRepo.Delete(role); err != nil {
		middleware.InternalErrorResponse(c, "Gagal menghapus role")
		return
	}

	logActivity(c, h.activityLogSvc, "role", "Role deleted: "+role.Name, "delete", map[string]interface{}{
		"role_id":   role.ID,
		"role_name": role.Name,
	})
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

	perm, err := h.permRepo.FindByID(req.PermissionID)
	if err != nil {
		middleware.NotFoundResponse(c, "Permission tidak ditemukan")
		return
	}

	if err := h.roleRepo.AssignPermission(uint(id), req.PermissionID); err != nil {
		middleware.InternalErrorResponse(c, "Gagal menetapkan permission")
		return
	}

	logActivity(c, h.activityLogSvc, "role", "Permission "+perm.Name+" assigned to role", "assign_permission", map[string]interface{}{
		"role_id":       id,
		"permission_id": perm.ID,
		"permission":    perm.Name,
	})
	middleware.SuccessResponse(c, "Permission berhasil ditetapkan", nil)
}

func (h *RoleHandler) RevokePermission(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID role tidak valid")
		return
	}

	permID, err := strconv.ParseUint(c.Param("permId"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID permission tidak valid")
		return
	}

	perm, err := h.permRepo.FindByID(uint(permID))
	if err == nil {
		logActivity(c, h.activityLogSvc, "role", "Permission "+perm.Name+" revoked from role", "revoke_permission", map[string]interface{}{
			"role_id":       roleID,
			"permission_id": permID,
			"permission":    perm.Name,
		})
	}

	if err := h.roleRepo.RevokePermission(uint(roleID), uint(permID)); err != nil {
		middleware.InternalErrorResponse(c, "Gagal mencabut permission")
		return
	}

	middleware.SuccessResponse(c, "Permission berhasil dicabut", nil)
}

func (h *RoleHandler) SyncPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID role tidak valid")
		return
	}

	var req struct {
		Permissions []string `json:"permissions" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	var permIDs []uint
	for _, name := range req.Permissions {
		perm, err := h.permRepo.FindByName(name)
		if err != nil {
			perm = &models.Permission{
				Name:      name,
				GuardName: "web",
			}
			if err := h.permRepo.Create(perm); err != nil {
				middleware.InternalErrorResponse(c, "Gagal membuat permission: "+name)
				return
			}
		}
		permIDs = append(permIDs, perm.ID)
	}

	if err := h.roleRepo.RevokeAllPermissions(uint(id)); err != nil {
		middleware.InternalErrorResponse(c, "Gagal menyinkronisasi permission")
		return
	}

	for _, permID := range permIDs {
		if err := h.roleRepo.AssignPermission(uint(id), permID); err != nil {
			middleware.InternalErrorResponse(c, "Gagal menetapkan permission")
			return
		}
	}

	logActivity(c, h.activityLogSvc, "role", "Permissions synced for role", "sync_permissions", map[string]interface{}{
		"role_id":     id,
		"permissions": req.Permissions,
	})
	middleware.SuccessResponse(c, "Permission berhasil disinkronisasi", nil)
}

package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/repository"
	"github.com/anraaa/visual-mesin/internal/services"
)

type UserHandler struct {
	userRepo       *repository.UserRepository
	roleRepo       *repository.RoleRepository
	activityLogSvc *services.ActivityLogService
}

func NewUserHandler(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository, activityLogSvc *services.ActivityLogService) *UserHandler {
	return &UserHandler{userRepo: userRepo, roleRepo: roleRepo, activityLogSvc: activityLogSvc}
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))

	users, total, err := h.userRepo.List(page, limit)
	if err != nil {
		middleware.InternalErrorResponse(c, "Gagal mengambil data user")
		return
	}

	type userItem struct {
		ID        uint          `json:"id"`
		NIP       *string       `json:"nip"`
		UserName  string        `json:"user_name"`
		Email     *string       `json:"email"`
		UserLevel string        `json:"user_level"`
		Roles     []models.Role `json:"roles"`
	}

	items := make([]userItem, 0)
	for _, u := range users {
		roles, _ := h.userRepo.GetRoles(u.ID)
		items = append(items, userItem{
			ID:        u.ID,
			NIP:       u.NIP,
			UserName:  u.UserName,
			Email:     u.Email,
			UserLevel: u.UserLevel,
			Roles:     roles,
		})
	}

	middleware.SuccessPaginated(c, "Data user berhasil diambil", items, total, page, limit)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID tidak valid")
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		middleware.NotFoundResponse(c, "User tidak ditemukan")
		return
	}

	roles, _ := h.userRepo.GetRoles(user.ID)
	perms, _ := h.userRepo.GetPermissions(user.ID)
	permNames := make([]string, 0, len(perms))
	for _, p := range perms {
		permNames = append(permNames, p.Name)
	}

	c.JSON(200, models.APIResponse{
		Success: true,
		Message: "Data user berhasil diambil",
		Data: gin.H{
			"id":          user.ID,
			"nip":         user.NIP,
			"user_name":   user.UserName,
			"email":       user.Email,
			"user_level":  user.UserLevel,
			"department":  user.Department,
			"jabatan":     user.Jabatan,
			"roles":       roles,
			"permissions": permNames,
		},
	})
}

type CreateUserRequest struct {
	NIP       string `json:"nip"`
	UserName  string `json:"user_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	UserLevel string `json:"user_level" binding:"omitempty,oneof=admin eng tech prod"`
	RoleIDs   []uint `json:"role_ids"`
}

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	existing, _ := h.userRepo.FindByEmail(req.Email)
	if existing != nil {
		middleware.BadRequestResponse(c, "Email sudah terdaftar")
		return
	}

	userLevel := req.UserLevel
	if userLevel == "" {
		userLevel = "prod"
	}

	user := &models.User{
		UserName:  req.UserName,
		Email:     &req.Email,
		Password:  req.Password,
		UserLevel: userLevel,
	}
	if req.NIP != "" {
		user.NIP = &req.NIP
	}

	if err := h.userRepo.Create(user); err != nil {
		middleware.InternalErrorResponse(c, "Gagal membuat user")
		return
	}

	for _, roleID := range req.RoleIDs {
		if err := h.userRepo.AssignRole(user.ID, roleID); err != nil {
			continue
		}
	}

	logActivity(c, h.activityLogSvc, "user", "User created: "+user.UserName, "create", map[string]interface{}{
		"user_id":   user.ID,
		"user_name": user.UserName,
	})

	roles, _ := h.userRepo.GetRoles(user.ID)
	middleware.CreatedResponse(c, "User berhasil dibuat", gin.H{
		"id":         user.ID,
		"user_name":  user.UserName,
		"email":      user.Email,
		"user_level": user.UserLevel,
		"roles":      roles,
	})
}

type UpdateUserRequest struct {
	NIP        *string `json:"nip"`
	UserName   *string `json:"user_name"`
	Email      *string `json:"email"`
	Password   *string `json:"password"`
	OldPassword *string `json:"old_password"`
	UserLevel  *string `json:"user_level"`
	Department *string `json:"department"`
	Jabatan    *string `json:"jabatan"`
	RoleIDs    *[]uint `json:"role_ids"`
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID tidak valid")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		middleware.NotFoundResponse(c, "User tidak ditemukan")
		return
	}

	if req.NIP != nil {
		user.NIP = req.NIP
	}
	if req.UserName != nil {
		user.UserName = *req.UserName
	}
	if req.Email != nil {
		user.Email = req.Email
	}
	if req.UserLevel != nil {
		validLevels := map[string]bool{"admin": true, "eng": true, "tech": true, "prod": true}
		if !validLevels[*req.UserLevel] {
			middleware.BadRequestResponse(c, "Level user tidak valid")
			return
		}
		user.UserLevel = *req.UserLevel
	}
	if req.Department != nil {
		user.Department = req.Department
	}
	if req.Jabatan != nil {
		user.Jabatan = req.Jabatan
	}
	if req.Password != nil && *req.Password != "" {
		if req.OldPassword == nil || *req.OldPassword == "" {
			middleware.BadRequestResponse(c, "Password lama wajib diisi untuk mengganti password")
			return
		}
		if user.Password != *req.OldPassword {
			middleware.BadRequestResponse(c, "Password lama tidak sesuai")
			return
		}
		user.Password = *req.Password
	}
	if req.RoleIDs != nil {
		if err := h.userRepo.RevokeAllRoles(user.ID); err != nil {
			middleware.InternalErrorResponse(c, "Gagal mensinkronisasi role")
			return
		}
		for _, roleID := range *req.RoleIDs {
			if err := h.userRepo.AssignRole(user.ID, roleID); err != nil {
				continue
			}
		}
	}

	if err := h.userRepo.Update(user); err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	logActivity(c, h.activityLogSvc, "user", "User updated: "+user.UserName, "update", map[string]interface{}{
		"user_id":   user.ID,
		"user_name": user.UserName,
	})

	roles, _ := h.userRepo.GetRoles(user.ID)
	middleware.SuccessResponse(c, "User berhasil diupdate", gin.H{
		"id":         user.ID,
		"user_name":  user.UserName,
		"email":      user.Email,
		"user_level": user.UserLevel,
		"roles":      roles,
	})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID tidak valid")
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		middleware.NotFoundResponse(c, "User tidak ditemukan")
		return
	}

	if err := h.userRepo.Delete(uint(id)); err != nil {
		middleware.InternalErrorResponse(c, "Gagal menghapus user")
		return
	}

	logActivity(c, h.activityLogSvc, "user", "User deleted: "+user.UserName, "delete", map[string]interface{}{
		"user_id":   user.ID,
		"user_name": user.UserName,
	})
	middleware.SuccessResponse(c, "User berhasil dihapus", nil)
}

type AssignRoleRequest struct {
	RoleID uint `json:"role_id" binding:"required"`
}

func (h *UserHandler) AssignRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID tidak valid")
		return
	}

	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	role, err := h.roleRepo.FindByID(req.RoleID)
	if err != nil {
		middleware.NotFoundResponse(c, "Role tidak ditemukan")
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		middleware.NotFoundResponse(c, "User tidak ditemukan")
		return
	}

	if err := h.userRepo.AssignRole(uint(id), req.RoleID); err != nil {
		middleware.InternalErrorResponse(c, "Gagal menetapkan role")
		return
	}

	logActivity(c, h.activityLogSvc, "user", "Role "+role.Name+" assigned to "+user.UserName, "assign_role", map[string]interface{}{
		"user_id":   user.ID,
		"user_name": user.UserName,
		"role_id":   role.ID,
		"role_name": role.Name,
	})
	middleware.SuccessResponse(c, "Role berhasil ditetapkan", nil)
}

type SyncRolesRequest struct {
	RoleIDs []uint `json:"role_ids" binding:"required"`
}

func (h *UserHandler) SyncRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID tidak valid")
		return
	}

	var req SyncRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		middleware.NotFoundResponse(c, "User tidak ditemukan")
		return
	}

	if err := h.userRepo.RevokeAllRoles(user.ID); err != nil {
		middleware.InternalErrorResponse(c, "Gagal mensinkronisasi role")
		return
	}
	for _, roleID := range req.RoleIDs {
		if err := h.userRepo.AssignRole(user.ID, roleID); err != nil {
			continue
		}
	}

	logActivity(c, h.activityLogSvc, "user", "Roles synced for "+user.UserName, "sync_roles", map[string]interface{}{
		"user_id":   user.ID,
		"user_name": user.UserName,
		"role_ids":  req.RoleIDs,
	})

	roles, _ := h.userRepo.GetRoles(user.ID)
	middleware.SuccessResponse(c, "Role berhasil disinkronisasi", gin.H{
		"roles": roles,
	})
}

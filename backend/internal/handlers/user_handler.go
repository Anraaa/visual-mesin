package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/repository"
)

type UserHandler struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewUserHandler(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo, roleRepo: roleRepo}
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
		ID        uint            `json:"id"`
		NIP       *string         `json:"nip"`
		UserName  string          `json:"user_name"`
		Email     *string         `json:"email"`
		UserLevel string          `json:"user_level"`
		Roles     []models.Role   `json:"roles"`
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

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Data user berhasil diambil",
		Data: gin.H{
			"id":         user.ID,
			"nip":        user.NIP,
			"user_name":  user.UserName,
			"email":      user.Email,
			"user_level": user.UserLevel,
			"department": user.Department,
			"jabatan":    user.Jabatan,
			"roles":      roles,
		},
	})
}

type UpdateUserRequest struct {
	NIP       *string `json:"nip"`
	UserName  *string `json:"user_name"`
	Email     *string `json:"email"`
	UserLevel *string `json:"user_level" binding:"omitempty,oneof=admin eng tech prod"`
	Department *string `json:"department"`
	Jabatan   *string `json:"jabatan"`
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
		user.UserLevel = *req.UserLevel
	}
	if req.Department != nil {
		user.Department = req.Department
	}
	if req.Jabatan != nil {
		user.Jabatan = req.Jabatan
	}

	if err := h.userRepo.Update(user); err != nil {
		middleware.InternalErrorResponse(c, "Gagal mengupdate user")
		return
	}

	middleware.SuccessResponse(c, "User berhasil diupdate", user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		middleware.BadRequestResponse(c, "ID tidak valid")
		return
	}

	if err := h.userRepo.Delete(uint(id)); err != nil {
		middleware.InternalErrorResponse(c, "Gagal menghapus user")
		return
	}

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

	if err := h.userRepo.AssignRole(uint(id), req.RoleID); err != nil {
		middleware.InternalErrorResponse(c, "Gagal menetapkan role")
		return
	}

	middleware.SuccessResponse(c, "Role berhasil ditetapkan", nil)
}

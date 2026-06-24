package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/services"
)

type AuthHandler struct {
	authSvc        *services.AuthService
	activityLogSvc *services.ActivityLogService
}

func NewAuthHandler(authSvc *services.AuthService, activityLogSvc *services.ActivityLogService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, activityLogSvc: activityLogSvc}
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.APIResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Validasi gagal: " + err.Error(),
		})
		return
	}

	resp, err := h.authSvc.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if resp.User.ID > 0 && h.activityLogSvc != nil {
		h.activityLogSvc.Log(resp.User.ID, "auth", "User logged in", "login", map[string]interface{}{
			"user_name": resp.User.UserName,
			"email":     resp.User.Email,
		})
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login berhasil",
		Data:    resp,
	})
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration data"
// @Success 201 {object} models.APIResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Validasi gagal: " + err.Error(),
		})
		return
	}

	user, err := h.authSvc.Register(req)
	if err != nil {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if h.activityLogSvc != nil {
		h.activityLogSvc.Log(user.ID, "auth", "User registered", "register", map[string]interface{}{
			"user_name": user.UserName,
			"email":     user.Email,
		})
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Registrasi berhasil",
		Data:    user,
	})
}

// Me godoc
// @Summary Get current user profile
// @Description Return authenticated user profile with roles
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")

	profile, err := h.authSvc.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "User tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile berhasil diambil",
		Data:    profile,
	})
}

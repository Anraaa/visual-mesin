package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/anraaa/visual-mesin/internal/models"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userLevel, exists := c.Get("user_level")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Message: "Akses ditolak",
			})
			return
		}

		level := userLevel.(string)
		for _, role := range roles {
			if level == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Anda tidak memiliki akses ke resource ini",
		})
	}
}

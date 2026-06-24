package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/anraaa/visual-mesin/internal/services"
)

func logActivity(c *gin.Context, svc *services.ActivityLogService, logName, description, event string, properties map[string]interface{}) {
	if svc == nil {
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		return
	}
	svc.Log(userID.(uint), logName, description, event, properties)
}

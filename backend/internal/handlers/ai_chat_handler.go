package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/anraaa/visual-mesin/internal/ai"
	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/services"
)

type AiChatHandler struct {
	pipeline *ai.ChatPipeline
	histSvc  *services.AiChatHistoryService
}

func NewAiChatHandler(pipeline *ai.ChatPipeline, histSvc *services.AiChatHistoryService) *AiChatHandler {
	return &AiChatHandler{pipeline: pipeline, histSvc: histSvc}
}

func (h *AiChatHandler) Chat(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.BadRequestResponse(c, err.Error())
		return
	}

	if req.SessionID == "" {
		req.SessionID = time.Now().Format("session_150405")
	}

	result, err := h.pipeline.Process(req.Question, userID.(uint), req.SessionID)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessResponse(c, "OK", gin.H{
		"session_id":      req.SessionID,
		"answer":          result.Answer,
		"detected_intent": result.DetectedIntent,
		"generated_sql":   result.GeneratedSQL,
		"sql_status":      result.SQLStatus,
		"query_result":    result.QueryResult,
		"latency":         result.Latency,
	})
}

func (h *AiChatHandler) History(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		middleware.BadRequestResponse(c, "Session ID required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	items, total, err := h.histSvc.ListBySession(sessionID, page, limit)
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessPaginated(c, "Chat history retrieved", items, total, page, limit)
}

func (h *AiChatHandler) Sessions(c *gin.Context) {
	userID, _ := c.Get("user_id")

	sessions, err := h.histSvc.ListSessions(userID.(uint))
	if err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessResponse(c, "Sessions retrieved", sessions)
}

func (h *AiChatHandler) DeleteSession(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		middleware.BadRequestResponse(c, "Session ID required")
		return
	}

	if err := h.histSvc.DeleteSession(sessionID); err != nil {
		middleware.InternalErrorResponse(c, err.Error())
		return
	}

	middleware.SuccessResponse(c, "Session deleted", nil)
}

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/anraaa/visual-mesin/internal/models"
)

func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessPaginated(c *gin.Context, message string, data interface{}, total int64, page, limit int) {
	lastPage := int(total) / limit
	if int(total)%limit > 0 {
		lastPage++
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: &models.Meta{
			CurrentPage: page,
			PerPage:     limit,
			Total:       total,
			LastPage:    lastPage,
		},
	})
}

func CreatedResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func BadRequestResponse(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, models.APIResponse{
		Success: false,
		Message: message,
	})
}

func NotFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, models.APIResponse{
		Success: false,
		Message: message,
	})
}

func InternalErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, models.APIResponse{
		Success: false,
		Message: message,
	})
}

func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, models.APIResponse{
		Success: false,
		Message: message,
	})
}

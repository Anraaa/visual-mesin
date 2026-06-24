package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/anraaa/visual-mesin/internal/services"
)

func setupTest() (*gin.Engine, *services.JWTService) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	jwtSvc := services.NewJWTService("test-secret-key", time.Hour)
	return r, jwtSvc
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	r, _ := setupTest()
	r.GET("/test", AuthMiddleware(nil), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "Missing authorization header")
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	r, _ := setupTest()
	r.GET("/test", AuthMiddleware(nil), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat token123")
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid authorization format")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	r, jwtSvc := setupTest()
	r.GET("/test", AuthMiddleware(jwtSvc), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-here")
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid or expired token")
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	r, jwtSvc := setupTest()
	r.GET("/test", AuthMiddleware(jwtSvc), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		userLevel, _ := c.Get("user_level")
		email, _ := c.Get("email")
		c.JSON(200, gin.H{
			"user_id":    userID,
			"user_level": userLevel,
			"email":      email,
		})
	})

	token, _ := jwtSvc.GenerateToken(1, "admin", "admin@test.com")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"user_id":1`)
	assert.Contains(t, w.Body.String(), `"user_level":"admin"`)
	assert.Contains(t, w.Body.String(), `"email":"admin@test.com"`)
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRBACTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestRequireRole_NoUserLevel(t *testing.T) {
	r := setupRBACTest()
	r.GET("/admin", RequireRole("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
	assert.Contains(t, w.Body.String(), "Akses ditolak")
}

func TestRequireRole_Allowed(t *testing.T) {
	r := setupRBACTest()
	r.GET("/admin", func(c *gin.Context) {
		c.Set("user_level", "admin")
		c.Next()
	}, RequireRole("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestRequireRole_Denied(t *testing.T) {
	r := setupRBACTest()
	r.GET("/admin", func(c *gin.Context) {
		c.Set("user_level", "prod")
		c.Next()
	}, RequireRole("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
	assert.Contains(t, w.Body.String(), "akses")
}

func TestRequireRole_MultipleRoles(t *testing.T) {
	r := setupRBACTest()
	r.GET("/eng", func(c *gin.Context) {
		c.Set("user_level", "eng")
		c.Next()
	}, RequireRole("admin", "eng"), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/eng", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/database"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/repository"
	"github.com/stretchr/testify/assert"
)

func TestAuditHandler(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试数据库
	dbConfig := &database.Config{
		Type:       "sqlite",
		SQLiteFile: ":memory:",
	}
	db, err := database.NewConnection(dbConfig)
	assert.NoError(t, err)

	// 执行数据库迁移
	err = models.AutoMigrate(db.DB())
	assert.NoError(t, err)

	// 创建Repository管理器
	repoManager := repository.NewRepositoryManager(db.DB())

	// 创建测试审计日志
	testLog := &models.AuditLog{
		TenantID:  "test-tenant-id",
		UserID:    "test-user-id",
		Action:    "create",
		Resource:  "hosts",
		Details:   `{"method":"POST","path":"/api/v1/hosts"}`,
		Result:    "success",
		CreatedAt: time.Now(),
	}
	err = repoManager.AuditLog().Create(context.Background(), testLog)
	assert.NoError(t, err)

	// 创建处理器
	handler := NewAuditHandler(repoManager)

	t.Run("GetAuditLogs", func(t *testing.T) {
		router := gin.New()
		router.GET("/audit/logs", func(c *gin.Context) {
			// 模拟认证中间件设置的租户ID
			c.Set("tenant_id", "test-tenant-id")
			handler.GetAuditLogs(c)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/audit/logs", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "test-user-id")
		assert.Contains(t, w.Body.String(), "create")
		assert.Contains(t, w.Body.String(), "hosts")
	})

	t.Run("GetAuditLogsByUser", func(t *testing.T) {
		router := gin.New()
		router.GET("/audit/users/:user_id/logs", handler.GetAuditLogsByUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/audit/users/test-user-id/logs", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "test-user-id")
	})

	t.Run("GetAuditLogsByResource", func(t *testing.T) {
		router := gin.New()
		router.GET("/audit/resources/:resource/logs", func(c *gin.Context) {
			// 模拟认证中间件设置的租户ID
			c.Set("tenant_id", "test-tenant-id")
			handler.GetAuditLogsByResource(c)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/audit/resources/hosts/logs", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "hosts")
	})

	t.Run("GetAuditLogs without tenant_id should return 401", func(t *testing.T) {
		router := gin.New()
		router.GET("/audit/logs", handler.GetAuditLogs)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/audit/logs", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "缺少租户信息")
	})

	t.Run("GetAuditLogsByUser with empty user_id should return 400", func(t *testing.T) {
		router := gin.New()
		router.GET("/audit/users/:user_id/logs", handler.GetAuditLogsByUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/audit/users//logs", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

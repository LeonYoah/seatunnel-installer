package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/database"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/repository"
	"github.com/seatunnel/enterprise-platform/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestAuditMiddleware(t *testing.T) {
	// 初始化日志
	err := logger.Init(logger.DefaultConfig())
	assert.NoError(t, err)
	log := logger.Get()

	// 创建测试数据库
	dbConfig := &database.Config{
		Type:       "sqlite",
		SQLiteFile: ":memory:",
	}
	db, err := database.NewConnection(dbConfig)
	assert.NoError(t, err)

	// 执行数据库迁移
	err = db.Migrate()
	if err != nil {
		t.Logf("Migration error: %v", err)
		// 尝试手动创建表
		err = models.AutoMigrate(db.DB())
		assert.NoError(t, err)
	}

	// 创建Repository管理器
	repoManager := repository.NewRepositoryManager(db.DB())

	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	router := gin.New()
	router.Use(AuditMiddleware(repoManager, log))

	// 添加测试路由
	router.POST("/api/v1/test", func(c *gin.Context) {
		// 设置用户信息（模拟认证中间件的行为）
		c.Set("user_id", "test-user-id")
		c.Set("tenant_id", "test-tenant-id")

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	router.GET("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// 测试写操作（应该记录审计日志）
	t.Run("POST request should create audit log", func(t *testing.T) {
		requestBody := map[string]string{"name": "test"}
		bodyBytes, _ := json.Marshal(requestBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/test", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 等待异步审计日志写入完成
		time.Sleep(100 * time.Millisecond)

		// 验证审计日志是否创建
		logs, total, err := repoManager.AuditLog().GetByTenantID(context.Background(), "test-tenant-id", 0, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		if assert.Len(t, logs, 1) {
			auditLog := logs[0]
			assert.Equal(t, "test-user-id", auditLog.UserID)
			assert.Equal(t, "test-tenant-id", auditLog.TenantID)
			assert.Equal(t, "create", auditLog.Action)
			assert.Equal(t, "test", auditLog.Resource)
			assert.Equal(t, "success", auditLog.Result)
			assert.NotEmpty(t, auditLog.Details)
		}
	})

	// 测试读操作（不应该记录审计日志）
	t.Run("GET request should not create audit log", func(t *testing.T) {
		// 清空之前的审计日志
		db.DB().Exec("DELETE FROM audit_logs")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 等待一段时间确保没有异步写入
		time.Sleep(100 * time.Millisecond)

		// 验证没有审计日志创建
		logs, total, err := repoManager.AuditLog().GetByTenantID(context.Background(), "test-tenant-id", 0, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, logs, 0)
	})

	// 测试没有用户信息的请求（不应该记录审计日志）
	t.Run("Request without user info should not create audit log", func(t *testing.T) {
		// 清空之前的审计日志
		db.DB().Exec("DELETE FROM audit_logs")

		// 创建不设置用户信息的路由
		router2 := gin.New()
		router2.Use(AuditMiddleware(repoManager, log))
		router2.POST("/api/v1/test2", func(c *gin.Context) {
			// 不设置用户信息
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/test2", nil)
		router2.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 等待一段时间确保没有异步写入
		time.Sleep(100 * time.Millisecond)

		// 验证没有审计日志创建
		logs, total, err := repoManager.AuditLog().GetByTenantID(context.Background(), "test-tenant-id", 0, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, logs, 0)
	})
}

func TestGetActionFromMethod(t *testing.T) {
	tests := []struct {
		method   string
		expected string
	}{
		{"POST", "create"},
		{"PUT", "update"},
		{"PATCH", "update"},
		{"DELETE", "delete"},
		{"GET", "unknown"},
		{"HEAD", "unknown"},
	}

	for _, test := range tests {
		result := getActionFromMethod(test.method)
		assert.Equal(t, test.expected, result, "Method %s should return %s", test.method, test.expected)
	}
}

func TestGetResourceFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/api/v1/hosts", "hosts"},
		{"/api/v1/clusters/123", "clusters"},
		{"/api/v1/tasks/456/runs", "tasks"},
		{"/hosts", "hosts"},
		{"/", "unknown"},
		{"", "unknown"},
	}

	for _, test := range tests {
		result := getResourceFromPath(test.path)
		assert.Equal(t, test.expected, result, "Path %s should return %s", test.path, test.expected)
	}
}

func TestExtractResourceID(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/api/v1/hosts/550e8400-e29b-41d4-a716-446655440000", "550e8400-e29b-41d4-a716-446655440000"},
		{"/api/v1/clusters/123", ""},
		{"/api/v1/tasks", ""},
		{"/hosts/550e8400-e29b-41d4-a716-446655440001/status", "550e8400-e29b-41d4-a716-446655440001"},
	}

	for _, test := range tests {
		result := extractResourceID(test.path)
		assert.Equal(t, test.expected, result, "Path %s should return %s", test.path, test.expected)
	}
}

func TestIsWriteOperation(t *testing.T) {
	tests := []struct {
		method   string
		expected bool
	}{
		{"POST", true},
		{"PUT", true},
		{"PATCH", true},
		{"DELETE", true},
		{"GET", false},
		{"HEAD", false},
		{"OPTIONS", false},
	}

	for _, test := range tests {
		result := isWriteOperation(test.method)
		assert.Equal(t, test.expected, result, "Method %s should return %t", test.method, test.expected)
	}
}

func TestContainsSensitiveData(t *testing.T) {
	tests := []struct {
		body     string
		expected bool
	}{
		{`{"username": "test", "password": "secret"}`, true},
		{`{"name": "test", "value": "normal"}`, false},
		{`{"api_key": "secret123"}`, true},
		{`{"token": "abc123"}`, true},
		{`{"credential": "secret"}`, true},
		{`{"normal": "data"}`, false},
	}

	for _, test := range tests {
		result := containsSensitiveData([]byte(test.body))
		assert.Equal(t, test.expected, result, "Body %s should return %t", test.body, test.expected)
	}
}

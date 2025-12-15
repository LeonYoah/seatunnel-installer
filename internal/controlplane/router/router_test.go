package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/auth"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/database"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/repository"
	"github.com/seatunnel/enterprise-platform/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
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

	// 创建Repository管理器
	repoManager := repository.NewRepositoryManager(db.DB())

	// 创建JWT服务
	jwtService := auth.NewJWTService("test-secret", time.Hour, time.Hour*24)

	router := NewRouter(log, repoManager, jwtService)

	assert.NotNil(t, router)
	assert.NotNil(t, router.engine)
	assert.NotNil(t, router.logger)
}

func TestHealthEndpoint(t *testing.T) {
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

	// 创建Repository管理器
	repoManager := repository.NewRepositoryManager(db.DB())

	// 创建JWT服务
	jwtService := auth.NewJWTService("test-secret", time.Hour, time.Hour*24)

	router := NewRouter(log, repoManager, jwtService)
	router.SetupRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "SeaTunnel Control Plane is running")
}

func TestAPIV1Routes(t *testing.T) {
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

	// 创建Repository管理器
	repoManager := repository.NewRepositoryManager(db.DB())

	// 创建JWT服务
	jwtService := auth.NewJWTService("test-secret", time.Hour, time.Hour*24)

	router := NewRouter(log, repoManager, jwtService)
	router.SetupRoutes()

	// 测试主机管理路由
	testCases := []struct {
		method   string
		path     string
		expected int
	}{
		{"GET", "/api/v1/hosts", http.StatusOK},
		{"POST", "/api/v1/hosts", http.StatusOK},
		{"GET", "/api/v1/hosts/123", http.StatusOK},
		{"PUT", "/api/v1/hosts/123", http.StatusOK},
		{"DELETE", "/api/v1/hosts/123", http.StatusOK},
		{"GET", "/api/v1/clusters", http.StatusOK},
		{"POST", "/api/v1/clusters", http.StatusOK},
		{"GET", "/api/v1/tasks", http.StatusOK},
		{"POST", "/api/v1/tasks", http.StatusOK},
		{"GET", "/api/v1/auth/login", http.StatusNotFound},    // POST only
		{"POST", "/api/v1/auth/login", http.StatusBadRequest}, // 需要请求体
		{"GET", "/api/v1/agent/install.sh", http.StatusOK},
		{"GET", "/api/v1/plugins", http.StatusOK},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tc.method, tc.path, nil)
		router.engine.ServeHTTP(w, req)

		assert.Equal(t, tc.expected, w.Code, "Failed for %s %s", tc.method, tc.path)
	}
}

func TestCORSHeaders(t *testing.T) {
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

	// 创建Repository管理器
	repoManager := repository.NewRepositoryManager(db.DB())

	// 创建JWT服务
	jwtService := auth.NewJWTService("test-secret", time.Hour, time.Hour*24)

	router := NewRouter(log, repoManager, jwtService)
	router.SetupRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	router.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

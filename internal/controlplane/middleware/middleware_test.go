package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/seatunnel/enterprise-platform/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestZapLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 初始化日志
	err := logger.Init(logger.DefaultConfig())
	assert.NoError(t, err)

	log := logger.Get()

	router := gin.New()
	router.Use(ZapLogger(log))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 初始化日志
	err := logger.Init(logger.DefaultConfig())
	assert.NoError(t, err)

	log := logger.Get()

	router := gin.New()
	router.Use(ErrorHandler(log))

	// 测试正常请求
	router.GET("/normal", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// 测试panic处理
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	// 测试错误处理
	router.GET("/error", func(c *gin.Context) {
		c.Error(errors.New("test error"))
	})

	// 测试正常请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/normal", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 测试panic处理
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/panic", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "服务器内部错误")

	// 测试错误处理
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/error", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

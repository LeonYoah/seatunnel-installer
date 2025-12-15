package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/router"
	"github.com/seatunnel/enterprise-platform/pkg/logger"
	"go.uber.org/zap"
)

// Server HTTP服务器
type Server struct {
	httpServer *http.Server
	router     *router.Router
	logger     *zap.Logger
	port       int
}

// NewServer 创建新的服务器实例
func NewServer(port int) (*Server, error) {
	// 初始化日志
	err := logger.Init(&logger.Config{
		Level:      "info",
		OutputPath: "logs/control-plane.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
		Console:    true,
	})
	if err != nil {
		return nil, fmt.Errorf("初始化日志失败: %w", err)
	}

	log := logger.Get()

	// 创建路由器
	r := router.NewRouter(log)
	r.SetupRoutes()

	// 创建HTTP服务器
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r.GetEngine(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		router:     r,
		logger:     log,
		port:       port,
	}, nil
}

// Start 启动服务器
func (s *Server) Start() error {
	s.logger.Info("启动Control Plane服务器",
		zap.Int("port", s.port),
		zap.String("addr", s.httpServer.Addr),
	)

	// 在goroutine中启动服务器
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	s.logger.Info("Control Plane服务器启动成功",
		zap.String("url", fmt.Sprintf("http://localhost:%d", s.port)),
		zap.String("health_check", fmt.Sprintf("http://localhost:%d/health", s.port)),
		zap.String("api_base", fmt.Sprintf("http://localhost:%d/api/v1", s.port)),
	)

	// 等待中断信号
	return s.waitForShutdown()
}

// waitForShutdown 等待关闭信号
func (s *Server) waitForShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	s.logger.Info("收到关闭信号，开始优雅关闭服务器...")

	// 创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("服务器关闭失败", zap.Error(err))
		return err
	}

	s.logger.Info("服务器已优雅关闭")
	return nil
}

// Stop 停止服务器
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

// GetLogger 获取日志器
func (s *Server) GetLogger() *zap.Logger {
	return s.logger
}

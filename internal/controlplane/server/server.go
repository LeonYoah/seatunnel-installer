package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/auth"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/database"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/repository"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/router"
	"github.com/seatunnel/enterprise-platform/pkg/config"
	"github.com/seatunnel/enterprise-platform/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Server HTTP服务器
type Server struct {
	httpServer  *http.Server
	router      *router.Router
	logger      *zap.Logger
	db          *gorm.DB
	repoManager repository.RepositoryManager
	jwtService  *auth.JWTService
	rbacService *auth.RBACService
	port        int
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

	// 加载配置
	cfg, err := config.Load("")
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	// 转换数据库配置
	dbConfig := &database.Config{
		Type:            cfg.Database.Type,
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		Database:        cfg.Database.Database,
		Username:        cfg.Database.Username,
		Password:        cfg.Database.Password,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: time.Duration(cfg.Database.ConnMaxLifetime) * time.Second,
		ConnMaxIdleTime: time.Duration(cfg.Database.ConnMaxIdleTime) * time.Second,
		SQLiteFile:      cfg.Database.SQLiteFile,
		SSLMode:         cfg.Database.SSLMode,
		ServiceName:     cfg.Database.ServiceName,
		SID:             cfg.Database.SID,
	}

	// 初始化数据库
	dbManager, err := database.NewManager(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("初始化数据库失败: %w", err)
	}

	// 获取数据库连接
	db := dbManager.DB()

	// 创建Repository管理器
	repoManager := dbManager.Repository()

	// 创建JWT服务
	jwtService := auth.NewJWTService(
		cfg.JWT.SecretKey,
		time.Duration(cfg.JWT.AccessTokenTTL)*time.Minute,
		time.Duration(cfg.JWT.RefreshTokenTTL)*time.Hour,
	)

	// 创建RBAC服务
	rbacService := auth.NewRBACService(repoManager)

	// 初始化内置角色和权限
	if err := rbacService.InitializeBuiltInRoles(context.Background()); err != nil {
		log.Warn("初始化内置角色失败", zap.Error(err))
	}

	// 创建路由器
	r := router.NewRouter(log, repoManager, jwtService)
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
		httpServer:  httpServer,
		router:      r,
		logger:      log,
		db:          db,
		repoManager: repoManager,
		jwtService:  jwtService,
		rbacService: rbacService,
		port:        port,
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

// CreateDefaultAdminUser 创建默认管理员用户
func (s *Server) CreateDefaultAdminUser() error {
	ctx := context.Background()

	// 检查是否已存在admin用户
	existingUser, err := s.repoManager.User().GetByUsername(ctx, "admin")
	if err == nil && existingUser != nil {
		s.logger.Info("默认管理员用户已存在，跳过创建")
		return nil
	}

	// 创建默认租户（如果不存在）
	var defaultTenant *models.Tenant
	tenants, _, err := s.repoManager.Tenant().List(ctx, 0, 1)
	if err != nil || len(tenants) == 0 {
		defaultTenant = &models.Tenant{
			Name:        "默认租户",
			Description: "系统默认租户",
		}
		if err := s.repoManager.Tenant().Create(ctx, defaultTenant); err != nil {
			return fmt.Errorf("创建默认租户失败: %w", err)
		}
		s.logger.Info("创建默认租户成功", zap.String("tenant_id", defaultTenant.ID))
	} else {
		defaultTenant = tenants[0]
	}

	// 创建默认工作空间（如果不存在）
	var defaultWorkspace *models.Workspace
	workspaces, _, err := s.repoManager.Workspace().GetByTenantID(ctx, defaultTenant.ID, 0, 1)
	if err != nil || len(workspaces) == 0 {
		defaultWorkspace = &models.Workspace{
			TenantID: defaultTenant.ID,
			Name:     "默认工作空间",
		}
		if err := s.repoManager.Workspace().Create(ctx, defaultWorkspace); err != nil {
			return fmt.Errorf("创建默认工作空间失败: %w", err)
		}
		s.logger.Info("创建默认工作空间成功", zap.String("workspace_id", defaultWorkspace.ID))
	} else {
		defaultWorkspace = workspaces[0]
	}

	// 获取owner角色
	ownerRole, err := s.repoManager.Role().GetByName(ctx, "", "owner")
	if err != nil {
		return fmt.Errorf("获取owner角色失败: %w", err)
	}

	// 创建默认管理员用户
	adminUser := &models.User{
		TenantID:    defaultTenant.ID,
		WorkspaceID: defaultWorkspace.ID,
		Username:    "admin",
		Email:       "admin@seatunnel.local",
		Password:    "admin123", // 明文密码，会在BeforeCreate中自动加密
		Status:      "active",
	}

	if err := s.repoManager.User().Create(ctx, adminUser); err != nil {
		return fmt.Errorf("创建管理员用户失败: %w", err)
	}

	// 为管理员用户分配owner角色
	if err := s.repoManager.User().AssignRole(ctx, adminUser.ID, ownerRole.ID); err != nil {
		return fmt.Errorf("分配角色失败: %w", err)
	}

	s.logger.Info("创建默认管理员用户成功",
		zap.String("user_id", adminUser.ID),
		zap.String("username", adminUser.Username),
	)

	return nil
}

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/auth"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/handlers"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/middleware"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/repository"
	"go.uber.org/zap"
)

// Router API路由器
type Router struct {
	engine      *gin.Engine
	logger      *zap.Logger
	repoManager repository.RepositoryManager
	jwtService  *auth.JWTService
	authHandler *handlers.AuthHandler
}

// NewRouter 创建新的路由器
func NewRouter(logger *zap.Logger, repoManager repository.RepositoryManager, jwtService *auth.JWTService) *Router {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	// 添加中间件
	engine.Use(middleware.ZapLogger(logger))
	engine.Use(middleware.ErrorHandler(logger))
	engine.Use(middleware.CORS())
	engine.Use(gin.Recovery())

	// 创建认证服务和处理器
	authService := auth.NewAuthService(repoManager, jwtService)
	authHandler := handlers.NewAuthHandler(authService)

	return &Router{
		engine:      engine,
		logger:      logger,
		repoManager: repoManager,
		jwtService:  jwtService,
		authHandler: authHandler,
	}
}

// GetEngine 获取Gin引擎
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

// SetupRoutes 设置路由
func (r *Router) SetupRoutes() {
	// 健康检查
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "SeaTunnel Control Plane is running",
		})
	})

	// API版本控制
	v1 := r.engine.Group("/api/v1")
	{
		r.setupV1Routes(v1)
	}
}

// setupV1Routes 设置v1版本的路由
func (r *Router) setupV1Routes(v1 *gin.RouterGroup) {
	// 认证相关路由
	auth := v1.Group("/auth")
	{
		r.setupAuthRoutes(auth)
	}

	// 主机管理路由
	hosts := v1.Group("/hosts")
	{
		r.setupHostRoutes(hosts)
	}

	// 集群管理路由
	clusters := v1.Group("/clusters")
	{
		r.setupClusterRoutes(clusters)
	}

	// 任务管理路由
	tasks := v1.Group("/tasks")
	{
		r.setupTaskRoutes(tasks)
	}

	// 部署管理路由
	deploy := v1.Group("/deploy")
	{
		r.setupDeployRoutes(deploy)
	}

	// Agent管理路由
	agent := v1.Group("/agent")
	{
		r.setupAgentRoutes(agent)
	}

	// 插件市场路由
	plugins := v1.Group("/plugins")
	{
		r.setupPluginRoutes(plugins)
	}

	// 审计日志路由
	audit := v1.Group("/audit")
	{
		r.setupAuditRoutes(audit)
	}
}

// setupAuthRoutes 设置认证路由
func (r *Router) setupAuthRoutes(authGroup *gin.RouterGroup) {
	// 公开路由（不需要认证）
	authGroup.POST("/login", r.authHandler.Login)
	authGroup.POST("/refresh", r.authHandler.RefreshToken)

	// 需要认证的路由
	authenticated := authGroup.Group("")
	authenticated.Use(middleware.AuthMiddleware(r.jwtService, r.repoManager))
	{
		authenticated.GET("/me", r.authHandler.GetCurrentUser)
		authenticated.POST("/logout", r.authHandler.Logout)
	}
}

// setupHostRoutes 设置主机管理路由
func (r *Router) setupHostRoutes(hosts *gin.RouterGroup) {
	// TODO: 实现主机管理路由
	hosts.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "获取主机列表接口待实现"})
	})
	hosts.POST("", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "注册主机接口待实现"})
	})
	hosts.GET("/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "获取主机详情接口待实现"})
	})
	hosts.PUT("/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "更新主机接口待实现"})
	})
	hosts.DELETE("/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "删除主机接口待实现"})
	})
}

// setupClusterRoutes 设置集群管理路由
func (r *Router) setupClusterRoutes(clusters *gin.RouterGroup) {
	// TODO: 实现集群管理路由
	clusters.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "获取集群列表接口待实现"})
	})
	clusters.POST("", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "注册集群接口待实现"})
	})
	clusters.GET("/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "获取集群详情接口待实现"})
	})
	clusters.PUT("/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "更新集群接口待实现"})
	})
	clusters.DELETE("/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "删除集群接口待实现"})
	})

	// 节点管理子路由
	clusters.GET("/:id/nodes", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "获取集群节点列表接口待实现"})
	})
	clusters.POST("/:id/nodes", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "添加集群节点接口待实现"})
	})
	clusters.DELETE("/:id/nodes/:nodeId", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "删除集群节点接口待实现"})
	})
}

// setupTaskRoutes 设置任务管理路由
func (r *Router) setupTaskRoutes(tasks *gin.RouterGroup) {
	// TODO: 实现任务管理路由
	tasks.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "获取任务列表接口待实现"})
	})
	tasks.POST("", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "创建任务接口待实现"})
	})
	tasks.GET("/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "获取任务详情接口待实现"})
	})
	tasks.PUT("/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "更新任务接口待实现"})
	})
	tasks.DELETE("/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "删除任务接口待实现"})
	})

	// 任务运行相关路由
	tasks.POST("/:id/runs", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "提交任务运行接口待实现"})
	})
	tasks.GET("/:id/runs", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "获取任务运行历史接口待实现"})
	})
	tasks.GET("/:id/runs/:runId", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "获取任务运行详情接口待实现"})
	})
	tasks.POST("/:id/runs/:runId/stop", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "停止任务运行接口待实现"})
	})
}

// setupDeployRoutes 设置部署管理路由
func (r *Router) setupDeployRoutes(deploy *gin.RouterGroup) {
	// TODO: 实现部署管理路由
	deploy.POST("/precheck", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "环境预检查接口待实现"})
	})
	deploy.POST("/start", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "开始部署接口待实现"})
	})
	deploy.GET("/status/:deployId", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "获取部署状态接口待实现"})
	})
	deploy.GET("/logs/:deployId", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "获取部署日志接口待实现"})
	})
}

// setupAgentRoutes 设置Agent管理路由
func (r *Router) setupAgentRoutes(agent *gin.RouterGroup) {
	// TODO: 实现Agent管理路由
	agent.GET("/install.sh", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Agent安装脚本接口待实现"})
	})
	agent.GET("/download", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Agent二进制下载接口待实现"})
	})
	agent.POST("/register", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Agent注册接口待实现"})
	})
	agent.POST("/heartbeat", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Agent心跳接口待实现"})
	})
}

// setupPluginRoutes 设置插件市场路由
func (r *Router) setupPluginRoutes(plugins *gin.RouterGroup) {
	// TODO: 实现插件市场路由
	plugins.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "获取插件列表接口待实现"})
	})
	plugins.GET("/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "获取插件详情接口待实现"})
	})
	plugins.POST("/:id/install", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "安装插件接口待实现"})
	})
	plugins.POST("/:id/upgrade", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "升级插件接口待实现"})
	})
}

// setupAuditRoutes 设置审计日志路由
func (r *Router) setupAuditRoutes(audit *gin.RouterGroup) {
	// TODO: 实现审计日志路由
	audit.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "获取审计日志接口待实现"})
	})
	audit.GET("/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "获取审计日志详情接口待实现"})
	})
}

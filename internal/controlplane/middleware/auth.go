package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/api"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/auth"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/repository"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(jwtService *auth.JWTService, repoManager repository.RepositoryManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			api.Unauthorized(c, "缺少认证令牌")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			api.Unauthorized(c, "无效的认证令牌格式")
			c.Abort()
			return
		}

		// 验证令牌
		claims, err := jwtService.ValidateAccessToken(tokenParts[1])
		if err != nil {
			if err == auth.ErrExpiredToken {
				api.Unauthorized(c, "令牌已过期")
			} else {
				api.Unauthorized(c, "无效的认证令牌")
			}
			c.Abort()
			return
		}

		// 获取用户信息（包含角色和权限）
		user, err := repoManager.User().GetWithRoles(c.Request.Context(), claims.UserID)
		if err != nil {
			api.Unauthorized(c, "用户不存在")
			c.Abort()
			return
		}

		// 检查用户状态
		if !user.IsActive() {
			api.Unauthorized(c, "用户已被禁用")
			c.Abort()
			return
		}

		// 将用户信息和声明存储到上下文
		c.Set("user", user)
		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("workspace_id", claims.WorkspaceID)

		c.Next()
	}
}

// RequirePermission 权限检查中间件
func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息
		userInterface, exists := c.Get("user")
		if !exists {
			api.Unauthorized(c, "未认证的用户")
			c.Abort()
			return
		}

		user, ok := userInterface.(*models.User)
		if !ok {
			api.InternalError(c, "用户信息格式错误", nil)
			c.Abort()
			return
		}

		// 检查权限
		if !user.HasPermission(resource, action) {
			api.Forbidden(c, "权限不足：需要"+resource+":"+action+"权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole 角色检查中间件
func RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息
		userInterface, exists := c.Get("user")
		if !exists {
			api.Unauthorized(c, "未认证的用户")
			c.Abort()
			return
		}

		user, ok := userInterface.(*models.User)
		if !ok {
			api.InternalError(c, "用户信息格式错误", nil)
			c.Abort()
			return
		}

		// 检查角色
		if !user.HasRole(roleName) {
			api.Forbidden(c, "角色权限不足：需要"+roleName+"角色")
			c.Abort()
			return
		}

		c.Next()
	}
}

// TenantIsolation 租户隔离中间件
func TenantIsolation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取租户ID
		tenantID, exists := c.Get("tenant_id")
		if !exists {
			api.Unauthorized(c, "缺少租户信息")
			c.Abort()
			return
		}

		// 将租户ID添加到查询参数中，确保所有查询都包含租户隔离
		c.Set("filter_tenant_id", tenantID)
		c.Next()
	}
}

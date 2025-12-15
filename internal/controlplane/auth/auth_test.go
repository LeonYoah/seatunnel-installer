package auth

import (
	"testing"
	"time"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"github.com/stretchr/testify/assert"
)

func TestJWTService(t *testing.T) {
	// 创建JWT服务
	jwtService := NewJWTService("test-secret-key", 15*time.Minute, 24*time.Hour)

	// 创建测试用户
	user := &models.User{
		ID:          "test-user-id",
		TenantID:    "test-tenant-id",
		WorkspaceID: "test-workspace-id",
		Username:    "testuser",
		Roles: []models.Role{
			{Name: "admin"},
			{Name: "viewer"},
		},
	}

	t.Run("生成和验证访问令牌", func(t *testing.T) {
		// 生成访问令牌
		token, err := jwtService.GenerateAccessToken(user)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// 验证访问令牌
		claims, err := jwtService.ValidateAccessToken(token)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.TenantID, claims.TenantID)
		assert.Equal(t, user.WorkspaceID, claims.WorkspaceID)
		assert.Equal(t, user.Username, claims.Username)
		assert.Equal(t, []string{"admin", "viewer"}, claims.Roles)
	})

	t.Run("生成和验证刷新令牌", func(t *testing.T) {
		// 生成刷新令牌
		token, err := jwtService.GenerateRefreshToken(user.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// 验证刷新令牌
		userID, err := jwtService.ValidateRefreshToken(token)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, userID)
	})

	t.Run("验证无效令牌", func(t *testing.T) {
		// 验证无效的访问令牌
		_, err := jwtService.ValidateAccessToken("invalid-token")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidToken, err)

		// 验证无效的刷新令牌
		_, err = jwtService.ValidateRefreshToken("invalid-token")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidToken, err)
	})

	t.Run("刷新访问令牌", func(t *testing.T) {
		// 生成刷新令牌
		refreshToken, err := jwtService.GenerateRefreshToken(user.ID)
		assert.NoError(t, err)

		// 刷新访问令牌
		newAccessToken, err := jwtService.RefreshAccessToken(refreshToken, user)
		assert.NoError(t, err)
		assert.NotEmpty(t, newAccessToken)

		// 验证新的访问令牌
		claims, err := jwtService.ValidateAccessToken(newAccessToken)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
	})
}

func TestUserPermissions(t *testing.T) {
	// 创建测试用户和角色
	user := &models.User{
		ID:       "test-user-id",
		Username: "testuser",
		Status:   "active",
		Roles: []models.Role{
			{
				Name: "admin",
				Permissions: []models.Permission{
					{Permission: "user:read"},
					{Permission: "user:create"},
					{Permission: "cluster:manage"},
				},
			},
			{
				Name: "viewer",
				Permissions: []models.Permission{
					{Permission: "task:read"},
				},
			},
		},
	}

	t.Run("检查用户权限", func(t *testing.T) {
		// 用户应该有的权限
		assert.True(t, user.HasPermission("user", "read"))
		assert.True(t, user.HasPermission("user", "create"))
		assert.True(t, user.HasPermission("cluster", "manage"))
		assert.True(t, user.HasPermission("task", "read"))

		// 用户不应该有的权限
		assert.False(t, user.HasPermission("user", "delete"))
		assert.False(t, user.HasPermission("system", "config"))
	})

	t.Run("检查用户角色", func(t *testing.T) {
		// 用户应该有的角色
		assert.True(t, user.HasRole("admin"))
		assert.True(t, user.HasRole("viewer"))

		// 用户不应该有的角色
		assert.False(t, user.HasRole("owner"))
		assert.False(t, user.HasRole("guest"))
	})

	t.Run("检查用户状态", func(t *testing.T) {
		// 活跃用户
		assert.True(t, user.IsActive())

		// 非活跃用户
		user.Status = "inactive"
		assert.False(t, user.IsActive())

		user.Status = "locked"
		assert.False(t, user.IsActive())
	})
}

package auth

import (
	"context"

	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/repository"
	"gorm.io/gorm"
)

// RBACService RBAC服务
type RBACService struct {
	repoManager repository.RepositoryManager
}

// NewRBACService 创建RBAC服务
func NewRBACService(repoManager repository.RepositoryManager) *RBACService {
	return &RBACService{
		repoManager: repoManager,
	}
}

// InitializeBuiltInRoles 初始化内置角色和权限
func (s *RBACService) InitializeBuiltInRoles(ctx context.Context) error {
	// 首先初始化权限
	if err := s.initializePermissions(ctx); err != nil {
		return err
	}

	// 然后初始化角色
	if err := s.initializeRoles(ctx); err != nil {
		return err
	}

	return nil
}

// initializePermissions 初始化权限
func (s *RBACService) initializePermissions(ctx context.Context) error {
	permissions := []models.Permission{
		// 租户管理权限
		{Permission: "tenant:read", Resource: "tenant", Action: "read", Description: "查看租户信息"},
		{Permission: "tenant:create", Resource: "tenant", Action: "create", Description: "创建租户"},
		{Permission: "tenant:update", Resource: "tenant", Action: "update", Description: "更新租户信息"},
		{Permission: "tenant:delete", Resource: "tenant", Action: "delete", Description: "删除租户"},

		// 工作空间管理权限
		{Permission: "workspace:read", Resource: "workspace", Action: "read", Description: "查看工作空间"},
		{Permission: "workspace:create", Resource: "workspace", Action: "create", Description: "创建工作空间"},
		{Permission: "workspace:update", Resource: "workspace", Action: "update", Description: "更新工作空间"},
		{Permission: "workspace:delete", Resource: "workspace", Action: "delete", Description: "删除工作空间"},

		// 用户管理权限
		{Permission: "user:read", Resource: "user", Action: "read", Description: "查看用户信息"},
		{Permission: "user:create", Resource: "user", Action: "create", Description: "创建用户"},
		{Permission: "user:update", Resource: "user", Action: "update", Description: "更新用户信息"},
		{Permission: "user:delete", Resource: "user", Action: "delete", Description: "删除用户"},

		// 角色管理权限
		{Permission: "role:read", Resource: "role", Action: "read", Description: "查看角色信息"},
		{Permission: "role:create", Resource: "role", Action: "create", Description: "创建角色"},
		{Permission: "role:update", Resource: "role", Action: "update", Description: "更新角色信息"},
		{Permission: "role:delete", Resource: "role", Action: "delete", Description: "删除角色"},

		// 主机管理权限
		{Permission: "host:read", Resource: "host", Action: "read", Description: "查看主机信息"},
		{Permission: "host:create", Resource: "host", Action: "create", Description: "注册主机"},
		{Permission: "host:update", Resource: "host", Action: "update", Description: "更新主机信息"},
		{Permission: "host:delete", Resource: "host", Action: "delete", Description: "删除主机"},

		// 集群管理权限
		{Permission: "cluster:read", Resource: "cluster", Action: "read", Description: "查看集群信息"},
		{Permission: "cluster:create", Resource: "cluster", Action: "create", Description: "创建集群"},
		{Permission: "cluster:update", Resource: "cluster", Action: "update", Description: "更新集群信息"},
		{Permission: "cluster:delete", Resource: "cluster", Action: "delete", Description: "删除集群"},
		{Permission: "cluster:deploy", Resource: "cluster", Action: "deploy", Description: "部署集群"},
		{Permission: "cluster:manage", Resource: "cluster", Action: "manage", Description: "管理集群（启动、停止、重启）"},

		// 任务管理权限
		{Permission: "task:read", Resource: "task", Action: "read", Description: "查看任务信息"},
		{Permission: "task:create", Resource: "task", Action: "create", Description: "创建任务"},
		{Permission: "task:update", Resource: "task", Action: "update", Description: "更新任务信息"},
		{Permission: "task:delete", Resource: "task", Action: "delete", Description: "删除任务"},
		{Permission: "task:execute", Resource: "task", Action: "execute", Description: "执行任务"},

		// 运行管理权限
		{Permission: "run:read", Resource: "run", Action: "read", Description: "查看运行记录"},
		{Permission: "run:stop", Resource: "run", Action: "stop", Description: "停止运行"},

		// 凭证管理权限
		{Permission: "secret:read", Resource: "secret", Action: "read", Description: "查看凭证信息"},
		{Permission: "secret:create", Resource: "secret", Action: "create", Description: "创建凭证"},
		{Permission: "secret:update", Resource: "secret", Action: "update", Description: "更新凭证"},
		{Permission: "secret:delete", Resource: "secret", Action: "delete", Description: "删除凭证"},

		// 审计日志权限
		{Permission: "audit:read", Resource: "audit", Action: "read", Description: "查看审计日志"},
		{Permission: "audit:export", Resource: "audit", Action: "export", Description: "导出审计日志"},

		// 系统管理权限
		{Permission: "system:read", Resource: "system", Action: "read", Description: "查看系统信息"},
		{Permission: "system:config", Resource: "system", Action: "config", Description: "配置系统参数"},
	}

	for _, perm := range permissions {
		// 检查权限是否已存在
		existing, err := s.repoManager.Permission().GetByPermission(ctx, perm.Permission)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if existing != nil {
			continue // 权限已存在，跳过
		}

		// 创建权限
		if err := s.repoManager.Permission().Create(ctx, &perm); err != nil {
			return err
		}
	}

	return nil
}

// initializeRoles 初始化角色
func (s *RBACService) initializeRoles(ctx context.Context) error {
	// 定义内置角色
	roles := []struct {
		Name        string
		DisplayName string
		Description string
		Permissions []string
	}{
		{
			Name:        "owner",
			DisplayName: "所有者",
			Description: "拥有所有权限的超级管理员",
			Permissions: []string{
				"tenant:read", "tenant:create", "tenant:update", "tenant:delete",
				"workspace:read", "workspace:create", "workspace:update", "workspace:delete",
				"user:read", "user:create", "user:update", "user:delete",
				"role:read", "role:create", "role:update", "role:delete",
				"host:read", "host:create", "host:update", "host:delete",
				"cluster:read", "cluster:create", "cluster:update", "cluster:delete", "cluster:deploy", "cluster:manage",
				"task:read", "task:create", "task:update", "task:delete", "task:execute",
				"run:read", "run:stop",
				"secret:read", "secret:create", "secret:update", "secret:delete",
				"audit:read", "audit:export",
				"system:read", "system:config",
			},
		},
		{
			Name:        "admin",
			DisplayName: "管理员",
			Description: "拥有大部分管理权限，但不能管理租户和系统配置",
			Permissions: []string{
				"workspace:read", "workspace:create", "workspace:update", "workspace:delete",
				"user:read", "user:create", "user:update", "user:delete",
				"role:read", "role:create", "role:update", "role:delete",
				"host:read", "host:create", "host:update", "host:delete",
				"cluster:read", "cluster:create", "cluster:update", "cluster:delete", "cluster:deploy", "cluster:manage",
				"task:read", "task:create", "task:update", "task:delete", "task:execute",
				"run:read", "run:stop",
				"secret:read", "secret:create", "secret:update", "secret:delete",
				"audit:read",
			},
		},
		{
			Name:        "viewer",
			DisplayName: "查看者",
			Description: "只有查看权限，不能进行修改操作",
			Permissions: []string{
				"workspace:read",
				"user:read",
				"role:read",
				"host:read",
				"cluster:read",
				"task:read",
				"run:read",
				"secret:read",
				"audit:read",
			},
		},
	}

	for _, roleData := range roles {
		// 检查角色是否已存在（使用空租户ID查找内置角色）
		existing, err := s.repoManager.Role().GetByName(ctx, "", roleData.Name)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if existing != nil {
			continue // 角色已存在，跳过
		}

		// 创建角色
		role := models.Role{
			TenantID:    "", // 内置角色不属于任何租户
			Name:        roleData.Name,
			DisplayName: roleData.DisplayName,
			Description: roleData.Description,
			IsBuiltIn:   true,
			Status:      "active",
		}

		if err := s.repoManager.Role().Create(ctx, &role); err != nil {
			return err
		}

		// 为角色分配权限
		for _, permName := range roleData.Permissions {
			perm, err := s.repoManager.Permission().GetByPermission(ctx, permName)
			if err != nil {
				continue // 权限不存在，跳过
			}

			if err := s.repoManager.Role().AssignPermission(ctx, role.ID, perm.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

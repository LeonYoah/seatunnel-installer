package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seatunnel/enterprise-platform/internal/controlplane/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// 使用内存SQLite数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}

	// 自动迁移
	err = models.AutoMigrate(db)
	if err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	return db
}

func TestBaseRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := newBaseRepository[models.Tenant](db)
	ctx := context.Background()

	// 测试创建
	tenant := &models.Tenant{
		ID:          uuid.New().String(),
		Name:        "test-tenant",
		Description: "测试租户",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, tenant)
	if err != nil {
		t.Fatalf("创建租户失败: %v", err)
	}

	// 测试根据ID获取
	found, err := repo.GetByID(ctx, tenant.ID)
	if err != nil {
		t.Fatalf("根据ID获取租户失败: %v", err)
	}
	if found == nil {
		t.Fatal("未找到租户")
	}
	if found.Name != tenant.Name {
		t.Errorf("租户名称不匹配: 期望 %s, 实际 %s", tenant.Name, found.Name)
	}

	// 测试更新
	found.Description = "更新后的描述"
	err = repo.Update(ctx, found)
	if err != nil {
		t.Fatalf("更新租户失败: %v", err)
	}

	// 验证更新
	updated, err := repo.GetByID(ctx, tenant.ID)
	if err != nil {
		t.Fatalf("获取更新后的租户失败: %v", err)
	}
	if updated.Description != "更新后的描述" {
		t.Errorf("租户描述未更新: 期望 '更新后的描述', 实际 '%s'", updated.Description)
	}

	// 测试列表查询
	tenants, total, err := repo.List(ctx, 0, 10)
	if err != nil {
		t.Fatalf("获取租户列表失败: %v", err)
	}
	if total != 1 {
		t.Errorf("租户总数不匹配: 期望 1, 实际 %d", total)
	}
	if len(tenants) != 1 {
		t.Errorf("租户列表长度不匹配: 期望 1, 实际 %d", len(tenants))
	}

	// 测试删除
	err = repo.Delete(ctx, tenant.ID)
	if err != nil {
		t.Fatalf("删除租户失败: %v", err)
	}

	// 验证删除
	deleted, err := repo.GetByID(ctx, tenant.ID)
	if err != nil {
		t.Fatalf("验证删除失败: %v", err)
	}
	if deleted != nil {
		t.Error("租户应该已被删除")
	}

	t.Log("基础Repository测试通过")
}

func TestTenantRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := newTenantRepository(db)
	ctx := context.Background()

	// 创建测试租户
	tenant := &models.Tenant{
		ID:          uuid.New().String(),
		Name:        "test-tenant",
		Description: "测试租户",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, tenant)
	if err != nil {
		t.Fatalf("创建租户失败: %v", err)
	}

	// 测试根据名称获取租户
	found, err := repo.GetByName(ctx, tenant.Name)
	if err != nil {
		t.Fatalf("根据名称获取租户失败: %v", err)
	}
	if found == nil {
		t.Fatal("未找到租户")
	}
	if found.ID != tenant.ID {
		t.Errorf("租户ID不匹配: 期望 %s, 实际 %s", tenant.ID, found.ID)
	}

	// 测试获取活跃租户列表
	activeTenants, total, err := repo.ListActive(ctx, 0, 10)
	if err != nil {
		t.Fatalf("获取活跃租户列表失败: %v", err)
	}
	if total != 1 {
		t.Errorf("活跃租户数量不匹配: 期望 1, 实际 %d", total)
	}
	if len(activeTenants) != 1 {
		t.Errorf("活跃租户列表长度不匹配: 期望 1, 实际 %d", len(activeTenants))
	}

	t.Log("租户Repository测试通过")
}

func TestRepositoryManager(t *testing.T) {
	db := setupTestDB(t)
	manager := NewRepositoryManager(db)

	// 测试获取各个Repository
	if manager.Tenant() == nil {
		t.Error("租户Repository为空")
	}
	if manager.Workspace() == nil {
		t.Error("工作空间Repository为空")
	}
	if manager.Host() == nil {
		t.Error("主机Repository为空")
	}
	if manager.Cluster() == nil {
		t.Error("集群Repository为空")
	}
	if manager.Node() == nil {
		t.Error("节点Repository为空")
	}
	if manager.Task() == nil {
		t.Error("任务Repository为空")
	}
	if manager.Run() == nil {
		t.Error("运行Repository为空")
	}
	if manager.AuditLog() == nil {
		t.Error("审计日志Repository为空")
	}
	if manager.Secret() == nil {
		t.Error("凭证Repository为空")
	}

	t.Log("Repository管理器测试通过")
}
